package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	// Application Packages
	"github.com/AnaelBerrouet/movie-night/pkg/app_errors"

	// Third Party Packages
	"github.com/google/uuid"
)

type AppHandler func(
	context.Context,
	http.ResponseWriter,
	*http.Request,
) error

type JSONResponse interface {
	ToJSON() ([]byte, error)
}

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		ctx = context.WithValue(ctx, "bps_request_id", uuid.New())
	} else {
		ctx = context.WithValue(ctx, "bps_request_id", requestID)
	}

	defer func() {
		if r := recover(); r != nil {
			buffer := make([]byte, 2048) // should be greater than max stack trace length
			stackTraceSize := runtime.Stack(buffer, false)
			stackTrace := string(buffer[:stackTraceSize])

			err := errors.New(fmt.Sprintf("Panic in Handler: %q -- \n %s", r, stackTrace))
			log.Println(err)

			SendJSONResponse(w, http.StatusInternalServerError, app_errors.WrapError(err))
		}
	}()

	err := fn(ctx, w, r)
	if err == nil {
		return
	}

	response, code := app_errors.WrapError(err), http.StatusInternalServerError
	SendJSONResponse(w, code, response)
}

func SendTextResponse(w http.ResponseWriter, httpStatus int, body string) {
	// Always add secure browser headers to response
	AddSecureHeaders(w)
	// Add correct content type to response
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	w.WriteHeader(httpStatus)
	w.Write([]byte(body))

	return
}

func SendInternalServerErrorResponse(w http.ResponseWriter, err error) {
	// Always add secure browser headers to response
	AddSecureHeaders(w)
	// Add correct content type to response
	w.Header().Set("Content-Type", "application/json")

	error_response, _ := json.Marshal(app_errors.WrapError(err))

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(error_response)

	return
}

func SendJSONResponse(w http.ResponseWriter, httpStatus int, body JSONResponse) {

	js, err := body.ToJSON()
	if err != nil {
		SendInternalServerErrorResponse(w, err)
		return
	}

	// Always add secure browser headers to response
	AddSecureHeaders(w)
	// Add correct content type to response
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(httpStatus)
	w.Write(js)

	return
}

func DecodeBodyAsMap(req *http.Request) (map[string]interface{}, error) {
	requestBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	var params map[string]interface{}

	if err := json.Unmarshal(requestBody, &params); err != nil {
		return nil, err
	}

	return params, nil
}

// AddSecureHeaders Adds Secure Headers to HTTP responses
func AddSecureHeaders(w http.ResponseWriter) {

	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("X-Xss-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("x-Download-Options", "noopen")
	w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
	w.Header().Set("Cross-Origin-Window-Policy", "deny")

	// HSTS Header teller clients to connect with HTTPS
	w.Header().Set("Strict-Transport-Security", "max-age=31536000")

	return
}
