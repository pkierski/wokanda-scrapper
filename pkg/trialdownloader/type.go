package trialdownloader

import "time"

type Trial struct {
	CaseID     string    `json:"case_id"`
	Department string    `json:"department"`
	Judges     []string  `json:"judges"`
	Date       time.Time `json:"date"`
	Room       string    `json:"room"`
}
