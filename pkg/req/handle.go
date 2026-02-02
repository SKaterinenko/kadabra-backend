package req

import (
	"kadabra/pkg/res"
	"net/http"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		res.Json(w, res.ResDTO{Message: err.Error(), Ok: false}, http.StatusInternalServerError)
		return nil, err
	}

	err = IsValid(body)
	if err != nil {
		res.Json(w, res.ResDTO{Message: err.Error(), Ok: false}, http.StatusBadRequest)
		return nil, err
	}

	return &body, nil
}
