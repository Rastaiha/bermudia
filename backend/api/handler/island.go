package handler

import (
	"errors"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
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
	data, ok := r.MultipartForm.File["data"]
	if !ok {
		sendError(w, http.StatusBadRequest, "Missing 'data' part in multipart form")
		return
	}
	if len(data) != 1 {
		sendError(w, http.StatusBadRequest, "Incorrect number of files in 'data' field in multipart form")
		return
	}
	filename := data[0].Filename
	file, err := data[0].Open()
	if err != nil {
		handleError(w, err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Error closing multipart file", err)
		}
	}()

	result, err := h.islandService.SubmitAnswer(r.Context(), user.ID, id, &tempReadCloser{file}, filename)
	if err != nil {
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
