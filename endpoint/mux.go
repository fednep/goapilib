package endpoint

import (
	"net/http"
)

// HandlerFuncWithData allows to have a unified way to return data and errors
// from handlers.
//
// Any object returned will be serialized as JSON
// Any error returned will be serialized as ApiError
type HandlerFuncWithData func(w http.ResponseWriter, r *http.Request) (any, error)

func (f HandlerFuncWithData) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := f(w, r)
	if err != nil {
		handleError(w, err)
		return
	}

	if data != nil {
		// TODO: think how to handle error from this
		// 	     it seems to be similar with WriteAPIError in the handleError function
		_ = Success(w, data)
	}
}

// HandlerFuncWithError allow to have a unified way to return errors from handlers
// doing so, it becomes easier to centralize handling errors in the implementing service.
// for example to log them to Prometheus, Sentry, or any other logging tool.
type HandlerFuncWithError func(w http.ResponseWriter, r *http.Request) error

func (f HandlerFuncWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)

	if err != nil {
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	apiError, ok := err.(ApiError)
	if !ok {
		// TODO: implement a handler to allow central behavior of these errors,
		// 		 like logging (for example to Sentry, Grafana, or similar)
		// 	     for the initial implementation panic is enough
		panic(err)
	}

	// TODO: consider making it possible to handle error returned by WriteApiError
	//		 however, I'm not sure that any meaningful error can be there, except "Pipe closed"
	// 		 or similar
	_ = WriteApiError(w, apiError)
}
