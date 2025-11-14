package report

import "github.com/google/uuid"

type ReportDTO struct {
	Id         string `csv:"id" json:"id" validate:"notblank,number"`
	Name       string `csv:"name" json:"name" validate:"notblank"`
	Category   string `csv:"category" json:"category" validate:"notblank"`
	Price      string `csv:"price" json:"price" validate:"notblank,numeric"`
	CreateDate string `csv:"create_date" json:"create_date" validate:"notblank,date"`
}

type ReportRecordDTO struct {
	UUID uuid.UUID `csv:"-" json:"-"`
	ReportDTO
}
