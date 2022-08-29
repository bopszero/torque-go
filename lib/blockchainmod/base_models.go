package blockchainmod

import "gitlab.com/snap-clickstaff/torque-go/lib/utils"

type baseTransaction struct {
	ownerAddress string
}

func (this *baseTransaction) SetOwnerAddress(address string) {
	this.ownerAddress = address
}

func (this *baseTransaction) GetTypeCode() string {
	return ""
}

func (this *baseTransaction) GetInputs() ([]Input, error) {
	return nil, utils.IssueErrorf("this kind of txn doesn't have inputs concept")
}

func (this *baseTransaction) GetOutputs() ([]Output, error) {
	return nil, utils.IssueErrorf("this kind of txn doesn't have outputs concept")
}

func (this *baseTransaction) GetRC20Transfers(TokenMeta) ([]RC20Transfer, error) {
	return nil, utils.IssueErrorf("this kind of txn doesn't have RC20 transfer concept")
}
