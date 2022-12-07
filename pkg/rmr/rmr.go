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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/policy"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

const (
	a1SourceName      = "service-ricplt-a1mediator-http"
	a1PolicyRequest   = 20010
	ecsServiceHost    = "http://ecs-service:8083"
	ecsEiTypePath     = ecsServiceHost + "/A1-EI/v1/eitypes"
	ecsEiJobPath      = ecsServiceHost + "/A1-EI/v1/eijobs/"
	a1EiQueryAllResp  = 20014
	a1EiCreateJobResp = 20016
	jobCreationData   = `{"ei_job_id": %s.}`
)

type RmrSender struct {
	rmrclient     *xapp.RMRClient
	policyManager *policy.PolicyManager
}

type IRmrSender interface {
	RmrSendToXapp(httpBodyString string, messagetype int) bool
}

func NewRMRSender(policyManager *policy.PolicyManager) IRmrSender {
	RMRclient := xapp.NewRMRClientWithParams(&xapp.RMRClientParams{
		StatDesc: "",
		RmrData: xapp.PortData{
			//TODO: Read configuration from config file
			Name:              "",
			MaxSize:           65534,
			ThreadType:        0,
			LowLatency:        false,
			FastAck:           false,
			MaxRetryOnFailure: 1,
			Port:              4561,
		},
	})

	rmrsender := &RmrSender{
		rmrclient:     RMRclient,
		policyManager: policyManager,
	}

	rmrsender.RmrRecieveStart()
	return rmrsender
}

var RICMessageTypes = map[string]int{
	"A1_POLICY_REQ":         20010,
	"A1_POLICY_RESP":        20011,
	"A1_POLICY_QUERY":       20012,
	"A1_EI_QUERY_ALL":       20013,
	"AI_EI_QUERY_ALL_RESP":  20014,
	"A1_EI_CREATE_JOB":      20015,
	"A1_EI_CREATE_JOB_RESP": 20016,
	"A1_EI_DATA_DELIVERY":   20017,
}

func (rmr *RmrSender) GetRicMessageName(id int) (s string) {
	for k, v := range RICMessageTypes {
		if id == v {
			return k
		}
	}
	return
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
	id := rmr.GetRicMessageName(msg.Mtype)
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
			return err
		}
		a1.Logger.Debug("message recieved for %d and %d with status : %s", result["policy_type_id"], result["policy_instance_id"], result["status"])
		rmr.policyManager.SetPolicyInstanceStatus(int(result["policy_type_id"].(float64)), int(result["policy_instance_id"].(float64)), result["status"].(string))
	case "A1_POLICY_QUERY":
		a1.Logger.Debug("Recived policy query")
		a1.Logger.Debug("message recieved ", msg.Payload)
		payload := msg.Payload
		var result map[string]interface{}
		json.Unmarshal([]byte(payload), &result)
		a1.Logger.Debug("message recieved : %s for %d and %d", result, result["policy_type_id"], result["policy_instance_id"])
		policytypeid := (result["policy_type_id"].(float64))
		instanceList, err1 := rmr.policyManager.GetAllPolicyInstance(int(policytypeid))
		if err1 != nil {
			a1.Logger.Error("Error : %+v", err1)
			return err1
		}
		a1.Logger.Debug("instanceList ", instanceList)
		a1.Logger.Debug("Received a query for a known policy type: %d", policytypeid)
		for _, policyinstanceid := range instanceList {
			policyinstance, err2 := rmr.policyManager.GetPolicyInstance(models.PolicyTypeID(policytypeid), policyinstanceid)
			if err2 != nil {
				a1.Logger.Error("Error : %+v", err2)
				return err2
			}
			a1.Logger.Debug("policyinstance ", policyinstance.(string))
			message := Message{}
			rmrMessage, err1 := message.PolicyMessage(strconv.FormatInt((int64(policytypeid)), 10), string(policyinstanceid), policyinstance.(string), "CREATE")
			if err1 != nil {
				a1.Logger.Error("error : %v", err1)
				return err1
			}
			a1.Logger.Debug("rmrMessage ", rmrMessage)
			isSent := rmr.RmrSendToXapp(rmrMessage, a1PolicyRequest)
			if isSent {
				a1.Logger.Debug("rmrSendToXapp : message sent")
			} else {
				a1.Logger.Error("rmrSendToXapp : message not sent")
			}
		}

	case "A1_EI_QUERY_ALL":
		a1.Logger.Debug("message recieved ", msg.Payload)
		resp, err := http.Get(ecsEiTypePath)
		if err != nil {
			a1.Logger.Error("Received error while fetching health info: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			a1.Logger.Warning("Received no reponse from A1-EI service1")
		}
		a1.Logger.Debug("response from A1-EI service : ", resp)

		defer resp.Body.Close()
		respByte, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			a1.Logger.Debug("error in response: %+v", respByte)
		}

		a1.Logger.Debug("response : %+v", string(respByte))

		isSent := rmr.RmrSendToXapp(string(respByte), a1EiQueryAllResp)
		if isSent {
			a1.Logger.Debug("rmrSendToXapp : message sent")
		} else {
			a1.Logger.Error("rmrSendToXapp : message not sent")
		}
	case "A1_EI_CREATE_JOB":
		payload := msg.Payload
		a1.Logger.Debug("message recieved : %s", payload)

		var result map[string]interface{}

		err := json.Unmarshal([]byte(payload), &result)
		if err != nil {
			a1.Logger.Error("Unmarshal error : %+v", err)
			return err
		}
		a1.Logger.Debug("Unmarshaled message recieved : %s ", result)

		jobIdStr := strconv.FormatInt((int64(result["job-id"].(float64))), 10)
		jsonReq, err := json.Marshal(result)
		if err != nil {
			a1.Logger.Error("marshal error : %v", err)
			return err
		}

		a1.Logger.Debug("url to send to :", ecsEiJobPath+jobIdStr)
		req, err := http.NewRequest(http.MethodPut, ecsEiJobPath+jobIdStr, bytes.NewBuffer(jsonReq))
		if err != nil {
			a1.Logger.Error("http error : %v", err)
			return err
		}

		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{}
		resp, err3 := client.Do(req)
		if err3 != nil {
			a1.Logger.Error("error:", err3)
			return err
		}

		defer resp.Body.Close()

		a1.Logger.Debug("response status : ", resp.StatusCode)
		if resp.StatusCode == 200 || resp.StatusCode == 201 {
			a1.Logger.Debug("received successful response for ei-job-id : ", jobIdStr)
			rmrData := fmt.Sprintf(jobCreationData, jobIdStr)
			a1.Logger.Debug("rmr_Data to send: ", rmrData)

			isSent := rmr.RmrSendToXapp(rmrData, a1EiCreateJobResp)
			if isSent {
				a1.Logger.Debug("rmrSendToXapp : message sent")
			} else {
				a1.Logger.Error("rmrSendToXapp : message not sent")
			}
		} else {
			a1.Logger.Warning("failed to create EIJOB ")
		}

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
	go rmr.rmrclient.Start(rmr)
	a1.Logger.Debug("Reciever started")
}
