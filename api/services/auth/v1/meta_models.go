package v1

type BannedInfo struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Message string `json:"message"`
}

type Registration struct {
	BannedInfo BannedInfo `json:"banned_info"`
	IsBanned   bool       `json:"is_banned"`
}

type MetaStatusResponse struct {
	Registration Registration `json:"registration"`
}
