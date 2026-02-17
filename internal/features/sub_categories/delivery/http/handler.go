package sub_categories_http

import (
	sub_categories_service "kadabra/internal/features/sub_categories/service"
	pkg "kadabra/pkg/lang"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"kadabra/pkg/utils"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	Service SubCategoriesService
}

type Handler struct {
	service SubCategoriesService
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /api/sub-categories", handler.Create)
	router.HandleFunc("GET /api/sub-categories", handler.GetAll)
	router.HandleFunc("GET /api/sub-categories/{id}", handler.GetById)
	router.HandleFunc("DELETE /api/sub-categories/{id}", handler.Delete)
	router.HandleFunc("PATCH /api/sub-categories/{id}", handler.Patch)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](w, r)
	if err != nil {
		return
	}
	input := &sub_categories_service.CreateInput{
		Translations: body.Translations,
		CategoryId:   body.CategoryId,
	}
	newCategory, err := h.service.Create(r.Context(), input)
	if utils.CheckErr(w, err) {
		return
	}
	res.Json(w, newCategory, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	lang := pkg.GetLang(r)
	categories, err := h.service.GetAll(r.Context(), lang)
	if utils.CheckErr(w, err) {
		return
	}
	res.Json(w, categories, http.StatusOK)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	lang := pkg.GetLang(r)
	if utils.CheckErr(w, err) {
		return
	}
	category, err := h.service.GetById(r.Context(), id, lang)
	if utils.CheckErr(w, err) {
		return
	}
	res.Json(w, category, http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if utils.CheckErr(w, err) {
		return
	}
	err = h.service.Delete(r.Context(), id)
	if utils.CheckErr(w, err) {
		return
	}
	res.Json(w, res.ResDTO{Ok: true, Message: "Delete successful"}, http.StatusOK)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[patchDTO](w, r)
	if err != nil {
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if utils.CheckErr(w, err) {
		return
	}
	input := &sub_categories_service.PatchInput{
		Name: body.Name,
	}
	category, err := h.service.Patch(r.Context(), id, input)
	if utils.CheckErr(w, err) {
		return
	}
	res.Json(w, category, http.StatusOK)
}
