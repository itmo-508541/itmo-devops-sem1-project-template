package save

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"project_sem/internal/models/price"
	"project_sem/internal/reader"
	"project_sem/internal/server"
)

type Handler struct {
	manager *price.Manager
}

func New(manager *price.Manager) *Handler {
	return &Handler{manager: manager}
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
	if err != nil {
		log.Println(fmt.Errorf("Save.ServeHTTP: %w", err))
		server.JSONBadRequestError(w)
		return
	}

	var accepted price.AcceptedDTO

	accepted, err = h.manager.AcceptCsv(bytes.NewReader(csv))
	if err != nil {
		log.Println(fmt.Errorf("manager.AcceptReader: %w", err))
		server.JSONInternalServerError(w)
		return
	}

	server.JSONResponse(w, accepted, http.StatusOK)
}
