package miner

import (
	"kadabra/internal/service/miner"
	"kadabra/pkg/check"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"net/http"
)

type MinerHandlerDeps struct {
	Service *service.MinerService
}

type MinerHandler struct {
	service *service.MinerService
}

func NewCoalHandler(router *http.ServeMux, deps MinerHandlerDeps) {
	handler := &MinerHandler{
		service: deps.Service,
	}
	router.HandleFunc("POST /miners", handler.Create)
	router.HandleFunc("GET /miners", handler.GetAll)
}

func (h *MinerHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[CreateMinerDTO](&w, r)
	if err != nil {
		return
	}
	input := &service.CreateMinerInput{Name: body.Name, Energy: body.Energy, Age: body.Age}
	miner, err := h.service.CreateMiner(r.Context(), input)
	if check.CheckErr(&w, err) {
		return
	} else {
		res.Json(w, miner, http.StatusOK)
		return
	}
}

func (h *MinerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	miners, err := h.service.GetAll(r.Context())
	if check.CheckErr(&w, err) {
		return
	}
	res.Json(w, miners, http.StatusOK)
}
