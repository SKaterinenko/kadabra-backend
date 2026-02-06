package reviews_http

import (
	"kadabra/internal/core/config"
	reviews_service "kadabra/internal/features/reviews/service"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"kadabra/pkg/utils"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	Service ReviewsService
	Cfg     *config.Config
}

type Handler struct {
	service ReviewsService
	cfg     *config.Config
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service, cfg: deps.Cfg}

	router.HandleFunc("POST /api/reviews", handler.Create)
	router.HandleFunc("GET /api/reviews/{id}", handler.GetAllById)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[createDTO](w, r)
	if utils.CheckErr(&w, err) {
		return
	}
	cookie, err := r.Cookie("access_token")
	if utils.CheckErr(&w, err) {
		return
	}

	claims, err := utils.ValidateToken(cookie.Value, h.cfg.JWTSecret)
	if utils.CheckErr(&w, err) {
		return
	}

	input := &reviews_service.CreateInput{
		UserId:      claims.UserID,
		ProductId:   body.ProductId,
		Description: body.Description,
		Rating:      body.Rating,
		Images:      body.Images,
	}

	review, err := h.service.Create(r.Context(), input)
	if utils.CheckErr(&w, err) {
		return
	}

	res.Json(w, review, http.StatusOK)
}

func (h *Handler) GetAllById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	q := r.URL.Query()
	limit, offset, err := utils.GetLimitOffset(q)

	if utils.CheckErr(&w, err) {
		return
	}

	review, err := h.service.GetAllById(r.Context(), id, limit, offset)
	if utils.CheckErr(&w, err) {
		return
	}

	res.Json(w, review, http.StatusOK)
}
