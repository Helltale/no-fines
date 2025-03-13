package domain

type Transaction struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	From       string  `json:"from_currency"`
	To         string  `json:"to_currency"`
	Amount     float64 `json:"amount"`
	Rate       float64 `json:"rate"`
	Commission float64 `json:"commission"`
	Timestamp  int64   `json:"timestamp"`
	Status     string  `json:"status"` // e.g., "completed", "failed"
}
