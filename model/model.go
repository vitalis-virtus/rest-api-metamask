package model

type User struct {
	Address string
	Nonce   string
}

type SignInPayload struct {
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
	Sig     string `json:"sig"`
}
