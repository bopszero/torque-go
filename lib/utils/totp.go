package utils

import "github.com/pquerna/otp/totp"

func IsValidateTOTP(input string, key string) bool {
	return totp.Validate(input, key)
}
