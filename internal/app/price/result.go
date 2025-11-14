package price

import (
	"project_sem/internal/app/report"

	"github.com/google/uuid"
)

type Accepted struct {
	TotalCount int `json:"total_count"`
	report.Accepted
	UUID   uuid.UUID        `json:"-"`
	Input  []PriceRecordDTO `json:"-"`
	Output []PriceRecordDTO `json:"-"`
}

func NewAccepted() *Accepted {
	accepted := &Accepted{UUID: uuid.New()}
	accepted.Input = make([]PriceRecordDTO, 0)
	accepted.Output = make([]PriceRecordDTO, 0)

	return accepted
}
