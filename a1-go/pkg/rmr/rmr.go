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
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

const (
	a1SourceName = "service-ricplt-a1mediator-http"
)

type RmrSender struct {
	rmrclient     *xapp.RMRClient
	policyManager *policy.PolicyManager
}

type IRmrSender interface {
	RmrSendToXapp(httpBodyString string, messagetype int) bool
	RmrRecieveStart()
}

func NewRMRSender(policyManager *policy.PolicyManager) IRmrSender {
	RMRclient := xapp.NewRMRClientWithParams(&xapp.RMRClientParams{
		StatDesc: "",
		RmrData: xapp.PortData{
			Name:              "",
			MaxSize:           65534,
			ThreadType:        0,
			LowLatency:        false,
			FastAck:           false,
			MaxRetryOnFailure: 1,
			Port:              4561,
			Policies:          []int{20012, 20013, 20017, 20015, 20011, 20010},
		},
	})
	return &RmrSender{
		rmrclient:     RMRclient,
		policyManager: policyManager,
	}
}

func (rmr *RmrSender) RmrSendToXapp(httpBodyString string, messagetype int) bool {

	params := &xapp.RMRParams{}
	params.Mtype = messagetype
	params.SubId = -1
	params.Xid = ""
	params.Meid = &xapp.RMRMeid{}
	params.Src = a1SourceName
	params.PayloadLen = len([]byte(httpBodyString))
	params.Payload = []byte(httpBodyString)
	a1.Logger.Debug("MSG to XAPP: %s ", params.String())
	a1.Logger.Debug("len payload %+v", len(params.Payload))
	s := rmr.rmrclient.SendMsg(params)
	a1.Logger.Debug("rmrSendToXapp: sending: %+v", s)
	return s
}

func (rmr *RmrSender) Consume(msg *xapp.RMRParams) (err error) {
	a1.Logger.Debug("In the Consume function")
	id := xapp.Rmr.GetRicMessageName(msg.Mtype)
	a1.Logger.Debug("Message received: name=%s meid=%s subId=%d txid=%s len=%d", id, msg.Meid.RanName, msg.SubId, msg.Xid, msg.PayloadLen)

	switch id {

	case "A1_POLICY_RESP":
		a1.Logger.Debug("Recived policy responose")
		payload := msg.Payload
		a1.Logger.Debug("message recieved : %s", payload)
		var result map[string]interface{}
		err := json.Unmarshal([]byte(payload), &result)
		if err != nil {
			a1.Logger.Error("Unmarshal error : %+v", err)
		}
		a1.Logger.Debug("message recieved for %d and %d with status : %s", result["policy_type_id"], result["policy_instance_id"], result["status"])
		rmr.policyManager.SetPolicyInstanceStatus(int(result["policy_type_id"].(float64)), int(result["policy_instance_id"].(float64)), result["status"].(string))
	default:
		xapp.Logger.Error("Unknown message type '%d', discarding", msg.Mtype)
	}

	defer func() {
		rmr.rmrclient.Free(msg.Mbuf)
		msg.Mbuf = nil
	}()
	return
}

func (rmr *RmrSender) RmrRecieveStart() {
	a1.Logger.Debug("Inside RmrRecieveStart function ")
	rmr.rmrclient.Start(rmr)
	a1.Logger.Debug("Reciever started")
}
