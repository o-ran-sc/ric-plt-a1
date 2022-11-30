/*
==================================================================================
  Copyright (c) 2022 Samsung

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

package policy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
)

var policyTypeNotFoundError = errors.New("Policy Type Not Found")
var policyInstanceNotFoundError = errors.New("Policy Instance Not Found")

const (
	a1HandlerPrefix  = "a1.policy_handler."
	a1PolicyPrefix   = "a1.policy_type."
	a1MediatorNs     = "A1m_ns"
	a1InstancePrefix = "a1.policy_instance."
)

func NewPolicyManager(sdl *sdlgo.SyncStorage) *PolicyManager {
	return createPolicyManager(sdl)
}

func createPolicyManager(sdlInst iSdl) *PolicyManager {
	pm := &PolicyManager{
		db: sdlInst,
	}
	return pm
}
func (pm *PolicyManager) SetPolicyInstanceStatus(policyTypeId int, policyInstanceID int, status string) error {
	a1.Logger.Debug("message recieved for %d and %d", policyTypeId, policyInstanceID)
	instancehandlerKey := a1HandlerPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + strconv.FormatInt((int64(policyInstanceID)), 10)
	err := pm.db.Set(a1MediatorNs, instancehandlerKey, status)
	if err != nil {
		a1.Logger.Error("error1 :%+v", err)
		return err
	}
	return nil
}

func (im *PolicyManager) GetAllPolicyInstance(policyTypeId int) ([]models.PolicyInstanceID, error) {
	a1.Logger.Debug("GetAllPolicyInstance")
	var policyTypeInstances = []models.PolicyInstanceID{}

	keys, err := im.db.GetAll("A1m_ns")

	if err != nil {
		a1.Logger.Error("error in retrieving policy. err: %v", err)
		return policyTypeInstances, err
	}
	a1.Logger.Debug("keys : %+v", keys)
	typekey := a1InstancePrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "."

	for _, key := range keys {
		if strings.HasPrefix(strings.TrimLeft(key, " "), typekey) {
			pti := strings.Split(strings.Trim(key, " "), typekey)[1]
			a1.Logger.Debug("pti %+v", pti)
			policyTypeInstances = append(policyTypeInstances, models.PolicyInstanceID(pti))
		}
	}

	if len(policyTypeInstances) == 0 {
		a1.Logger.Debug("policy instance Not Present  ")
		return policyTypeInstances, policyInstanceNotFoundError
	}

	a1.Logger.Debug("return : %+v", policyTypeInstances)
	return policyTypeInstances, nil
}

func (im *PolicyManager) GetPolicyInstance(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID) (interface{}, error) {
	a1.Logger.Debug("GetPolicyInstance1")

	var keys [1]string

	typekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
	keys[0] = typekey

	a1.Logger.Debug("key1 : %+v", typekey)

	valmap, err := im.db.Get(a1MediatorNs, keys[:])
	if len(valmap) == 0 {
		a1.Logger.Debug("policy type Not Present for policyid : %v", policyTypeId)
		return nil, policyTypeNotFoundError
	}

	if err != nil {
		a1.Logger.Error("error in retrieving policy type. err: %v", err)
		return nil, err
	}

	if valmap[typekey] == nil {
		a1.Logger.Debug("policy type Not Present for policyid : %v", policyTypeId)
		return nil, policyTypeNotFoundError
	}

	a1.Logger.Debug("keysmap : %+v", valmap[typekey])

	instancekey := a1InstancePrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
	a1.Logger.Debug("key2 : %+v", instancekey)
	keys[0] = instancekey
	instanceMap, err := im.db.Get(a1MediatorNs, keys[:])
	if err != nil {
		a1.Logger.Error("policy instance error : %v", err)
		return nil, err
	}
	a1.Logger.Debug("policyinstancetype map : %+v", instanceMap)

	if instanceMap[instancekey] == nil {
		a1.Logger.Debug("policy instance Not Present for policyinstaneid : %v", policyInstanceID)
		return nil, policyInstanceNotFoundError
	}

	valStr := fmt.Sprint(instanceMap[instancekey])
	return valStr, nil
}
