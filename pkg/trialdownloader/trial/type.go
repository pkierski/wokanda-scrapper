package trial

import "time"

type Trial struct {
	CaseID     string    `json:"case_id"`
	Department string    `json:"department"`
	Judge      string    `json:"judge"`
	Date       time.Time `json:"date"`
	Room       string    `json:"room"`
}
