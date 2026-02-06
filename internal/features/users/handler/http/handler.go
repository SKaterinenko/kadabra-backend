package users_http

import (
	"kadabra/internal/core/config"
	users_service "kadabra/internal/features/users/service"
	"kadabra/pkg/req"
	"kadabra/pkg/res"
	"kadabra/pkg/utils"
	"net/http"
)

type HandlerDeps struct {
	Service UsersService
	Cfg     *config.Config
}

type Handler struct {
	service UsersService
	cfg     *config.Config
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{service: deps.Service, cfg: deps.Cfg}

	router.HandleFunc("POST /api/auth/register", handler.Register)
	router.HandleFunc("POST /api/auth/login", handler.Login)
	router.HandleFunc("POST /api/auth/logout", handler.Logout)
	router.HandleFunc("GET /api/auth/refresh", handler.RefreshTokens)
	router.HandleFunc("GET /api/me", handler.Me)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[registerDTO](w, r)
	if err != nil {
		return
	}

	createUser := &users_service.CreateUserRequest{
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		BirthDate:   body.BirthDate,
		PhoneNumber: body.PhoneNumber,
		Gender:      body.Gender,
		Password:    body.Password,
	}

	authResp, err := h.service.Register(r.Context(), createUser)
	if utils.CheckErr(&w, err) {
		return
	}
	utils.SetAuthCookies(w, authResp.AccessToken, authResp.RefreshToken, h.cfg)
	res.Json(w, authResp.User, http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[loginDTO](w, r)
	if err != nil {
		return
	}

	authResp, err := h.service.Login(r.Context(), body.Email, body.Password)
	if utils.CheckErr(&w, err) {
		return
	}
	utils.SetAuthCookies(w, authResp.AccessToken, authResp.RefreshToken, h.cfg)
	res.Json(w, authResp.User, http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.ClearAuthCookies(w)
	res.Json(w, res.ResDTO{Message: "Successfully logged out", Ok: true}, http.StatusOK)
}

func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.ClearAuthCookies(w)
		res.Json(w, res.ResDTO{
			Message: "Refresh token not found",
			Ok:      false,
		}, http.StatusUnauthorized)
		return
	}

	authResp, err := h.service.RefreshTokens(r.Context(), cookie.Value)
	if utils.CheckErr(&w, err) {
		utils.ClearAuthCookies(w) // Очищаем невалидные токены
		return
	}

	// Устанавливаем новые токены
	utils.SetAuthCookies(w, authResp.AccessToken, authResp.RefreshToken, h.cfg)

	// Отдаём данные пользователя
	res.Json(w, authResp.User, http.StatusOK)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	// Получаем access token из куки
	cookie, err := r.Cookie("access_token")
	if err != nil {
		res.Json(w, res.ResDTO{Message: "Unauthorized", Ok: false}, http.StatusUnauthorized)
		return
	}

	// Валидируем и парсим токен
	claims, err := utils.ValidateToken(cookie.Value, h.cfg.JWTSecret)
	if utils.CheckErr(&w, err) {
		return
	}

	// Получаем пользователя из БД
	user, err := h.service.GetByID(r.Context(), claims.UserID)
	if utils.CheckErr(&w, err) {
		return
	}

	res.Json(w, user, http.StatusOK)
}
