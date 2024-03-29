// Code generated by go-swagger; DO NOT EDIT.

package todos

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/kaenova/go-playground/swag/models"
)

// GetTodosOKCode is the HTTP code returned for type GetTodosOK
const GetTodosOKCode int = 200

/*GetTodosOK list the todo operations

swagger:response getTodosOK
*/
type GetTodosOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Item `json:"body,omitempty"`
}

// NewGetTodosOK creates GetTodosOK with default headers values
func NewGetTodosOK() *GetTodosOK {

	return &GetTodosOK{}
}

// WithPayload adds the payload to the get todos o k response
func (o *GetTodosOK) WithPayload(payload []*models.Item) *GetTodosOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get todos o k response
func (o *GetTodosOK) SetPayload(payload []*models.Item) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTodosOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Item, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*GetTodosDefault generic error response

swagger:response getTodosDefault
*/
type GetTodosDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetTodosDefault creates GetTodosDefault with default headers values
func NewGetTodosDefault(code int) *GetTodosDefault {
	if code <= 0 {
		code = 500
	}

	return &GetTodosDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get todos default response
func (o *GetTodosDefault) WithStatusCode(code int) *GetTodosDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get todos default response
func (o *GetTodosDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get todos default response
func (o *GetTodosDefault) WithPayload(payload *models.Error) *GetTodosDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get todos default response
func (o *GetTodosDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTodosDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
