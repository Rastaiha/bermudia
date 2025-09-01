package handler

import (
	"errors"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) GetIsland(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	id := chi.URLParam(r, "islandID")
	if id == "" {
		sendError(w, http.StatusBadRequest, "Island ID is required")
		return
	}

	island, err := h.islandService.GetIsland(r.Context(), user.ID, id)
	if err != nil {
		// Check for specific error types using errors.Is
		if errors.Is(err, domain.ErrIslandNotFound) {
			sendError(w, http.StatusNotFound, "Island not found")
			return
		}

		handleError(w, err)
		return
	}

	sendResult(w, island)
}

func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	id := chi.URLParam(r, "inputID")
	if id == "" {
		sendError(w, http.StatusBadRequest, "input ID is required")
		return
	}
	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Error parsing multipart form")
		return
	}

	var file io.ReadCloser
	var filename string
	var textContent string

	if data, ok := r.MultipartForm.File["data"]; ok {
		if len(data) != 1 {
			sendError(w, http.StatusBadRequest, "Incorrect number of files in 'data' field in multipart form")
			return
		}
		if data[0].Size <= 0 {
			sendError(w, http.StatusBadRequest, "Empty file in 'data' field in multipart form")
			return
		}
		f, err := data[0].Open()
		if err != nil {
			handleError(w, err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				slog.Error("Error closing multipart file", err)
			}
		}()
		file = &tempReadCloser{f}
		filename = data[0].Filename
	} else if data, ok := r.MultipartForm.Value["data"]; ok {
		if len(data) == 0 {
			sendError(w, http.StatusBadRequest, "Incorrect number of values in 'data' field in multipart form")
			return
		}
		textContent = strings.Join(data, "\n")
	} else {
		sendError(w, http.StatusBadRequest, "Missing 'data' part in multipart form")
		return
	}

	result, err := h.islandService.SubmitAnswer(r.Context(), user.ID, id, file, filename, textContent)
	if err != nil {
		if errors.Is(err, domain.ErrQuestionNotRelatedToIsland) {
			sendError(w, http.StatusForbidden, "answer not related to player's current island")
			return
		}
		handleError(w, err)
		return
	}
	sendResult(w, result)
}

type tempReadCloser struct {
	r io.ReadCloser
}

func (t *tempReadCloser) Read(p []byte) (n int, err error) {
	return t.r.Read(p)
}

func (t *tempReadCloser) Close() error {
	return t.r.Close()
}
