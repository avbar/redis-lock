package handler

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Locker interface {
	Lock(ctx context.Context, key, token string) (bool, error)
	Unlock(ctx context.Context, key, token string) error
}

type Handler struct {
	locker Locker
	token  string
}

func NewHandler(locker Locker) *Handler {
	h := &Handler{
		locker: locker,
	}
	h.token = h.getToken()
	return h
}

func (h *Handler) getToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (h *Handler) Lock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	ctx := r.Context()
	res, err := h.locker.Lock(ctx, key, h.token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res {
		w.Write([]byte(fmt.Sprintf("%q is locked", key)))
	} else {
		w.Write([]byte(fmt.Sprintf("%q cannot be locked", key)))
	}
}

func (h *Handler) Unlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	ctx := r.Context()
	err := h.locker.Unlock(ctx, key, h.token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%q is unlocked", key)))
}
