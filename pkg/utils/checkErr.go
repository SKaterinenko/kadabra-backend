package utils

import (
	"kadabra/pkg/res"
	"net/http"
)

func CheckErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		res.Json(w, res.ResDTO{
			Message: err.Error(),
			Ok:      false,
		}, http.StatusBadRequest)
		return true
	}
	return false
}
