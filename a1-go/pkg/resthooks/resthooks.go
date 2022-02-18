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
package resthooks

import (
	"strconv"
	"strings"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
)

const (
	a1policyprefix = "a1.policy_type."
	a1MediatorNs   = "A1m_ns"
)

func NewResthook() *Resthook {
	return createResthook(sdlgo.NewSyncStorage())
}

func createResthook(sdlInst iSdl) *Resthook {
	return &Resthook{
		db: sdlInst,
	}
}

func (rh *Resthook) GetAllPolicyType() []models.PolicyTypeID {

	var policyTypeIDs []models.PolicyTypeID

	keys, err := rh.db.GetAll("A1m_ns")

	if err != nil {
		a1.Logger.Error("error in retrieving policy. err: %v", err)
		return policyTypeIDs
	}
	a1.Logger.Debug("keys : %+v", keys)

	for _, key := range keys {
		if strings.HasPrefix(strings.TrimLeft(key, " "), "a1.policy_type.") {
			pti := strings.Split(strings.Trim(key, " "), "a1.policy_type.")[1]
			ptii, _ := strconv.ParseInt(pti, 10, 64)
			policyTypeIDs = append(policyTypeIDs, models.PolicyTypeID(ptii))
		}
	}

	a1.Logger.Debug("return : %+v", policyTypeIDs)
	return policyTypeIDs
}

func (rh *Resthook) CreatePolicyType(policyTypeId models.PolicyTypeID, httprequest models.PolicyTypeSchema) bool {
	a1.Logger.Debug("CreatePolicyType function")
	if policyTypeId != models.PolicyTypeID(*httprequest.PolicyTypeID) {
		//error message
		a1.Logger.Debug("Policytype Mismatch")
		return false
	}
	key := a1policyprefix + strconv.FormatInt((int64(policyTypeId)), 10)
	a1.Logger.Debug("key %+v ", key)
	if data, err := httprequest.MarshalBinary(); err == nil {
		a1.Logger.Debug("Mrshaled String : %+v", string(data))
		success, err1 := rh.db.SetIfNotExists(a1MediatorNs, key, string(data))
		a1.Logger.Error("success:%+v", success)
		if err1 != nil {
			a1.Logger.Error("error :%+v", err1)
			return false
		}
		return success
	}
	return false
}
