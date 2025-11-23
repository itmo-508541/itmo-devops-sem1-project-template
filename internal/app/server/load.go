package server

import (
	"fmt"
	"log"
	"net/http"
	"project_sem/internal/app/price"
	"project_sem/internal/app/validate"

	"github.com/gocarina/gocsv"
)

func NewLoadHandler(repository DataFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var csv string
		var err error

		filter := price.NewRequestFilter(r)
		fileType := r.URL.Query().Get("type")
		if fileType == "request" {
			JSONResponse(w, filter, http.StatusOK)

			return
		} else {
			v, err := validate.New()
			if err != nil {
				log.Println(fmt.Errorf("validators.NewValidate: %w", err))
				JSONInternalServerError(w)

				return
			}
			err = v.Struct(filter)
			if err != nil {
				log.Println(fmt.Errorf("v.Struct: %w", err))
				JSONBadRequestError(w)

				return
			}
		}

		all, err := repository.FindByFilter(r.Context(), filter)
		if err == nil {
			csv, err = gocsv.MarshalString(all)
		}
		if err != nil {
			log.Println(fmt.Errorf("ServeHTTP: %w", err))
			JSONInternalServerError(w)

			return
		}

		if fileType == "csv" {
			TextResponse(w, csv, http.StatusOK)
		} else {
			ZipResponse(w, "prices.zip", csv, "data.csv", http.StatusOK)
		}
	}
}
