package handlers

import (
	"fmt"
	"log"
	"net/http"
	"project_sem/internal/models/report"
	"project_sem/internal/server"

	"github.com/gocarina/gocsv"
)

func NewLoadHandler(reportRepo *report.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var csv string

		fileType := r.URL.Query().Get("type")

		all, err := reportRepo.All(r.Context())
		if err == nil {
			csv, err = gocsv.MarshalString(all)
		}

		if err != nil {
			log.Println(fmt.Errorf("ServeHTTP: %w", err))
			server.JSONInternalServerError(w)
			return
		}

		if fileType == "csv" {
			server.TextResponse(w, csv, http.StatusOK)
		} else {
			server.ZipResponse(w, "prices.zip", csv, "data.csv", http.StatusOK)
		}
	}
}
