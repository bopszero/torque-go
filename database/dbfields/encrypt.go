package dbfields

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"gitlab.com/snap-clickstaff/go-common/comutils"
)

const (
	EncryptedPrefix      = "!"
	EncryptedSeparator   = ":"
	EncryptedNonceLength = 12
)

type EncryptedField struct {
	data   string
	secret []byte
}

func NewEncryptedField(value string) EncryptedField {
	return EncryptedField{data: value}
}

func (this *EncryptedField) GetValue() (string, error) {
	if !strings.HasPrefix(this.data, EncryptedPrefix) {
		return this.data, nil
	}

	encryptedVal := this.data[len(EncryptedPrefix):]
	decryptedText, err := this.Decrypt(encryptedVal)
	if err != nil {
		return "", err
	}

	return decryptedText, nil
}

func (this EncryptedField) Value() (driver.Value, error) {
	if strings.HasPrefix(this.data, EncryptedPrefix) {
		return this.data, nil
	}

	encryptedText, err := this.Encrypt(this.data)
	if err != nil {
		return nil, err
	}

	return EncryptedPrefix + encryptedText, nil
}

func (this *EncryptedField) Scan(input interface{}) (err error) {
	this.data = string(input.([]byte))
	return nil
}

func (this EncryptedField) MarshalText() ([]byte, error) {
	return this.MarshalBinary()
}

func (this *EncryptedField) UnmarshalText(text []byte) error {
	return this.UnmarshalBinary(text)
}

func (this EncryptedField) MarshalBinary() ([]byte, error) {
	return []byte(this.data), nil
}

func (this *EncryptedField) UnmarshalBinary(data []byte) error {
	this.data = string(data)
	return nil
}

func (this *EncryptedField) SetSecretHex(secretHex string) error {
	secret, err := comutils.HexDecode(secretHex)
	if err != nil {
		return err
	}
	if len(secret) != 32 {
		return fmt.Errorf("`EncryptedField` requires 32-length secret to use AES-256")
	}

	this.secret = secret
	return nil
}

func (this *EncryptedField) validateStateEncryption() error {
	if this.secret == nil {
		return fmt.Errorf("`EncryptedField` requires `secretHex` before execution")
	}
	return nil
}

func (this *EncryptedField) Encrypt(value string) (_ string, err error) {
	if err = this.validateStateEncryption(); err != nil {
		return
	}

	nonce, err := comutils.RandomBytes(EncryptedNonceLength)
	if err != nil {
		return
	}
	cipherBytes, err := comutils.AesGcmEncrypt(this.secret, nonce, []byte(value))
	if err != nil {
		return
	}

	nonceBase64 := comutils.Base64Encode(nonce)
	cipherBase64 := comutils.Base64Encode(cipherBytes)
	encryptedText := nonceBase64 + EncryptedSeparator + cipherBase64

	return encryptedText, nil
}

func (this *EncryptedField) Decrypt(value string) (_ string, err error) {
	if err = this.validateStateEncryption(); err != nil {
		return
	}

	encryptedBase64Parts := strings.Split(value, EncryptedSeparator)
	if len(encryptedBase64Parts) != 2 {
		err = fmt.Errorf(
			"`EncryptedField` requires 2 parts in encrypted value (not %v)",
			len(encryptedBase64Parts))
		return
	}
	nonceBytes, err := comutils.Base64Decode(encryptedBase64Parts[0])
	if err != nil {
		return
	}
	cipherBytes, err := comutils.Base64Decode(encryptedBase64Parts[1])
	if err != nil {
		return
	}
	decryptedBytes, err := comutils.AesGcmDecrypt(this.secret, nonceBytes, cipherBytes)
	if err != nil {
		return
	}
	return string(decryptedBytes), nil
}
