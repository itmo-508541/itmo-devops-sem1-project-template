package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"project_sem/internal/app/price"
	"project_sem/internal/app/report"
	"project_sem/internal/reader"
)

func NewSaveHandler(manager *price.Manager, priceRepo *price.Repository, reportRepo *report.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
				err = fmt.Errorf("Request.Get: unknown archive type '%s'", fileType)
			}
		}
		if err != nil {
			log.Println(fmt.Errorf("Save.ServeHTTP: %w", err))
			JSONBadRequestError(w)
			return
		}

		if csv[len(csv)-1] != 10 {
			csv = append(csv, 10)
		}
		accepted, err := manager.AcceptCsv(bytes.NewReader(csv))

		if err != nil {
			err = fmt.Errorf("manager.AcceptCsv: %w", err)
		}
		if err != nil {
			err = priceRepo.InsertAll(r.Context(), &accepted.Output)
			if err != nil {
				err = fmt.Errorf("priceRepo.InsertAll: %w", err)
			}
		}

		var result *report.Accepted
		if err != nil {
			result, err = reportRepo.Renew(r.Context(), accepted.UUID)
			if err != nil {
				err = fmt.Errorf("reportRepo.Renew: %w", err)
			}
		}

		if err != nil {
			log.Println(err)

			JSONInternalServerError(w)
		} else {
			accepted.DuplicatesCount = result.DuplicatesCount
			accepted.TotalItems = result.TotalItems
			accepted.TotalCategories = result.TotalCategories
			accepted.TotalPrice = result.TotalPrice

			JSONResponse(w, *accepted, http.StatusOK)
		}
	}
}
