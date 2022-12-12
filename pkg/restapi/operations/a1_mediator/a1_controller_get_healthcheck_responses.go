/*
==================================================================================
  Copyright (c) 2021 Samsung

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   This source code is part of the near-RT RIC (RAN Intelligent Controller)
   platform project (RICP).
==================================================================================
*/
// Code generated by go-swagger; DO NOT EDIT.

package a1_mediator

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// A1ControllerGetHealthcheckOKCode is the HTTP code returned for type A1ControllerGetHealthcheckOK
const A1ControllerGetHealthcheckOKCode int = 200

/*A1ControllerGetHealthcheckOK A1 is healthy. Anything other than a 200 should be considered a1 as failing


swagger:response a1ControllerGetHealthcheckOK
*/
type A1ControllerGetHealthcheckOK struct {
}

// NewA1ControllerGetHealthcheckOK creates A1ControllerGetHealthcheckOK with default headers values
func NewA1ControllerGetHealthcheckOK() *A1ControllerGetHealthcheckOK {

	return &A1ControllerGetHealthcheckOK{}
}

// WriteResponse to the client
func (o *A1ControllerGetHealthcheckOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// A1ControllerGetHealthcheckInternalServerErrorCode is the HTTP code returned for type A1ControllerGetHealthcheckInternalServerError
const A1ControllerGetHealthcheckInternalServerErrorCode int = 500

/*A1ControllerGetHealthcheckInternalServerError Internal error to signal A1 is not healthy. Client should attempt to retry later.

swagger:response a1ControllerGetHealthcheckInternalServerError
*/
type A1ControllerGetHealthcheckInternalServerError struct {
}

// NewA1ControllerGetHealthcheckInternalServerError creates A1ControllerGetHealthcheckInternalServerError with default headers values
func NewA1ControllerGetHealthcheckInternalServerError() *A1ControllerGetHealthcheckInternalServerError {

	return &A1ControllerGetHealthcheckInternalServerError{}
}

// WriteResponse to the client
func (o *A1ControllerGetHealthcheckInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}