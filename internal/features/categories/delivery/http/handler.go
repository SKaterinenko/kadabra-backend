package categories_http

import (
	"github.com/google/uuid"
	category_service "kadabra/internal/features/categories/service"
	"kadabra/pkg/check"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"net/http"
)

type HandlerDeps struct {
	Service CategoryService
}

type Handler struct {
	service CategoryService
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /categories", handler.Create)
	router.HandleFunc("GET /categories", handler.GetAll)
	router.HandleFunc("GET /categories/{id}", handler.GetById)
	router.HandleFunc("DELETE /categories/{id}", handler.Delete)
	router.HandleFunc("PATCH /categories/{id}", handler.Patch)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](&w, r)
	if err != nil {
		return
	}
	input := &category_service.CreateInput{
		Name: body.Name,
	}
	newCategory, err := h.service.Create(r.Context(), input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, newCategory, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll(r.Context())
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, categories, http.StatusOK)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if check.CheckErr(&w, err) {
		return
	}
	category, err := h.service.GetById(r.Context(), id)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, category, http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if check.CheckErr(&w, err) {
		return
	}
	err = h.service.Delete(r.Context(), id)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, res.ResDTO{Ok: true, Message: "Delete successful"}, http.StatusOK)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[patchDTO](&w, r)
	if err != nil {
		return
	}
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if check.CheckErr(&w, err) {
		return
	}
	input := &category_service.PatchInput{
		Name: body.Name,
	}
	category, err := h.service.Patch(r.Context(), id, input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, category, http.StatusOK)
}
