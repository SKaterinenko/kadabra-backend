package productHandler

import (
	"github.com/google/uuid"
	"kadabra/internal/service/productService"
	"kadabra/pkg/check"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"net/http"
)

type HandlerDeps struct {
	Service *productService.Service
}

type Handler struct {
	service *productService.Service
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /products", handler.Create)
	router.HandleFunc("GET /products", handler.GetAll)
	router.HandleFunc("GET /products/{id}", handler.GetById)
	router.HandleFunc("DELETE /products/{id}", handler.Delete)
	router.HandleFunc("PATCH /products/{id}", handler.Patch)
	router.HandleFunc("POST /productsByIds", handler.GetByCategoryIds)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](&w, r)
	if err != nil {
		return
	}
	input := &productService.CreateInput{
		Name:             body.Name,
		Description:      body.Description,
		ShortDescription: body.ShortDescription,
		ProductsTypeId:   body.ProductsTypeId,
		ManufacturerId:   body.ManufacturerId,
	}
	newProduct, err := h.service.Create(r.Context(), input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, newProduct, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll(r.Context())
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, products, http.StatusOK)
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
	input := &productService.PatchInput{
		Name: body.Name,
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
	products, err := h.service.GetByCategoryIds(r.Context(), body.Ids)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, products, http.StatusOK)
}
