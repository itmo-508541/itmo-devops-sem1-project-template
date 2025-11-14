package price

import (
	"gopkg.in/validator.v2"
)

type PriceProcessor func(price *PriceRecordDTO) error

func NewValidateProcessor() PriceProcessor {
	return func(price *PriceRecordDTO) error {
		return validator.Validate(price)
	}
}
