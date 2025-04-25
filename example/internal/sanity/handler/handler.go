package handler

import (
	"net/http"

	"github.com/mickamy/goimportmaps-example/internal/sanity/usecase"
)

type SanityHandler struct {
	usecase *usecase.FindUserUseCase
}

func (h *SanityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.usecase.Do("1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(user.Name))
}
