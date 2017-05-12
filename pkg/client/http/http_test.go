package http_test

import (
	"encoding/json"
	h "github.com/lavrs/dms/pkg/client/http"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGET(t *testing.T) {
	ts200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(map[string]string{
			"data": "data",
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		assert.NoError(t, err)
	}))
	defer ts200.Close()

	ts500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts500.Close()

	body, err := h.GET(ts200.URL)
	assert.NoError(t, err)
	assert.Equal(t, `{"data":"data"}`, string(body))

	_, err = h.GET(ts500.URL)
	assert.Error(t, err)
}
