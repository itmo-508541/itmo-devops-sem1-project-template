package price

import (
	"github.com/google/uuid"

	"project_sem/internal/app/report"
)

type PriceRecordDTO struct {
	report.ReportRecordDTO
	GroupUUID uuid.UUID `json:"-"`
}
