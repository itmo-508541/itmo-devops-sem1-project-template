package server

import (
	"context"
	"io"
	"project_sem/internal/app/price"
)

type CsvReceiver interface {
	AcceptCsv(context.Context, io.Reader) (*price.AcceptResultDTO, error)
}

type DataFinder interface {
	FindByFilter(context.Context, price.SqlFilter) (*[]price.PriceDTO, error)
}
