package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	basePath = "/api"
)

type handler struct {
}

func NewController() http.Handler {
	h := handler{}
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("%v/users/token-issue", basePath), responseHandler(h.login)).Methods(http.MethodPost)
	router.HandleFunc(fmt.Sprintf("%v/users/test-krakend", basePath), responseHandler(h.test)).Methods(http.MethodGet)
	return router
}

func responseHandler(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if status != http.StatusNoContent {
			if err := json.NewEncoder(w).Encode(data); err != nil {
				log.Printf("could not encode response to output: %v", err)
			}
		}
	}
}

func (h *handler) test(w io.Writer, r *http.Request) (interface{}, int, error) {
	return map[string]interface{}{"message": "if you see this message it's because of krakend pass the request through",}, 200, nil
}

func (h *handler) login(w io.Writer, r *http.Request) (interface{}, int, error) {
	fmt.Println("LOGIN CALLED")
	var fields map[string]interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	username, ok := fields["username"].(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("missing username field")
	}

	password, ok := fields["password"].(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("missing password field")
	}

	scopes, err := getScopes(username, password)
	now := time.Now()

	at := map[string]interface{}{
		"aud":    "http://mybackend.com",
		"iss":    "http://mybackend.com",
		"sub":    "someSubject",
		"nbf":    now.Unix(),
		"exp":    now.Add(time.Minute * 15).Unix(),
		"jti":    now,
		"scopes": scopes,
	}

	data := map[string]interface{}{
		"access_token": at,
		"exp":          now.Add(time.Minute * 15).Unix(),
	}

	return data, http.StatusOK, nil
}

func getScopes(username, password string) ([]string, error) {
	mockedPassword := "123456789"
	if password != mockedPassword {
		return nil, errors.New("either username or password are wrong")
	}

	switch username {
	case "all_scopes_user":
		return []string{"inventory", "payment", "order", "other"}, nil
	case "no_scopes_user":
		return []string{}, nil
	case "inventory_scopes_user":
		return []string{"inventory"}, nil
	default:
		return nil, errors.New("either username or password are wrong")
	}
}
