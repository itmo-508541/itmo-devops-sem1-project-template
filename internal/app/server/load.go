package server

import (
	"fmt"
	"log"
	"net/http"
	"project_sem/internal/app/report"

	"github.com/go-playground/validator/v10"
	"github.com/gocarina/gocsv"
)

// NewLoadHandler возвращает GET handler
// http://localhost:8080/api/v0/prices?type=csv&start=2023-01-01&end=2025-10-01&min=10&max=20
func NewLoadHandler(reportRepo *report.Repository, v *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var csv string
		var err error

		filter := report.NewRequestFilter(r)
		fileType := r.URL.Query().Get("type")
		if fileType == "request" {
			JSONResponse(w, filter, http.StatusOK)

			return
		} else {
			err = v.Struct(filter)
			if err != nil {
				log.Println(err)
				JSONBadRequestError(w)

				return
			}
		}

		all, err := reportRepo.All(r.Context(), filter)
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
