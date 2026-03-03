package categories_http

import (
	categories_service "kadabra/internal/features/categories/service"
	pkg "kadabra/pkg/lang"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"kadabra/pkg/utils"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	Service CategoryService
}

type Handler struct {
	service CategoryService
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /api/categories", handler.Create)
	router.HandleFunc("GET /api/categories", handler.GetAll)
	//router.HandleFunc("GET /api/categories/{id}", handler.GetById)
	router.HandleFunc("GET /api/categories/{slug}", handler.GetBySlug)
	router.HandleFunc("DELETE /api/categories/{id}", handler.Delete)
	router.HandleFunc("PATCH /api/categories/{id}", handler.Patch)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](w, r)
	if err != nil {
		return
	}

	translations := make([]categories_service.TranslationInput, len(body.Translations))
	for i, t := range body.Translations {
		translations[i] = categories_service.TranslationInput{
			LanguageCode: t.LanguageCode,
			Name:         t.Name,
		}
	}

	input := &categories_service.CreateInput{
		Translations: translations,
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

func (h *Handler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	lang := pkg.GetLang(r)

	category, err := h.service.GetBySlug(r.Context(), slug, lang)
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
	res.Json(w, res.ResDTO{Ok: true, Message: "Delete successful!"}, http.StatusOK)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	formData, err := req.ParseMultipartForm(r, 10<<20) // 10MB max
	if utils.CheckErr(w, err) {
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if utils.CheckErr(w, err) {
		return
	}
	imageFile, err := req.GetSingleFileFromForm(formData, "image")
	if utils.CheckErr(w, err) {
		return
	}

	input := &categories_service.PatchInput{
		Image: imageFile,
	}

	category, err := h.service.Patch(r.Context(), id, input)
	if utils.CheckErr(w, err) {
		return
	}
	res.Json(w, category, http.StatusOK)
}
