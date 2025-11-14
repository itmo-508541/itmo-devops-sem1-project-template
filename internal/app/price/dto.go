package price

import (
	"project_sem/internal/app/report"

	"github.com/google/uuid"
)

type PriceRecordDTO struct {
	report.ReportRecordDTO
	GroupUUID uuid.UUID `json:"-"`
}
