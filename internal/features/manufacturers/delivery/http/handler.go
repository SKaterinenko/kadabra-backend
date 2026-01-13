package manufacturers_http

import (
	manufacturers_service "kadabra/internal/features/manufacturers/service"
	"kadabra/pkg/check"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	Service *manufacturers_service.Service
}

type Handler struct {
	service *manufacturers_service.Service
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service}

	router.HandleFunc("POST /manufacturers", handler.Create)
	router.HandleFunc("GET /manufacturers", handler.GetAll)
	router.HandleFunc("GET /manufacturers/{id}", handler.GetById)
	router.HandleFunc("DELETE /manufacturers/{id}", handler.Delete)
	router.HandleFunc("PATCH /manufacturers/{id}", handler.Patch)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](&w, r)
	if err != nil {
		return
	}
	input := &manufacturers_service.CreateInput{
		Name: body.Name,
	}
	newManufacturer, err := h.service.Create(r.Context(), input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, newManufacturer, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	manufacturers, err := h.service.GetAll(r.Context())
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, manufacturers, http.StatusOK)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if check.CheckErr(&w, err) {
		return
	}
	manufacturer, err := h.service.GetById(r.Context(), id)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, manufacturer, http.StatusOK)
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
	input := &manufacturers_service.PatchInput{
		Name: body.Name,
	}
	manufacturer, err := h.service.Patch(r.Context(), id, input)
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, manufacturer, http.StatusOK)
}
