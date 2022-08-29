package authmod

import (
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func PasswordPrepareNonce(username string) (nonceInfo PasswordNonce, err error) {
	user, err := usermod.GetUserByUsername(username)
	if utils.IsOurError(err, constants.ErrorCodeUserNotFound) {
		// Prevent brute-force scan username
		dummyNonce := PasswordNonce{
			ID:       comutils.NewUuid4Code(),
			Value:    comutils.RandomBytesF(PasswordNonceSize),
			Username: username,
			Salt:     comutils.HashSha256String(username + SaltDummyString)[:PasswordSaltSize],
		}
		return dummyNonce, nil
	}
	if err != nil {
		return
	}
	passHash, err := LoadPasswordArgon2iHash(user.Password)
	if err != nil {
		return
	}

	nonceInfo = PasswordNonce{
		ID:       comutils.NewUuid4Code(),
		Value:    comutils.RandomBytesF(PasswordNonceSize),
		Username: username,
		Salt:     passHash.GetSalt(),
	}
	err = passwordSetNonce(nonceInfo)
	return
}

func PasswordValidate(nonceID, username, passwordEncryptedHex string) (user models.User, err error) {
	nonce, err := passwordGetNonce(username, nonceID)
	if err != nil {
		return
	}
	passwordEncrypted, err := comutils.HexDecode(passwordEncryptedHex)
	if err != nil {
		return
	}

	user, err = usermod.GetUserByUsername(username)
	if err != nil {
		return
	}
	passHash, err := LoadPasswordArgon2iHash(user.Password)
	if err != nil {
		return
	}
	password, err := comutils.AesGcmDecrypt(passHash.GetSalt(), nonce, passwordEncrypted)
	if err != nil {
		return
	}
	if !passHash.IsValidPassword(password) {
		err = utils.WrapError(constants.ErrorAuth)
		return
	}

	err = passwordDeleteNonce(username, nonceID)
	return
}
