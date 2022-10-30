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
package rmr

import (
	"encoding/json"
	"strconv"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
)

type Message struct {
}

type IMessage interface {
	PolicyMessage(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID, httpBody interface{}, operation string) string
}

func (m *Message) PolicyMessage(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID, httpBody interface{}, operation string) string {
	var datajson []interface{}
	datajson = append(datajson, map[string]string{"operation": operation,
		"policy_type_id":     strconv.FormatInt((int64(policyTypeId)), 10),
		"policy_instance_id": string(policyInstanceID),
		"payload":            httpBody.(string)})
	data, _ := json.Marshal(datajson)

	return string(data)
}
