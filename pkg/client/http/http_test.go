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
	var isInternalServerError = false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(map[string]string{
			"data": "data",
		})
		assert.NoError(t, err)

		if !isInternalServerError {
			isInternalServerError = true

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(data)
			assert.NoError(t, err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	body, err := h.GET(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, `{"data":"data"}`, string(body))

	_, err = h.GET(ts.URL)
	assert.Error(t, err)
}
