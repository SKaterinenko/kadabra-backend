package productsTypeHandler

import (
	"github.com/google/uuid"
	"kadabra/internal/service/productsTypeService"
	"kadabra/pkg/check"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"net/http"
)

type HandlerDeps struct {
	Service *productsTypeService.Service
}

type Handler struct {
	service *productsTypeService.Service
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /products-type", handler.Create)
	router.HandleFunc("GET /products-type", handler.GetAll)
	router.HandleFunc("GET /products-type/{id}", handler.GetById)
	router.HandleFunc("DELETE /products-type/{id}", handler.Delete)
	router.HandleFunc("PATCH /products-type/{id}", handler.Patch)
	router.HandleFunc("GET /products-type-by-category-id/{id}", handler.GetProductsTypeByCategoryId)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](&w, r)
	if err != nil {
		return
	}
	input := &productsTypeService.CreateInput{
		Name:          body.Name,
		SubCategoryId: body.SubCategoryId,
	}
	newProductsType, err := h.service.Create(r.Context(), input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, newProductsType, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	productsType, err := h.service.GetAll(r.Context())
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, productsType, http.StatusOK)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if check.CheckErr(&w, err) {
		return
	}
	productsType, err := h.service.GetById(r.Context(), id)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, productsType, http.StatusOK)
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
	input := &productsTypeService.PatchInput{
		Name: body.Name,
	}
	productsType, err := h.service.Patch(r.Context(), id, input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, productsType, http.StatusOK)
}

func (h *Handler) GetProductsTypeByCategoryId(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if check.CheckErr(&w, err) {
		return
	}
	productsType, err := h.service.GetProductsTypeByCategoryId(r.Context(), id)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, productsType, http.StatusOK)
}
