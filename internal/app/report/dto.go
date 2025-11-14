package report

import "github.com/google/uuid"

type ReportDTO struct {
	Id         int     `csv:"id" json:"id" validate:"nonzero"`
	Name       string  `csv:"name" json:"name" validate:"nonzero"`
	Category   string  `csv:"category" json:"category" validate:"nonzero"`
	Price      float32 `csv:"price" json:"price" validate:"nonzero"`
	CreateDate string  `csv:"create_date" json:"create_date" validate:"nonzero,regexp=^(((19|20)([2468][048]|[13579][26]|0[48])|2000)[/-]02[/-]29|((19|20)[0-9]{2}[/-](0[4678]|1[02])[/-](0[1-9]|[12][0-9]|30)|(19|20)[0-9]{2}[/-](0[1359]|11)[/-](0[1-9]|[12][0-9]|3[01])|(19|20)[0-9]{2}[/-]02[/-](0[1-9]|1[0-9]|2[0-8])))$"`
}

type ReportRecordDTO struct {
	UUID uuid.UUID `csv:"-" json:"-"`
	ReportDTO
}
