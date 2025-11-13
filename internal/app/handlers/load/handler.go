package load

import (
	_ "archive/zip"
	"fmt"
	"log"
	"net/http"
	"project_sem/internal/models/report"
	"project_sem/internal/server"

	"github.com/gocarina/gocsv"
)

type Handler struct {
	reportRepo *report.Repository
}

func New(r *report.Repository) *Handler {
	return &Handler{reportRepo: r}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var csv string

	fileType := r.URL.Query().Get("type")

	all, err := h.reportRepo.All(r.Context())
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
