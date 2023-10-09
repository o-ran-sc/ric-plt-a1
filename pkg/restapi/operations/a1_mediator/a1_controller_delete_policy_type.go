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
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// A1ControllerDeletePolicyTypeHandlerFunc turns a function with the right signature into a a1 controller delete policy type handler
type A1ControllerDeletePolicyTypeHandlerFunc func(A1ControllerDeletePolicyTypeParams) middleware.Responder

// Handle executing the request and returning a response
func (fn A1ControllerDeletePolicyTypeHandlerFunc) Handle(params A1ControllerDeletePolicyTypeParams) middleware.Responder {
	return fn(params)
}

// A1ControllerDeletePolicyTypeHandler interface for that can handle valid a1 controller delete policy type params
type A1ControllerDeletePolicyTypeHandler interface {
	Handle(A1ControllerDeletePolicyTypeParams) middleware.Responder
}

// NewA1ControllerDeletePolicyType creates a new http.Handler for the a1 controller delete policy type operation
func NewA1ControllerDeletePolicyType(ctx *middleware.Context, handler A1ControllerDeletePolicyTypeHandler) *A1ControllerDeletePolicyType {
	return &A1ControllerDeletePolicyType{Context: ctx, Handler: handler}
}

/* A1ControllerDeletePolicyType swagger:route DELETE /A1-P/v2/policytypes/{policy_type_id} A1 Mediator a1ControllerDeletePolicyType

Delete this policy type. Can only be performed if there are no instances of this type


*/
type A1ControllerDeletePolicyType struct {
	Context *middleware.Context
	Handler A1ControllerDeletePolicyTypeHandler
}

func (o *A1ControllerDeletePolicyType) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewA1ControllerDeletePolicyTypeParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
