package responses

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   string `json:"sum"`
}
