package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerCorrectRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/?city=moscow&count=4", nil)
	assert.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code)

	assert.NotEmpty(t, responseRecorder.Body)

	cafeCount := len(strings.Split(responseRecorder.Body.String(), ", "))
	assert.LessOrEqual(t, cafeCount, len(cafeList["moscow"]))
}

// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerInvalidCity(t *testing.T) {
	req, err := http.NewRequest("GET", "/?city=borovichi&count=2", nil)
	assert.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	expectedError := "wrong city value"
	assert.Equal(t, expectedError, responseRecorder.Body.String())
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerCountGreaterThanAvailable(t *testing.T) {
	req, err := http.NewRequest("GET", "/?city=moscow&count=10", nil)
	assert.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code)

	cafeCount := len(strings.Split(responseRecorder.Body.String(), ", "))
	assert.LessOrEqual(t, cafeCount, len(cafeList["moscow"]))
}
