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
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	"strconv"
)

const (
	a1HandlerPrefix = "a1.policy_handler."
	a1MediatorNs    = "A1m_ns"
)

func NewInstanceManager() *InstanceManager {
	return createInstanceManager(sdlgo.NewSyncStorage())
}

func createInstanceManager(sdlInst iSdl) *InstanceManager {
	im := &InstanceManager{
		db: sdlInst,
	}
	return im
}
func (im *InstanceManager) SetPolicyInstanceStatus(policyTypeId int, policyInstanceID int, status string) error {
	a1.Logger.debug("message recieved for %d and %d", policyTypeId, policyInstanceID)
	instancehandlerKey := a1HandlerPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + strconv.FormatInt((int64(policyInstanceID)), 10)
	err := im.db.Set(a1MediatorNs, instancehandlerKey, status)
	if err != nil {
		a1.Logger.Error("error1 :%+v", err)
		return err
	}
	return nil
}
