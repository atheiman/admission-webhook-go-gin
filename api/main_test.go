package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestAdmissionReviewSuccessRoute(t *testing.T) {
	router := setupRouter()

	var jsonStr = []byte(`{
		"uid":"1234",
		"kind": {"group": "", "version": "v1", "kind": "Pod"}
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/admission-review", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"uid\":\"1234\",\"allowed\":true}", w.Body.String())
}

func TestJsonBindingSuccessRoute(t *testing.T) {
	router := setupRouter()

	var jsonStr = []byte(`{"name":"austin", "age": 25}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/json-binding", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"welcome\":\"austin\"}", w.Body.String())
}

func TestJsonBindingDeniedRoute(t *testing.T) {
	router := setupRouter()

	var jsonStr = []byte(`{"name":"script-kiddie", "age": 14}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/json-binding", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"status\":\"unauthorized\"}", w.Body.String())
}

func TestJsonBindingBadRequestRoute(t *testing.T) {
	router := setupRouter()

	var jsonStr = []byte(`{"name":"bad-user", "age": "a string"}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/json-binding", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"json: cannot unmarshal string into Go struct field Person.age of type int\"}", w.Body.String())
}

func TestJsonBindingMissingRequiredRoute(t *testing.T) {
	router := setupRouter()

	var jsonStr = []byte(`{"age": 30}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/json-binding", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"Key: 'Person.Name' Error:Field validation for 'Name' failed on the 'required' tag\"}", w.Body.String())
}
