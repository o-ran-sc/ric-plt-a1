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
package restful

import (
	"log"
	"os"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/restapi"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/restapi/operations"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/restapi/operations/a1_e_i_data_delivery"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/restapi/operations/a1_mediator"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/resthooks"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

func NewRestful() *Restful {
	r := &Restful{
		rh: resthooks.NewResthook(),
	}
	r.api = r.setupHandler()
	return r
}

func (r *Restful) setupHandler() *operations.A1API {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		os.Exit(1)
	}

	api := operations.NewA1API(swaggerSpec)

	api.A1MediatorA1ControllerGetHealthcheckHandler = a1_mediator.A1ControllerGetHealthcheckHandlerFunc(func(param a1_mediator.A1ControllerGetHealthcheckParams) middleware.Responder {
		a1.Logger.Debug("handler for get Health Check of A1")
		resp := r.rh.GetA1Health()
		if resp == false {
			return a1_mediator.NewA1ControllerGetHealthcheckInternalServerError()
		}
		return a1_mediator.NewA1ControllerGetHealthcheckOK()
	})

	api.A1MediatorA1ControllerGetAllPolicyTypesHandler = a1_mediator.A1ControllerGetAllPolicyTypesHandlerFunc(func(param a1_mediator.A1ControllerGetAllPolicyTypesParams) middleware.Responder {
		a1.Logger.Debug("handler for get all policy type")
		return a1_mediator.NewA1ControllerGetAllPolicyTypesOK().WithPayload(r.rh.GetAllPolicyType())
	})

	api.A1MediatorA1ControllerCreatePolicyTypeHandler = a1_mediator.A1ControllerCreatePolicyTypeHandlerFunc(func(params a1_mediator.A1ControllerCreatePolicyTypeParams) middleware.Responder {
		a1.Logger.Debug("handler for Create policy type ")
		if err = r.rh.CreatePolicyType(models.PolicyTypeID(params.PolicyTypeID), *params.Body); err == nil {
			//Increase prometheus counter
			return a1_mediator.NewA1ControllerCreatePolicyTypeCreated()
		}
		if r.rh.IsTypeAlready(err) || r.rh.IsTypeMismatch(err) {
			return a1_mediator.NewA1ControllerCreatePolicyTypeBadRequest()
		}
		return a1_mediator.NewA1ControllerCreatePolicyTypeServiceUnavailable()

	})

	api.A1MediatorA1ControllerGetPolicyTypeHandler = a1_mediator.A1ControllerGetPolicyTypeHandlerFunc(func(params a1_mediator.A1ControllerGetPolicyTypeParams) middleware.Responder {
		a1.Logger.Debug("handler for get policy type from policytypeID")
		policyTypeSchema, err := r.rh.GetPolicyType(models.PolicyTypeID(params.PolicyTypeID))
		if err == nil {
			return a1_mediator.NewA1ControllerGetPolicyTypeOK().WithPayload(policyTypeSchema)
		}
		if r.rh.IsPolicyTypeNotFound(err) {
			return a1_mediator.NewA1ControllerGetPolicyTypeNotFound()
		}
		return a1_mediator.NewA1ControllerGetPolicyTypeServiceUnavailable()
	})

	api.A1MediatorA1ControllerCreateOrReplacePolicyInstanceHandler = a1_mediator.A1ControllerCreateOrReplacePolicyInstanceHandlerFunc(func(params a1_mediator.A1ControllerCreateOrReplacePolicyInstanceParams) middleware.Responder {
		a1.Logger.Debug("handler for create policy type instance ")
		if err = r.rh.CreatePolicyInstance(models.PolicyTypeID(params.PolicyTypeID), models.PolicyInstanceID(params.PolicyInstanceID), params.Body); err == nil {
			return a1_mediator.NewA1ControllerCreateOrReplacePolicyInstanceAccepted()
		}
		if r.rh.IsValidJson(err) {
			return a1_mediator.NewA1ControllerCreateOrReplacePolicyInstanceBadRequest()
		}
		if r.rh.IsPolicyTypeNotFound(err) {
			return a1_mediator.NewA1ControllerCreateOrReplacePolicyInstanceNotFound()
		}
		return a1_mediator.NewA1ControllerCreateOrReplacePolicyInstanceServiceUnavailable()

	})

	api.A1MediatorA1ControllerGetPolicyInstanceHandler = a1_mediator.A1ControllerGetPolicyInstanceHandlerFunc(func(params a1_mediator.A1ControllerGetPolicyInstanceParams) middleware.Responder {
		a1.Logger.Debug("handler for get policy instance from policytypeID")
		if resp, err := r.rh.GetPolicyInstance(models.PolicyTypeID(params.PolicyTypeID), models.PolicyInstanceID(params.PolicyInstanceID)); err == nil {
			return a1_mediator.NewA1ControllerGetPolicyInstanceOK().WithPayload(resp)
		} else if r.rh.IsPolicyInstanceNotFound(err) {
			return a1_mediator.NewA1ControllerGetPolicyInstanceNotFound()
		} else if r.rh.IsPolicyTypeNotFound(err) {
			return a1_mediator.NewA1ControllerGetPolicyTypeNotFound()
		}
		return a1_mediator.NewA1ControllerGetPolicyInstanceServiceUnavailable()
	})

	api.A1MediatorA1ControllerGetAllInstancesForTypeHandler = a1_mediator.A1ControllerGetAllInstancesForTypeHandlerFunc(func(params a1_mediator.A1ControllerGetAllInstancesForTypeParams) middleware.Responder {
		a1.Logger.Debug("handler for get all policy instance")
		if resp, err := r.rh.GetAllPolicyInstance(models.PolicyTypeID(params.PolicyTypeID)); err == nil {
			if resp != nil {
				return a1_mediator.NewA1ControllerGetAllInstancesForTypeOK().WithPayload(resp)
			}
		}
		if r.rh.IsPolicyInstanceNotFound(err) {
			return a1_mediator.NewA1ControllerGetPolicyInstanceNotFound()
		}
		return a1_mediator.NewA1ControllerGetAllInstancesForTypeServiceUnavailable()

	})

	api.A1MediatorA1ControllerDeletePolicyTypeHandler = a1_mediator.A1ControllerDeletePolicyTypeHandlerFunc(func(params a1_mediator.A1ControllerDeletePolicyTypeParams) middleware.Responder {
		a1.Logger.Debug("handler for delete policy type")
		if err := r.rh.DeletePolicyType(models.PolicyTypeID(params.PolicyTypeID)); err != nil {
			if r.rh.CanPolicyTypeBeDeleted(err) {
				return a1_mediator.NewA1ControllerDeletePolicyTypeBadRequest()
			}
			if r.rh.IsPolicyTypeNotFound(err) {
				return a1_mediator.NewA1ControllerDeletePolicyTypeNotFound()
			}
			return a1_mediator.NewA1ControllerDeletePolicyTypeServiceUnavailable()
		}
		return a1_mediator.NewA1ControllerDeletePolicyTypeNoContent()

	})

	api.A1MediatorA1ControllerGetPolicyInstanceStatusHandler = a1_mediator.A1ControllerGetPolicyInstanceStatusHandlerFunc(func(params a1_mediator.A1ControllerGetPolicyInstanceStatusParams) middleware.Responder {
		a1.Logger.Debug("handler for get policy instance status")
		if resp, err := r.rh.GetPolicyInstanceStatus(models.PolicyTypeID(params.PolicyTypeID), models.PolicyInstanceID(params.PolicyInstanceID)); err == nil {
			return a1_mediator.NewA1ControllerGetPolicyInstanceStatusOK().WithPayload(resp)
		}
		if r.rh.IsPolicyInstanceNotFound(err) {
			return a1_mediator.NewA1ControllerGetPolicyInstanceStatusNotFound()
		} else if r.rh.IsPolicyTypeNotFound(err) {
			return a1_mediator.NewA1ControllerGetPolicyTypeNotFound()
		}
		return a1_mediator.NewA1ControllerGetPolicyInstanceStatusServiceUnavailable()
	})

	api.A1MediatorA1ControllerDeletePolicyInstanceHandler = a1_mediator.A1ControllerDeletePolicyInstanceHandlerFunc(func(params a1_mediator.A1ControllerDeletePolicyInstanceParams) middleware.Responder {
		a1.Logger.Debug("handler for delete policy instance")
		if err := r.rh.DeletePolicyInstance(models.PolicyTypeID(params.PolicyTypeID), models.PolicyInstanceID(params.PolicyInstanceID)); err != nil {
			if r.rh.CanPolicyInstanceBeDeleted(err) {
				return a1_mediator.NewA1ControllerDeletePolicyTypeBadRequest()
			}
			if r.rh.IsPolicyInstanceNotFound(err) {
				return a1_mediator.NewA1ControllerDeletePolicyInstanceNotFound()
			}
			if r.rh.IsPolicyTypeNotFound(err) {
				return a1_mediator.NewA1ControllerDeletePolicyTypeNotFound()
			}
			return a1_mediator.NewA1ControllerDeletePolicyInstanceServiceUnavailable()
		}
		return a1_mediator.NewA1ControllerDeletePolicyInstanceAccepted()

	})

	api.A1eiDataDeliveryA1ControllerDataDeliveryHandler = a1_e_i_data_delivery.A1ControllerDataDeliveryHandlerFunc(func(params a1_e_i_data_delivery.A1ControllerDataDeliveryParams) middleware.Responder {
		a1.Logger.Debug("handler for EI data delivery")
		if err = r.rh.DataDelivery(params.Body); err != nil {
			return a1_e_i_data_delivery.NewA1ControllerDataDeliveryNotFound()
		}
		return a1_e_i_data_delivery.NewA1ControllerDataDeliveryOK()
	})

	return api

}

func (r *Restful) Run() {

	server := restapi.NewServer(r.api)
	defer server.Shutdown()
	server.Port = 10000
	server.Host = "0.0.0.0"
	if err := server.Serve(); err != nil {
		log.Fatal(err.Error())
	}
}
