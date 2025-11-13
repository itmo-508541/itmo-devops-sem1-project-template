package price

import (
	"fmt"
	"io"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

type Manager struct {
	processors []PriceProcessor
}

func NewManager() *Manager {
	manager := &Manager{}

	return manager
}

func (m *Manager) AddProcessor(processor PriceProcessor) {
	m.processors = append(m.processors, processor)
}

func (m *Manager) AcceptCsv(r io.Reader) (*Accepted, error) {
	accepted := NewAccepted()

	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := make(chan PriceRecordDTO)
	defer close(inputCh)

	go func(outputCh chan PriceRecordDTO) {
		for _, processor := range m.processors {
			outputCh = m.processStep(doneCh, outputCh, processor)
		}

		for price := range outputCh {
			accepted.Output = append(accepted.Output, price)
		}
	}(inputCh)

	err := gocsv.Unmarshal(r, &accepted.Input)
	if err != nil && false {
		fmt.Println(err)
		return nil, fmt.Errorf("gocsv.Unmarshal: %w", err)
	}
	accepted.TotalCount = len(accepted.Input)
	for _, price := range accepted.Input {
		price.GroupUUID = accepted.UUID
		price.UUID = uuid.New()
		inputCh <- price
	}

	return accepted, nil
}

func (m *Manager) processStep(doneCh chan struct{}, inputCh chan PriceRecordDTO, processor PriceProcessor) chan PriceRecordDTO {
	outputCh := make(chan PriceRecordDTO)

	go func() {
		defer close(outputCh)

		for price := range inputCh {
			err := processor(&price)
			if err != nil {
				continue
			}

			select {
			case <-doneCh:
				return
			case outputCh <- price:
			}
		}
	}()

	return outputCh
}
