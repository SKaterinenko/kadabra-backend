package products_http

import (
	products_service "kadabra/internal/features/products/service"
	"kadabra/pkg/check"
	pkg "kadabra/pkg/lang"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	Service ProductsService
}

type Handler struct {
	service ProductsService
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /products", handler.Create)
	router.HandleFunc("GET /products", handler.GetAll)
	router.HandleFunc("GET /products/{id}", handler.GetById)
	router.HandleFunc("DELETE /products/{id}", handler.Delete)
	router.HandleFunc("PATCH /products/{id}", handler.Patch)
	router.HandleFunc("POST /productsByIds", handler.GetByCategoryIds)
	router.HandleFunc("POST /productsByProductsTypeIds", handler.GetByProductsTypeIds)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](&w, r)
	if err != nil {
		return
	}

	translations := make([]products_service.TranslationInput, len(body.Translations))
	for i, t := range body.Translations {
		translations[i] = products_service.TranslationInput{
			LanguageCode:     t.LanguageCode,
			Name:             t.Name,
			ShortDescription: t.ShortDescription,
			Description:      t.Description,
		}
	}

	input := &products_service.CreateInput{
		Translations:   translations,
		ProductTypeId:  body.ProductTypeId,
		ManufacturerId: body.ManufacturerId,
	}
	newProduct, err := h.service.Create(r.Context(), input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, newProduct, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	lang := pkg.GetLang(r)
	products, err := h.service.GetAll(r.Context(), lang)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, products, http.StatusOK)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if check.CheckErr(&w, err) {
		return
	}
	lang := pkg.GetLang(r)
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
	input := &products_service.PatchInput{
		Name:             body.Name,
		Description:      body.Description,
		ShortDescription: body.ShortDescription,
		ProductTypeId:    body.ProductTypeId,
		ManufacturerId:   body.ManufacturerId,
	}
	category, err := h.service.Patch(r.Context(), id, input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, category, http.StatusOK)
}

func (h *Handler) GetByCategoryIds(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[getByIdsDTO](&w, r)
	if err != nil {
		return
	}
	lang := pkg.GetLang(r)
	products, err := h.service.GetByCategoryIds(r.Context(), body.Ids, lang)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, products, http.StatusOK)
}

func (h *Handler) GetByProductsTypeIds(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[getByIdsDTO](&w, r)
	if check.CheckErr(&w, err) {
		return
	}
	lang := pkg.GetLang(r)
	products, err := h.service.GetByProductsTypeIds(r.Context(), body.Ids, lang)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, products, http.StatusOK)
}
