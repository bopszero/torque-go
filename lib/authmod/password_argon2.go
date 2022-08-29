package authmod

import (
	"crypto/subtle"
	"fmt"
	"strings"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"golang.org/x/crypto/argon2"
)

type PasswordArgon2iConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

var PasswordDefaultArgon2iConfig = PasswordArgon2iConfig{
	Time:    3,
	Memory:  2048,
	Threads: 3,
	KeyLen:  32,
}

type PasswordArgon2iHash struct {
	config PasswordArgon2iConfig
	salt   []byte
	hash   []byte
}

func LoadPasswordArgon2iHash(text string) (hash PasswordArgon2iHash, err error) {
	var (
		parts = strings.Split(text, "$")
		conf  PasswordArgon2iConfig
	)
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &conf.Memory, &conf.Time, &conf.Threads)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	salt, err := comutils.Base64DecodeNoPadding(parts[4])
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	hashBytes, err := comutils.Base64DecodeNoPadding(parts[5])
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	conf.KeyLen = uint32(len(hashBytes))

	hash = PasswordArgon2iHash{
		config: conf,
		salt:   salt,
		hash:   hashBytes,
	}
	return hash, nil
}

func NewPasswordArgon2iHash(password string, c PasswordArgon2iConfig) (hash PasswordArgon2iHash, err error) {
	salt, err := comutils.RandomBytes(PasswordSaltSize)
	if err != nil {
		return
	}

	hash = PasswordArgon2iHash{
		config: c,
		salt:   salt,
		hash:   argon2.Key([]byte(password), salt, c.Time, c.Memory, c.Threads, c.KeyLen),
	}
	return
}

func (this PasswordArgon2iHash) String() string {
	var (
		saltBase64 = comutils.Base64EncodeNoPadding(this.salt)
		hashBase64 = comutils.Base64EncodeNoPadding(this.hash)
	)
	return fmt.Sprintf(
		"$argon2i$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, this.config.Memory, this.config.Time, this.config.Threads,
		saltBase64, hashBase64,
	)
}

func (this PasswordArgon2iHash) GetSalt() []byte {
	return this.salt
}

func (this PasswordArgon2iHash) GetSaltHex() string {
	return comutils.HexEncode(this.salt)
}

func (this PasswordArgon2iHash) IsValidPassword(offerPassword []byte) bool {
	offerHash := argon2.Key(
		offerPassword, this.salt,
		this.config.Time, this.config.Memory, this.config.Threads,
		this.config.KeyLen,
	)
	return subtle.ConstantTimeCompare(this.hash, offerHash) == 1
}
