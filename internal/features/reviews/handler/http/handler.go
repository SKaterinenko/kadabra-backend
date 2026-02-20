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
	router.HandleFunc("DELETE /api/reviews/{id}", handler.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	formData, err := req.ParseMultipartForm(r, 10<<20) // 10MB max
	if err != nil {
		res.Json(w, res.ResDTO{Message: err.Error(), Ok: false}, http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("access_token")
	if utils.CheckErr(w, err) {
		return
	}

	claims, err := utils.ValidateToken(cookie.Value, h.cfg.JWTSecret)
	if utils.CheckErr(w, err) {
		return
	}

	productId, err := req.GetIntFromForm(formData, "product_id")
	if err != nil {
		res.Json(w, res.ResDTO{Message: err.Error(), Ok: false}, http.StatusBadRequest)
		return
	}

	description, exists := formData.FormValues["description"]
	if !exists {
		res.Json(w, res.ResDTO{Message: "description is required", Ok: false}, http.StatusBadRequest)
		return
	}

	rating, err := req.GetIntFromForm(formData, "rating")
	if err != nil {
		res.Json(w, res.ResDTO{Message: err.Error(), Ok: false}, http.StatusBadRequest)
		return
	}

	// Получаем массив изображений
	images, exists := formData.Files["images"]
	//if !exists || len(images) == 0 {
	//	res.Json(w, res.ResDTO{Message: "at least one image is required", Ok: false}, http.StatusBadRequest)
	//	return
	//}

	if len(images) > 3 {
		res.Json(w, res.ResDTO{Message: "maximum 3 images allowed", Ok: false}, http.StatusBadRequest)
		return
	}

	input := &reviews_service.CreateInput{
		UserId:      claims.UserID,
		ProductId:   productId,
		Description: description,
		Rating:      rating,
		Images:      images,
	}

	review, err := h.service.Create(r.Context(), input)
	if utils.CheckErr(w, err) {
		return
	}

	res.Json(w, review, http.StatusOK)
}

func (h *Handler) GetAllById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if utils.CheckErr(w, err) {
		return
	}

	q := r.URL.Query()
	limit, offset, err := utils.GetLimitOffset(q)
	if utils.CheckErr(w, err) {
		return
	}

	review, err := h.service.GetAllById(r.Context(), id, limit, offset)
	if utils.CheckErr(w, err) {
		return
	}

	res.Json(w, review, http.StatusOK)
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

	res.Json(w, res.ResDTO{Message: "Delete successful", Ok: true}, http.StatusOK)
}
