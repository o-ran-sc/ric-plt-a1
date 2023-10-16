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
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
)

type Message struct {
}

func (m *Message) PolicyMessage(policyTypeId string, policyInstanceID string, httpBody string, operation string) (string, error) {
        policyTypeIdInt, _ := strconv.ParseInt(policyTypeId, 10, 64)
        
        payloadmap := make(map[string]interface{})

        if len(httpBody) != 0 {
	    err := json.Unmarshal([]byte(httpBody), &payloadmap)
	    if err != nil {
		a1.Logger.Error("Umarshal error : %v", err)
		return "", err
	    }
        }
	//var datajson interface{}
        
	datajson := map[string]interface{}{"operation": operation,
		"payload": payloadmap,
		"policy_instance_id": policyInstanceID,
		"policy_type_id": policyTypeIdInt}

        
        //datajson = map[string]interface{"operation": operation, "payload": payloadmap, "policy_instance_id": policyInstanceID, "policy_type_id": policyTypeIdInt}
	data, err := json.Marshal(datajson)
	if err != nil {
		a1.Logger.Error("marshal error : %v", err)
		return "", err
	}
	return string(data), nil
}

func (m *Message) A1EIMessage(eiJobId string, httpBody string) (string, error) {
	var datajson interface{}
	datajson = map[string]string{
		"ei_job_id": eiJobId,
		"payload":   httpBody}
	data, err := json.Marshal(datajson)

	if err != nil {
		a1.Logger.Error("marshal error : %v", err)
		return "", err
	}
	return string(data), nil
}
