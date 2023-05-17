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

// A1ControllerDeletePolicyInstanceAcceptedCode is the HTTP code returned for type A1ControllerDeletePolicyInstanceAccepted
const A1ControllerDeletePolicyInstanceAcceptedCode int = 202

/*
A1ControllerDeletePolicyInstanceAccepted policy instance deletion initiated

swagger:response a1ControllerDeletePolicyInstanceAccepted
*/
type A1ControllerDeletePolicyInstanceAccepted struct {
}

// NewA1ControllerDeletePolicyInstanceAccepted creates A1ControllerDeletePolicyInstanceAccepted with default headers values
func NewA1ControllerDeletePolicyInstanceAccepted() *A1ControllerDeletePolicyInstanceAccepted {

	return &A1ControllerDeletePolicyInstanceAccepted{}
}

// WriteResponse to the client
func (o *A1ControllerDeletePolicyInstanceAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// A1ControllerDeletePolicyInstanceNotFoundCode is the HTTP code returned for type A1ControllerDeletePolicyInstanceNotFound
const A1ControllerDeletePolicyInstanceNotFoundCode int = 404

/*
A1ControllerDeletePolicyInstanceNotFound there is no policy instance with this policy_instance_id or there is no policy type with this policy_type_id

swagger:response a1ControllerDeletePolicyInstanceNotFound
*/
type A1ControllerDeletePolicyInstanceNotFound struct {
}

// NewA1ControllerDeletePolicyInstanceNotFound creates A1ControllerDeletePolicyInstanceNotFound with default headers values
func NewA1ControllerDeletePolicyInstanceNotFound() *A1ControllerDeletePolicyInstanceNotFound {

	return &A1ControllerDeletePolicyInstanceNotFound{}
}

// WriteResponse to the client
func (o *A1ControllerDeletePolicyInstanceNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// A1ControllerDeletePolicyInstanceServiceUnavailableCode is the HTTP code returned for type A1ControllerDeletePolicyInstanceServiceUnavailable
const A1ControllerDeletePolicyInstanceServiceUnavailableCode int = 503

/*
A1ControllerDeletePolicyInstanceServiceUnavailable Potentially transient backend database error. Client should attempt to retry later.

swagger:response a1ControllerDeletePolicyInstanceServiceUnavailable
*/
type A1ControllerDeletePolicyInstanceServiceUnavailable struct {
}

// NewA1ControllerDeletePolicyInstanceServiceUnavailable creates A1ControllerDeletePolicyInstanceServiceUnavailable with default headers values
func NewA1ControllerDeletePolicyInstanceServiceUnavailable() *A1ControllerDeletePolicyInstanceServiceUnavailable {

	return &A1ControllerDeletePolicyInstanceServiceUnavailable{}
}

// WriteResponse to the client
func (o *A1ControllerDeletePolicyInstanceServiceUnavailable) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(503)
}
