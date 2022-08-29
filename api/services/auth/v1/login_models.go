package v1

import "gitlab.com/snap-clickstaff/torque-go/lib/authmod"

type LoginInputPrepareRequest struct {
	Username string `json:"username" validate:"required,max=64"`
}

type LoginInputPrepareResponse struct {
	NonceID string `json:"nonce_id"`
	Nonce   string `json:"nonce"`
	Salt    string `json:"salt"`
}

type LoginInputExecuteRequest struct {
	Username             string `json:"username" validate:"required,max=64"`
	PasswordEncryptedHex string `json:"password" validate:"required,min=8"`
	NonceID              string `json:"nonce_id" validate:"required,max=32"`
}

type LoginInputExecuteResponse struct {
	CommitToken string `json:"commit_token"`
}

type LoginInputCommitRequest struct {
	DeviceUID string `json:"device_uid" validate:"required,max=128"`
	AuthCode  string `json:"auth_code" validate:"max=6"`
}

type LoginInputCommitResponse struct {
	authmod.JwtKeyPair
}

type LoginRefreshRequest struct {
	DeviceUID string `json:"device_uid" validate:"required,max=128"`
	Rotate    bool   `json:"rotate"`
}

type LoginRefreshResponse struct {
	authmod.JwtKeyPair
}
