package sub_categories_http

import (
	sub_categories_service "kadabra/internal/features/sub_categories/service"
	"kadabra/pkg"
	"kadabra/pkg/check"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	Service *sub_categories_service.Service
}

type Handler struct {
	service *sub_categories_service.Service
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /sub-categories", handler.Create)
	router.HandleFunc("GET /sub-categories", handler.GetAll)
	router.HandleFunc("GET /sub-categories/{id}", handler.GetById)
	router.HandleFunc("DELETE /sub-categories/{id}", handler.Delete)
	router.HandleFunc("PATCH /sub-categories/{id}", handler.Patch)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](&w, r)
	if err != nil {
		return
	}
	input := &sub_categories_service.CreateInput{
		Translations: body.Translations,
		CategoryId:   body.CategoryId,
	}
	newCategory, err := h.service.Create(r.Context(), input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, newCategory, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	lang := pkg.GetLang(r)
	categories, err := h.service.GetAll(r.Context(), lang)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, categories, http.StatusOK)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	lang := pkg.GetLang(r)
	if check.CheckErr(&w, err) {
		return
	}
	category, err := h.service.GetById(r.Context(), id, lang)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, category, http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
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
	id, err := strconv.Atoi(idStr)
	if check.CheckErr(&w, err) {
		return
	}
	input := &sub_categories_service.PatchInput{
		Name: body.Name,
	}
	category, err := h.service.Patch(r.Context(), id, input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, category, http.StatusOK)
}
