package jsonutil

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Resp struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ReadJSON reads data into the data param(it assumes data is a reference type)
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) error {
	var payload Resp
	payload.Error = true
	payload.Message = err.Error()

	out, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

	return nil
}

// WriteJSON writes arbitrary data out as JSON
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

func ErrorJSON(w http.ResponseWriter, logger *zap.Logger, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	logger.Error(err.Error())

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload Resp
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(w, statusCode, payload)
}
