package utils

import (
	"errors"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestCheckErr_WithError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("boom")
	ok := CheckErr(w, err)

	require.True(t, ok)
	require.Equal(t, 400, w.Code)
	require.Contains(t, w.Body.String(), "boom")
}

func TestCheckErr_NoError(t *testing.T) {
	w := httptest.NewRecorder()
	ok := CheckErr(w, nil)

	require.False(t, ok)
	require.Equal(t, 200, w.Code)
}
