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

func NewSaveHandler(priceRepo *price.Repository, reportRepo *report.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
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
		} else if len(csv) > 0 && csv[len(csv)-1] != 10 {
			csv = append(csv, 10)
		}

		UUID, totalCount, err := priceRepo.AcceptCsv(r.Context(), bytes.NewReader(csv))
		if err != nil {
			err = fmt.Errorf("manager.AcceptCsv: %w", err)
		}

		sr := saveResultDTO{TotalCount: totalCount}
		if err == nil {
			sr.DuplicatesCount, sr.TotalItems, sr.TotalCategories, sr.TotalPrice, err = reportRepo.Renew(
				r.Context(),
				UUID,
			)
			if err != nil {
				err = fmt.Errorf("reportRepo.Renew: %w", err)
			}
		}

		if err != nil {
			log.Println(err)
			JSONInternalServerError(w)
		} else {
			JSONResponse(w, sr, http.StatusOK)
		}
	}
}
