package responses

type Response struct {
	Result string `json:"result"`
}

type Order struct {
	Number     int    `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}
