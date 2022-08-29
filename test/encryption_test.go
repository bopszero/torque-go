package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

const (
	AesGcmSecretHex = "df476ef5bceab612bb31e703ca2594f212384ae15b28e59cbfe67529c6af7a1a"
	AesGcmNonceHex  = "b30603f06720fa1b192ed5f7"
)

func TestAesGcmEncryption(t *testing.T) {
	text := "Hello World!"
	ciptherBytes, err := comutils.AesGcmEncryptString(AesGcmSecretHex, AesGcmNonceHex, text)
	assert.NoError(t, err)
	assert.Equal(t, "firh52aVDyr/4EILDem0SqN01EIe2ejL6w5Tpw==", comutils.Base64Encode(ciptherBytes))

	textBytes, err := comutils.AesGcmDecryptString(AesGcmSecretHex, AesGcmNonceHex, ciptherBytes)
	assert.Equal(t, string(textBytes), text)
}
