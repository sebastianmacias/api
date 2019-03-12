package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Response ...
type Response struct {
	Success   bool        `json:"success"`
	Error     bool        `json:"error"`
	Warning   bool        `json:"warning"`
	Type      string      `json:"type"`
	Msg       string      `json:"msg"`
	Code      int         `json:"code,omitempty"`
	Payload   interface{} `json:"payload"`
	Timestamp int64       `json:"timestamp"`
	Actions   []*Action   `json:"actions"`
}

// AddAction ...
func (r *Response) AddAction(actionType, name, code, method, url string, required bool) error {
	r.Actions = append(r.Actions, NewAction(actionType, name, code, method, url, required))
	return nil
}

// Action ...Must add payloadIn and payloadOut
type Action struct {
	Type     string `json:"type,omitempty"`
	Name     string `json:"name"`
	Code     string `json:"code,omitempty"`
	Method   string `json:"method,omitempty"`
	URL      string `json:"url,omitempty"`
	Required bool   `json:"required,omitempty"`
}

// NewAction ...
func NewAction(actionType, name, code, method, url string, required bool) *Action {
	return &Action{Type: actionType, Name: name, Code: code, Method: method, URL: url, Required: required}
}

// NewAPIRes ...
func NewAPIRes(Success, Error, Warning bool, Type, Msg string, Code int, Payload interface{}) *Response {
	return &Response{Success: Success, Error: Error, Warning: Warning, Type: Type, Msg: Msg, Code: Code, Payload: Payload, Timestamp: time.Now().Unix()}
}

// NewAPIResOk ...
func NewAPIResOk(Msg string, Code int, Payload interface{}) *Response {
	return NewAPIRes(true, false, false, "success", Msg, Code, Payload)
}

// NewAPIResErr ...
func NewAPIResErr(Msg string, Code int, Payload interface{}) *Response {
	return NewAPIRes(false, true, false, "error", Msg, Code, Payload)
}

// NewAPIResErrFromError ...
func NewAPIResErrFromError(Err error, Code int, Payload interface{}) *Response {
	return NewAPIResErr(Err.Error(), Code, Payload)
}

// NewAPIResInfo ...
func NewAPIResInfo(Msg string, Code int, Payload interface{}) *Response {
	return NewAPIRes(false, true, false, "info", Msg, Code, Payload)
}

// NewAPIResWarn ...
func NewAPIResWarn(Msg string, Code int, Payload interface{}) *Response {
	return NewAPIRes(false, false, true, "warning", Msg, Code, Payload)
}

// Return helpers

// returnAPIRes ...
func returnAPIRes(w http.ResponseWriter, statusCode int, outPayload []byte) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(outPayload)
	return
}

// ReturnAPIResErr ...
func returnAPIResErr(w http.ResponseWriter, Msg string, Code int) {
	outPayload, _ := Marshal(NewAPIResErr(Msg, Code, nil))
	returnAPIRes(w, http.StatusBadGateway, outPayload)
	return
}

// ReturnAPIResOk ...
func ReturnAPIResOk(w http.ResponseWriter, Msg string, Code int, Payload interface{}) {
	if outPayload, err := Marshal(NewAPIResOk(Msg, Code, Payload)); err != nil {
		returnAPIResErr(w, err.Error(), 0)
	} else {
		returnAPIRes(w, http.StatusOK, outPayload)
	}
	return
}

// ReturnAPIResErr ...
func ReturnAPIResErr(w http.ResponseWriter, Msg string, Code int, Payload interface{}) {
	if outPayload, err := Marshal(NewAPIResErr(Msg, Code, Payload)); err != nil {
		returnAPIResErr(w, err.Error(), 0)
	} else {
		returnAPIRes(w, http.StatusOK, outPayload)
	}
	return
}

// ReturnAPIResInfo ...
func ReturnAPIResInfo(w http.ResponseWriter, Msg string, Code int, Payload interface{}) {
	if outPayload, err := Marshal(NewAPIResInfo(Msg, Code, Payload)); err != nil {
		returnAPIResErr(w, err.Error(), 0)
	} else {
		returnAPIRes(w, http.StatusOK, outPayload)
	}
	return
}

// ReturnAPIResWarn ...
func ReturnAPIResWarn(w http.ResponseWriter, Msg string, Code int, Payload interface{}) {
	if outPayload, err := Marshal(NewAPIResWarn(Msg, Code, Payload)); err != nil {
		returnAPIResErr(w, err.Error(), 0)
	} else {
		returnAPIRes(w, http.StatusOK, outPayload)
	}
	return
}

//General utility functions

// Marshal converts interface to JSON
func Marshal(inPayload interface{}) ([]byte, error) {
	outPayload, err := json.Marshal(inPayload)
	if err != nil {
		return nil, errors.New("Error marshaling JSON")
	}
	return outPayload, nil
}
