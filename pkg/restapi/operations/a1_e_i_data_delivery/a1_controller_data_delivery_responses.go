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

package a1_e_i_data_delivery

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// A1ControllerDataDeliveryOKCode is the HTTP code returned for type A1ControllerDataDeliveryOK
const A1ControllerDataDeliveryOKCode int = 200

/*A1ControllerDataDeliveryOK successfully delivered data from data producer


swagger:response a1ControllerDataDeliveryOK
*/
type A1ControllerDataDeliveryOK struct {
}

// NewA1ControllerDataDeliveryOK creates A1ControllerDataDeliveryOK with default headers values
func NewA1ControllerDataDeliveryOK() *A1ControllerDataDeliveryOK {

	return &A1ControllerDataDeliveryOK{}
}

// WriteResponse to the client
func (o *A1ControllerDataDeliveryOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// A1ControllerDataDeliveryNotFoundCode is the HTTP code returned for type A1ControllerDataDeliveryNotFound
const A1ControllerDataDeliveryNotFoundCode int = 404

/*A1ControllerDataDeliveryNotFound no job id defined for this data delivery


swagger:response a1ControllerDataDeliveryNotFound
*/
type A1ControllerDataDeliveryNotFound struct {
}

// NewA1ControllerDataDeliveryNotFound creates A1ControllerDataDeliveryNotFound with default headers values
func NewA1ControllerDataDeliveryNotFound() *A1ControllerDataDeliveryNotFound {

	return &A1ControllerDataDeliveryNotFound{}
}

// WriteResponse to the client
func (o *A1ControllerDataDeliveryNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}