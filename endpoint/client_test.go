package endpoint

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

const testAddr = "localhost:53312"

type someStruct struct {
	Res string
}

var errInvalidCredentials = ApiError{HttpStatus: 403, ErrId: 10001, Message: "Invalid credentials provided"}

func TestErrorsWrapper(t *testing.T) {

	mux := &http.ServeMux{}
	mux.Handle("/err", HandlerFuncWithError(func(w http.ResponseWriter, r *http.Request) error {
		return errInvalidCredentials
	}))

	mux.Handle("/success", HandlerFuncWithData(func(w http.ResponseWriter, r *http.Request) (any, error) {
		return someStruct{Res: "test-string"}, nil
	}))

	srv := http.Server{Addr: testAddr, Handler: mux}
	go srv.ListenAndServe()

	api := RESTClient{BaseURL: fmt.Sprintf("http://%s/", testAddr)}

	var res = someStruct{}
	err := api.Get(context.Background(), "err", &res)
	if err != nil {
		if !IsError(err, errInvalidCredentials) {
			t.Fatalf("Received incorrect result: %#v", res)
		}
	}

	err = api.Get(context.Background(), "success", &res)
	if err != nil {
		t.Fatalf("Unexpected error received: %s", res)
	}

	if res.Res != "test-string" {
		t.Fatalf("Invalid result received: %s (expected test-string)", res.Res)
	}
	srv.Close()
}

func TestErrors(t *testing.T) {

	succResult := someStruct{Res: "Some-string"}

	mux := &http.ServeMux{}
	mux.HandleFunc("/err", func(w http.ResponseWriter, r *http.Request) {
		WriteApiError(w, errInvalidCredentials)
	})

	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		Success(w, succResult)
	})

	srv := http.Server{Addr: testAddr, Handler: mux}
	go srv.ListenAndServe()

	api := RESTClient{BaseURL: fmt.Sprintf("http://%s/", testAddr)}

	res := someStruct{}

	// Checking for error response
	err := api.Get(context.Background(), "err", &res)
	if err != nil {
		if !IsError(err, errInvalidCredentials) {
			t.Fatalf("Received incorrect result: %#v", res)
		}
	}

	// Checking for success response
	err = api.Get(context.Background(), "success", &res)
	if err != nil {
		t.Fatalf("Expected success result. Got error: %s", err)
	}

	if res.Res != "Some-string" {
		t.Fatalf("Unexpected result. Got: %s, expected: Some-string", res.Res)
	}

	srv.Close()
}
