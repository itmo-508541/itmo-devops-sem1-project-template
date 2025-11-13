package save

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"project_sem/internal/models/price"
	"project_sem/internal/models/report"
	"project_sem/internal/reader"
	"project_sem/internal/server"
)

type Handler struct {
	manager    *price.Manager
	priceRepo  *price.Repository
	reportRepo *report.Repository
}

func New(m *price.Manager, p *price.Repository, r *report.Repository) *Handler {
	return &Handler{manager: m, priceRepo: p, reportRepo: r}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var csv []byte

	fileType := r.URL.Query().Get("type")
	file, header, err := r.FormFile("file")
	if err == nil {
		switch fileType {
		case "zip", "":
			zip := reader.ZipArchive{}
			csv, err = zip.ReadContents(file, header.Size)
		case "tar":
			tar := reader.TarArchive{}
			csv, err = tar.ReadContents(file)
		default:
			err = fmt.Errorf("unknown archive type '%s'", fileType)
		}
	}
	if csv[len(csv)-1] != 10 {
		csv = append(csv, 10)
	}

	if err != nil {
		log.Println(fmt.Errorf("Save.ServeHTTP: %w", err))
		server.JSONBadRequestError(w)
		return
	}

	accepted, err := h.manager.AcceptCsv(bytes.NewReader(csv))
	if err != nil {
		log.Println(fmt.Errorf("manager.AcceptCsv: %w", err))
		server.JSONInternalServerError(w)
		return
	}
	err = h.priceRepo.InsertAll(r.Context(), &accepted.Output)
	if err != nil {
		log.Println(fmt.Errorf("priceRepo.InsertAll: %w", err))
		server.JSONInternalServerError(w)
		return
	}
	result, err := h.reportRepo.Renew(r.Context(), accepted.UUID)
	if err != nil {
		log.Println(fmt.Errorf("reportRepo.Renew: %w", err))
		server.JSONInternalServerError(w)
		return
	}

	accepted.DuplicatesCount = result.DuplicatesCount
	accepted.TotalItems = result.TotalItems
	accepted.TotalCategories = result.TotalCategories
	accepted.TotalPrice = result.TotalPrice

	server.JSONResponse(w, *accepted, http.StatusOK)
}
