package deribit

type Transaction struct {
	SessionRpl float64 `json:"session_rpl"`
	Timestamp  int64   `json:"timestamp"`
	Price      float64 `json:"price"`
}

type transactionsResponse struct {
	Result struct {
		Logs []Transaction `json:"logs"`
	} `json:"result"`
}
