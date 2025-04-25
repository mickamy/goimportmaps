package handler

import (
	"net/http"

	"github.com/mickamy/goimportmaps-example/internal/insanity/repository"
)

type InsanityHandler struct {
	repo repository.User
}

func (h *InsanityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.repo.Find("1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(user.Name))
}
