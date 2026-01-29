package products_http

import (
	products_service "kadabra/internal/features/products/service"
	pkg "kadabra/pkg/lang"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"kadabra/pkg/utils"
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
	router.HandleFunc("POST /products-by-category-ids", handler.GetByCategoryIds)
	router.HandleFunc("POST /products-by-products-type-ids", handler.GetByProductsTypeIds)
	router.HandleFunc("GET /products-by-category-slug/{slug}", handler.GetByCategorySlug)
	router.HandleFunc("GET /products-by-manufacturer-id/{id}", handler.GetByManufacturersIds)
	router.HandleFunc("GET /product/{slug}", handler.GetProductBySlug)
	router.HandleFunc("POST /product-variations", handler.CreateProductVariations)
	router.HandleFunc("DELETE /product-variations/{id}", handler.DeleteProductVariations)
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
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, newProduct, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	categories, err := utils.ParseIntSlice(q.Get("categories"))
	if utils.CheckErr(&w, err) {
		return
	}

	types, err := utils.ParseIntSlice(q.Get("types"))
	if utils.CheckErr(&w, err) {
		return
	}

	manufacturers, err := utils.ParseIntSlice(q.Get("manufacturers"))
	if utils.CheckErr(&w, err) {
		return
	}

	limit, offset, err := utils.GetLimitOffset(q)
	if utils.CheckErr(&w, err) {
		return
	}

	lang := pkg.GetLang(r)
	products, err := h.service.GetAll(r.Context(), lang, categories, types, manufacturers, limit, offset)
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, products, http.StatusOK)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if utils.CheckErr(&w, err) {
		return
	}
	lang := pkg.GetLang(r)
	category, err := h.service.GetById(r.Context(), id, lang)
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, category, http.StatusOK)
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
	input := &products_service.PatchInput{
		Name:             body.Name,
		Description:      body.Description,
		ShortDescription: body.ShortDescription,
		ProductTypeId:    body.ProductTypeId,
		ManufacturerId:   body.ManufacturerId,
	}
	category, err := h.service.Patch(r.Context(), id, input)
	if utils.CheckErr(&w, err) {
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
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, products, http.StatusOK)
}

func (h *Handler) GetByProductsTypeIds(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[getByIdsDTO](&w, r)
	if err != nil {
		return
	}
	lang := pkg.GetLang(r)
	products, err := h.service.GetByProductsTypeIds(r.Context(), body.Ids, lang)
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, products, http.StatusOK)
}

func (h *Handler) GetByCategorySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	lang := pkg.GetLang(r)
	products, err := h.service.GetByCategorySlug(r.Context(), lang, slug)
	if utils.CheckErr(&w, err) {
		return
	}

	res.Json(w, products, http.StatusOK)
}

func (h *Handler) GetByManufacturersIds(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[getByIdsDTO](&w, r)
	if err != nil {
		return
	}
	lang := pkg.GetLang(r)
	products, err := h.service.GetByManufacturersIds(r.Context(), body.Ids, lang)
	if utils.CheckErr(&w, err) {
		return
	}

	res.Json(w, products, http.StatusOK)
}

func (h *Handler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	lang := pkg.GetLang(r)
	product, err := h.service.GetBySlug(r.Context(), slug, lang)
	if utils.CheckErr(&w, err) {
		return
	}

	res.Json(w, product, http.StatusOK)
}

func (h *Handler) CreateProductVariations(w http.ResponseWriter, r *http.Request) {
	formData, err := req.ParseMultipartForm(r, 10<<20) // 10MB max
	if err != nil {
		res.Json(w, res.ResDTO{Message: err.Error(), Ok: false}, http.StatusBadRequest)
		return
	}

	productId, err := req.GetIntFromForm(formData, "product_id")
	if err != nil {
		res.Json(w, res.ResDTO{Message: err.Error(), Ok: false}, http.StatusBadRequest)
		return
	}

	price, err := req.GetDecimalFromForm(formData, "price")
	if err != nil {
		res.Json(w, res.ResDTO{Message: err.Error(), Ok: false}, http.StatusBadRequest)
		return
	}

	imageFile, exists := formData.Files["image"]
	if !exists {
		res.Json(w, res.ResDTO{Message: "image file is required", Ok: false}, http.StatusBadRequest)
		return
	}

	variation := &products_service.VariationInput{
		ProductId: productId,
		Image:     imageFile,
		Price:     price,
	}

	product, err := h.service.CreateProductVariations(r.Context(), variation)
	if utils.CheckErr(&w, err) {
		return
	}

	res.Json(w, product, http.StatusOK)
}

func (h *Handler) DeleteProductVariations(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if utils.CheckErr(&w, err) {
		return
	}
	err = h.service.DeleteProductVariation(r.Context(), id)
	if utils.CheckErr(&w, err) {
		return
	}
	res.Json(w, res.ResDTO{Ok: true, Message: "Delete successful"}, http.StatusOK)
}
