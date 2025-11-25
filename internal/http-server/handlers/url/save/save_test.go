package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shorter/internal/http-server/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler_Success(t *testing.T) {
	mockSaver := mocks.NewURLSaver(t)
	mockGetter := mocks.NewURLGetter(t)

	mockGetter.On("GetAllAliases").Return([]string{}, nil)
	mockSaver.On("SaveURL", "https://google.com/long", "https://google.com/short").Return(int64(1), nil)

	handler := New(slog.Default(), mockSaver, mockGetter)

	body := `{"url": "https://google.com/long", "alias":"https://google.com/short"}`
	req := httptest.NewRequest("POST", "/save", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "OK", response.Status)
	assert.Equal(t, "https://google.com/short", response.Alias)
}

func TestSaveHandler_Error(t *testing.T) {
	mockSaver := mocks.NewURLSaver(t)
	mockGetter := mocks.NewURLGetter(t)

	mockGetter.On("GetAllAliases").Return([]string{"https://google.com/short"}, nil)
	// Возвращаем ошибку
	mockSaver.On("SaveURL", "https://google.com/long", "https://google.com/short").
		Return(int64(0), errors.New("db fail"))

	handler := New(slog.Default(), mockSaver, mockGetter)

	body := `{"url": "https://google.com/long", "alias":"https://google.com/short"}`
	req := httptest.NewRequest("POST", "/save", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Error", response.Status)
	assert.Equal(t, "failed to save url", response.Error)
}
