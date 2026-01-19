package products_type_http

import (
	products_type_service "kadabra/internal/features/products_type/service"
	"kadabra/pkg/check"
	pkg "kadabra/pkg/lang"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	Service ProductsTypeService
}

type Handler struct {
	service ProductsTypeService
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /products-type", handler.Create)
	router.HandleFunc("GET /products-type", handler.GetAll)
	router.HandleFunc("GET /products-type/{id}", handler.GetById)
	router.HandleFunc("DELETE /products-type/{id}", handler.Delete)
	router.HandleFunc("PATCH /products-type/{id}", handler.Patch)
	router.HandleFunc("GET /products-type-by-category-slug/{slug}", handler.GetProductsTypeByCategorySlug)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](&w, r)
	if err != nil {
		return
	}

	translations := make([]products_type_service.TranslationInput, len(body.Translations))
	for i, t := range body.Translations {
		translations[i] = products_type_service.TranslationInput{
			LanguageCode: t.LanguageCode,
			Name:         t.Name,
		}
	}

	input := &products_type_service.CreateInput{
		Translations:  translations,
		SubCategoryId: body.SubCategoryId,
	}
	newProductsType, err := h.service.Create(r.Context(), input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, newProductsType, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	lang := pkg.GetLang(r)
	productsType, err := h.service.GetAll(r.Context(), lang)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, productsType, http.StatusOK)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if check.CheckErr(&w, err) {
		return
	}
	lang := pkg.GetLang(r)
	productsType, err := h.service.GetById(r.Context(), id, lang)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, productsType, http.StatusOK)
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
	input := &products_type_service.PatchInput{
		Name: body.Name,
	}
	productsType, err := h.service.Patch(r.Context(), id, input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, productsType, http.StatusOK)
}

func (h *Handler) GetProductsTypeByCategorySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	lang := pkg.GetLang(r)
	productsType, err := h.service.GetProductsTypeByCategorySlug(r.Context(), slug, lang)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, productsType, http.StatusOK)
}
