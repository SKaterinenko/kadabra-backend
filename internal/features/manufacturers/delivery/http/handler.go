package manufacturers_http

import (
	manufacturers_service "kadabra/internal/features/manufacturers/service"
	pkg "kadabra/pkg/lang"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"kadabra/pkg/utils"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	Service ManufacturerService
}

type Handler struct {
	service ManufacturerService
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /manufacturers", handler.Create)
	router.HandleFunc("GET /manufacturers", handler.GetAll)
	router.HandleFunc("GET /manufacturers/{id}", handler.GetById)
	router.HandleFunc("DELETE /manufacturers/{id}", handler.Delete)
	router.HandleFunc("PATCH /manufacturers/{id}", handler.Patch)
	router.HandleFunc("GET /manufacturers-by-category-slug/{slug}", handler.GetByCategorySlug)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](&w, r)
	if err != nil {
		return
	}

	translations := make([]manufacturers_service.TranslationInput, len(body.Translations))
	for i, t := range body.Translations {
		translations[i] = manufacturers_service.TranslationInput{
			LanguageCode: t.LanguageCode,
			Description:  t.Description,
		}
	}

	input := &manufacturers_service.CreateInput{
		Name:         body.Name,
		Translations: translations,
		CategoryIds:  body.CategoryIds,
	}
	newManufacturer, err := h.service.Create(r.Context(), input)
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, newManufacturer, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	lang := pkg.GetLang(r)
	manufacturers, err := h.service.GetAll(r.Context(), lang)
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, manufacturers, http.StatusOK)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if utils.CheckErr(&w, err) {
		return
	}
	lang := pkg.GetLang(r)
	manufacturer, err := h.service.GetById(r.Context(), id, lang)
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, manufacturer, http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if utils.CheckErr(&w, err) {
		return
	}
	err = h.service.Delete(r.Context(), id)
	if utils.CheckErr(&w, err) {
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
	if utils.CheckErr(&w, err) {
		return
	}

	var translations *[]manufacturers_service.TranslationInput

	if body.Translations != nil {
		temp := make([]manufacturers_service.TranslationInput, len(*body.Translations))
		for i, t := range *body.Translations {
			temp[i] = manufacturers_service.TranslationInput{
				LanguageCode: t.LanguageCode,
				Description:  t.Description,
			}
		}
		translations = &temp
	}

	input := &manufacturers_service.PatchInput{
		CategoryIds:  body.CategoryIds,
		Translations: translations,
	}

	manufacturer, err := h.service.Patch(r.Context(), id, input)
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, manufacturer, http.StatusOK)
}

func (h *Handler) GetByCategorySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	lang := pkg.GetLang(r)
	manufacturers, err := h.service.GetByCategorySlug(r.Context(), slug, lang)
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, manufacturers, http.StatusOK)
}
