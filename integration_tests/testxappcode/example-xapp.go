/*
==================================================================================
  Copyright (c) 2020 AT&T Intellectual Property.
  Copyright (c) 2020 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/
package main

import (
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

type A1TestXapp struct {
	msgChan  chan *xapp.RMRParams
	appReady bool
	rmrReady bool
}

type policyRequest struct {
	Operation        string      `json:"operation"`
	PolicyTypeID     int         `json:"policy_type_id"`
	PolicyInstanceID string      `json:"policy_instance_id"`
	Pay              interface{} `json:"payload"`
}

type policyRequestResponse struct {
	PolicyTypeID     int    `json:"policy_type_id"`
	PolicyInstanceID string `json:"policy_instance_id"`
	HandlerID        string `json:"handler_id"`
	Status           string `json:"status"`
}

func (e *A1TestXapp) handlePolicyReq(msg *xapp.RMRParams) {

	// unmarshal the request
	var dat policyRequest
	if err := json.Unmarshal(msg.Payload, &dat); err != nil {
		panic(err)
	}

	var status string
	switch dat.Operation {
	case "CREATE":
		status = "OK"
	case "DELETE":
		status = "DELETED"
	}

	// form the response
	res := &policyRequestResponse{
		dat.PolicyTypeID,
		dat.PolicyInstanceID,
		"test_receiver",
		status,
	}

	outgoing, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	params := &xapp.RMRParams{
		Mtype:   20011,
		Payload: outgoing,
	}

	if ok := xapp.Rmr.SendMsg(params); !ok {
		xapp.Logger.Info("Rmr.SendMsg failed ...")
		panic("Policy response failed to send")
	}
}

func (e *A1TestXapp) messageLoop() {
	for {
		fmt.Println("Waiting for message..")

		msg := <-e.msgChan
		defer xapp.Rmr.Free(msg.Mbuf)

		switch msg.Mtype {
		case 20010:
			e.handlePolicyReq(msg)
		default:
			panic("Unexpected message type!")
		}
	}
}

func (e *A1TestXapp) Consume(rp *xapp.RMRParams) (err error) {
	e.msgChan <- rp
	return
}

func (e *A1TestXapp) Run() {
	// Set MDC (read: name visible in the logs)
	xapp.Logger.SetMdc("test receiver xapp", "0.1.0")

	// Register various callback functions for application management
	xapp.SetReadyCB(func(d interface{}) { e.rmrReady = true }, true)

	go e.messageLoop()

	xapp.Run(e)
}

func NewA1TestXapp(appReady, rmrReady bool) *A1TestXapp {
	return &A1TestXapp{
		msgChan:  make(chan *xapp.RMRParams),
		rmrReady: rmrReady,
		appReady: appReady,
	}
}

func main() {
	NewA1TestXapp(true, false).Run()
}
