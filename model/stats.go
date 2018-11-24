package model

// Income keeps information about cash income for every month
type Income struct {
	Year  int     `json:"year" db:"year"`
	Month int     `json:"month" db:"month"`
	Total float64 `json:"total" db:"total"`
	Tax   float64 `json:"tax" db:"tax"`
}
