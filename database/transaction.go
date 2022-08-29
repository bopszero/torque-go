package database

import (
	"database/sql"
	"fmt"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gorm.io/gorm"
)

type (
	TransactionFunc             func(dbTxn *gorm.DB) error
	TransactionOnCommitCallback func() error
)

func TransactionFromContext(ctx comcontext.Context, alias string, execFunc TransactionFunc) error {
	var (
		ctxKey        = genContextKey(alias)
		ctxTxnWrapper = ctx.Get(ctxKey)
		txnWrapper    *tTransactionWrapper
	)
	if ctxTxnWrapper != nil {
		txnWrapper = ctxTxnWrapper.(*tTransactionWrapper)
	} else {
		txnWrapper = newTransactionWrapper(ctx, alias, nil)
		ctx.Set(ctxKey, txnWrapper)
		defer ctx.Set(ctxKey, nil)
	}
	return txnWrapper.execute(execFunc)
}

func genContextKey(alias string) string {
	return fmt.Sprintf("db:txn:%s", alias)
}

func Atomic(ctx comcontext.Context, alias string, execFunc TransactionFunc) error {
	return TransactionFromContext(ctx, alias, execFunc)
}

func OnCommit(ctx comcontext.Context, alias string, handler TransactionOnCommitCallback) error {
	var (
		ctxKey        = genContextKey(alias)
		ctxTxnWrapper = ctx.Get(ctxKey)
	)
	if ctxTxnWrapper == nil {
		return handler()
	}
	ctxTxnWrapper.(*tTransactionWrapper).RegisterOnCommitCallback(handler)
	return nil
}

type tTransactionWrapper struct {
	ctx               comcontext.Context
	curTxn            *gorm.DB
	txnOpts           *sql.TxOptions
	onCommitCallbacks []TransactionOnCommitCallback
}

func newTransactionWrapper(ctx comcontext.Context, dbAlias string, opts *sql.TxOptions) *tTransactionWrapper {
	return &tTransactionWrapper{
		ctx:     ctx,
		curTxn:  GetDbF(dbAlias).DB,
		txnOpts: opts,
	}
}

func (this *tTransactionWrapper) execute(execFunc TransactionFunc) (err error) {
	var (
		panicked = true
		db       = this.curTxn
	)
	defer func() {
		this.curTxn = db
	}()

	if committer, ok := db.Statement.ConnPool.(gorm.TxCommitter); ok && committer != nil {
		if !db.DisableNestedTransaction {
			spID := this.genSavePointID(execFunc)
			err = db.SavePoint(spID).Error
			defer func() {
				if panicked || err != nil {
					db.RollbackTo(spID)
				}
			}()
		}
		if err == nil {
			txn := db.Session(&gorm.Session{})
			this.curTxn = txn
			err = execFunc(txn)
		}
	} else {
		txn := db.Begin(this.txnOpts)
		defer func() {
			if panicked || err != nil {
				txn.Rollback()
			}
		}()
		if err = txn.Error; err == nil {
			this.curTxn = txn
			err = execFunc(txn)
		}
		if err == nil {
			err = txn.Commit().Error
			if err == nil {
				this.executeOnCommitCallbacks()
			}
		}
	}

	panicked = false
	return
}

func (this *tTransactionWrapper) genSavePointID(execFunc TransactionFunc) string {
	return fmt.Sprintf("sid_%p_%s", execFunc, comutils.NewUuidCode())
}

func (this *tTransactionWrapper) RegisterOnCommitCallback(callback TransactionOnCommitCallback) {
	this.onCommitCallbacks = append(this.onCommitCallbacks, callback)
}

func (this *tTransactionWrapper) executeOnCommitCallbacks() {
	if len(this.onCommitCallbacks) == 0 {
		return
	}
	logger := comlogging.GetLogger()
	defer func() {
		if errLike := recover(); errLike != nil {
			err := comutils.ToError(errLike)
			logger.
				WithType("db").
				WithContext(this.ctx).
				WithError(comlogging.NewPreloadSentryError(err)).
				Error("transaction on commit panic")
		}
	}()
	for _, callback := range this.onCommitCallbacks {
		if err := callback(); err != nil {
			logger.
				WithType("db").
				WithContext(this.ctx).
				WithError(comlogging.NewPreloadSentryError(err)).
				Error("transaction on commit callback failed")
		}
	}
}
