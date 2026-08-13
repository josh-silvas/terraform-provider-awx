package awx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil" //nolint: staticcheck
	"net/http"
	"net/url"
	"strings"
)

// APIRequest represents the http api communication way.
type APIRequest struct {
	Method   string
	Endpoint string
	Payload  io.Reader
	Headers  http.Header
	Suffix   string
}

// SetHeader sets http header by passing k,v.
func (ar *APIRequest) SetHeader(key string, value string) *APIRequest {
	ar.Headers.Set(key, value)
	return ar
}

// NewAPIRequest news an APIRequest object.
func NewAPIRequest(method string, endpoint string, payload io.Reader) *APIRequest {
	var headers = http.Header{}
	var suffix string
	ar := &APIRequest{method, endpoint, payload, headers, suffix}
	return ar
}

// Authenticator represents an authentication interface.
type Authenticator interface {
	addAuthenticationHeaders(*http.Request)
}

// BasicAuth represents http basic auth.
type BasicAuth struct {
	Username string
	Password string
}

func (ba *BasicAuth) addAuthenticationHeaders(r *http.Request) {
	r.SetBasicAuth(ba.Username, ba.Password)
}

// TokenAuth represents token authentication.
type TokenAuth struct {
	Token string
}

func (ta *TokenAuth) addAuthenticationHeaders(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ta.Token))
}

// Requester implemented a base http client.
// It supports do POST/GET via an human-readable way,
// in other word, all data is in `application/json` format.
// It also originally supports basic auth.
// For production usage, It would be better to wrapper
// an another rest client on this requester.
type Requester struct {
	Base          string
	Authenticator Authenticator
	Client        *http.Client
}

// Do : Performs the actual http request.
func (r *Requester) Do(ar *APIRequest, responseStruct interface{}, options ...interface{}) (*http.Response, error) { //nolint:funlen
	if !strings.HasSuffix(ar.Endpoint, "/") && ar.Method != "POST" {
		ar.Endpoint += "/"
	}

	URL, err := url.Parse(r.Base + ar.Endpoint + ar.Suffix)
	if err != nil {
		return nil, err
	}

	for _, o := range options {
		switch v := o.(type) {
		case map[string]string:
			querystring := make(url.Values)
			for key, val := range v {
				querystring.Set(key, val)
			}

			URL.RawQuery = querystring.Encode()
		}
	}

	var req *http.Request
	req, err = http.NewRequest(ar.Method, URL.String(), ar.Payload)
	if err != nil {
		return nil, err
	}

	r.Authenticator.addAuthenticationHeaders(req)

	for k := range ar.Headers {
		req.Header.Add(k, ar.Headers.Get(k))
	}

	response, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Do.Request: %v", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		body, readErr := ioutil.ReadAll(response.Body)
		if closeErr := response.Body.Close(); closeErr != nil {
			fmt.Println(closeErr)
		}
		if readErr != nil {
			return response, readErr
		}
		// Restore the body so callers calling CheckResponse still see it.
		response.Body = ioutil.NopCloser(bytes.NewReader(body))

		if response.StatusCode == http.StatusBadRequest {
			errorList := map[string]interface{}{}
			if err := json.Unmarshal(body, &errorList); err == nil {
				errorString := "Errors:"
				for k, v := range errorList {
					errorString = fmt.Sprintf("%s\n- %s: %+v", errorString, k, v)
				}

				return response, errors.New(errorString)
			}
		}

		return response, HTTPError(response.StatusCode, body)
	}

	// If there is no response body, or if the response is of type `No Content` bypass the decode process.
	if responseStruct == nil || response.StatusCode == 204 { //nolint:gomnd
		return response, nil
	}

	switch responseStruct.(type) {
	case *string:
		return r.ReadRawResponse(response, responseStruct)
	default:
		return r.ReadJSONResponse(response, responseStruct)
	}
}

// ReadRawResponse reads the http raw response and store it into `responseStruct`.
func (r *Requester) ReadRawResponse(response *http.Response, responseStruct interface{}) (*http.Response, error) {
	defer func() {
		if err := response.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if str, ok := responseStruct.(*string); ok {
		*str = string(content)
	} else {
		return nil, fmt.Errorf("could not cast responseStruct to *string")
	}

	return response, nil
}

// ReadJSONResponse reads the http raw response and decodes into json.
func (r *Requester) ReadJSONResponse(response *http.Response, responseStruct interface{}) (*http.Response, error) {
	if response.ContentLength == 0 {
		return response, nil
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if err := json.NewDecoder(response.Body).Decode(responseStruct); err != nil {
		return response, err
	}
	return response, nil
}

// Get performs http get request.
func (r *Requester) Get(endpoint string, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	ar := NewAPIRequest("GET", endpoint, nil)
	ar.Suffix = ""
	return r.Do(ar, responseStruct, querystring)
}

// GetJSON performs http get request with json response.
func (r *Requester) GetJSON(endpoint string, responseStruct interface{}, query map[string]string) (*http.Response, error) {
	ar := NewAPIRequest("GET", endpoint, nil)
	ar.SetHeader("Content-Type", "application/json")
	ar.Suffix = ""
	return r.Do(ar, &responseStruct, query)
}

// Post performs http post request.
func (r *Requester) Post(endpoint string, payload io.Reader, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	ar := NewAPIRequest("POST", endpoint, payload)
	ar.SetHeader("Content-Type", "application/json")
	ar.Suffix = ""
	return r.Do(ar, &responseStruct, querystring)
}

// PutJSON perform http PUT request with json response.
func (r *Requester) PutJSON(endpoint string, payload io.Reader, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	ar := NewAPIRequest("PUT", endpoint, payload)
	ar.SetHeader("Content-Type", "application/json")
	ar.Suffix = ""
	return r.Do(ar, &responseStruct, querystring)
}

// PostJSON performs http post request with json response.
func (r *Requester) PostJSON(endpoint string, payload io.Reader, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	ar := NewAPIRequest("POST", endpoint, payload)
	ar.SetHeader("Content-Type", "application/json")
	ar.Suffix = ""
	return r.Do(ar, &responseStruct, querystring)
}

// PatchJSON perform http patch request with json response.
func (r *Requester) PatchJSON(endpoint string, payload io.Reader, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	ar := NewAPIRequest("PATCH", endpoint, payload)
	ar.SetHeader("Content-Type", "application/json")
	ar.Suffix = ""
	return r.Do(ar, &responseStruct, querystring)
}

// Delete performs http Delete request.
func (r *Requester) Delete(endpoint string, responseStruct interface{}, querystring map[string]string) (*http.Response, error) {
	ar := NewAPIRequest("DELETE", endpoint, nil)
	ar.Suffix = ""
	return r.Do(ar, responseStruct, querystring)
}
