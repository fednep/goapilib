package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// RESTClient helps to work with REST APIs which consumes and produces JSONs
type RESTClient struct {
	httpClient http.Client

	// Function which will construct errors can be overriden
	ErrParser func(*http.Response, []byte) error

	// baseURL should not contain trailing "/"
	BaseURL string
}

type ApiError struct {
	HttpStatus int

	// ErrId represents a unique error id from specific service
	ErrId   int    `json:"id"`
	Message string `json:"message"`
}

func IsError(err error, e ApiError) bool {
	aerr, ok := err.(ApiError)
	if !ok {
		return false
	}

	return aerr.ErrId == e.ErrId && aerr.HttpStatus == e.HttpStatus
}

func (e ApiError) Error() string {
	return fmt.Sprintf("error (id: %d): %s", e.ErrId, e.Message)
}

type HttpError struct {
	HttpStatus int

	// The content of the whole response body
	ResponseBody []byte
}

// Some servers may return big HTML page for error status (for example 404 page).
// We don't want to treat the whole page as an error message
func (e HttpError) Error() string {

	// We don't want to accidentally dump big pages into the error logs
	if len(e.ResponseBody) < 512 {
		return fmt.Sprintf("HTTP error (%d): %s", e.HttpStatus, string(e.ResponseBody))
	}

	return fmt.Sprintf("HTTP error (%d): <Body length: %d>", e.HttpStatus, len(e.ResponseBody))
}

// Get performs GET request on the specified URI
// URI should be provided without leading '/'
func (c *RESTClient) Get(ctx context.Context, uri string, res any) error {

	// This ensures that no double // appears in the request.
	uri = strings.TrimPrefix(uri, "/")

	furl := fmt.Sprintf("%s/%s", c.BaseURL, uri)

	req, err := http.NewRequestWithContext(ctx, "GET", furl, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	// Parse response body into the res object.
	err = c.do(req, furl, res)
	if err != nil {
		return err
	}

	return nil
}

// parseError tries to unmarshal the JSON output of the response and return an
// appropriate error struct.
//
// The output can be one of the following:
//  1. ApiError - in case both ErrId and Message is provided
//  2. HttpError - in case response was not a JSON response, or not of the ApiError type
//     (ErrId or Message was not specified)
//  3. error - in case JSON cannot be unmarshalled
func (c *RESTClient) parseError(resp *http.Response, reqBody []byte) error {

	// Content type can contain more that just a content type, for example
	// encoding which will be provided after semicolon
	ct := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
	if ct == "application/json" {
		apiError := ApiError{
			HttpStatus: resp.StatusCode,
		}
		err := json.Unmarshal(reqBody, &apiError)

		// If response body cannot be unmarshalled (when Content-Type is "application/json",
		// something went wrong and we should not be very clever about how handle that.
		// Client code should not try to recover from this situation, and probably have to give with and error
		if err != nil {
			return fmt.Errorf("HTTP error (%d, cannot unmarshal JSON error response: %s)", resp.StatusCode, err)
		}

		// If the unmarshalled JSON is probably not of the ApiError type
		if apiError.ErrId == 0 || strings.TrimSpace(apiError.Message) == "" {
			return fmt.Errorf(
				"HTTP error (%d, cannot unmarshal ApiError - no ErrId or Message specified)",
				resp.StatusCode)
		}

		return apiError
	}

	return HttpError{HttpStatus: resp.StatusCode, ResponseBody: reqBody}
}

func (c *RESTClient) do(req *http.Request, furl string, res any) error {

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("error creating request for %q: %w", furl, err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read from response: %w", err)
	}

	//
	if resp.StatusCode != http.StatusOK {

		if c.ErrParser != nil {
			return c.ErrParser(resp, data)
		}

		return c.parseError(resp, data)
	}

	log.Printf("Data: %#v", string(data))

	err = json.Unmarshal(data, res)
	if err != nil {
		return fmt.Errorf("cannot unmarshal response: %w", err)
	}
	return nil
}
