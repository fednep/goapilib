package endpoint

import (
	"encoding/json"
	"net/http"
)

// Usage pattern
//
// func (...) handlePage(w ResponseWriter, r *http.Request) (any, error) {
//      obj := User{}
//      err := endpoint.In(&obj)
//      ...
//      return outObj, nil
// }
func Success(w http.ResponseWriter, obj any) error {
	w.Header().Add("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	return encoder.Encode(obj)
}

func WriteApiError(w http.ResponseWriter, err ApiError) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.HttpStatus)

	encoder := json.NewEncoder(w)
	return encoder.Encode(err)
}
