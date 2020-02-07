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
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"os"
	"strconv"
	"time"
)

var DELAY int         // used for the delay receiver
var HANDLER_ID string // used for the delay receiver too
var DO_QUERY bool     // used for the query receiver

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

type policyQuery struct {
	PolicyTypeID int `json:"policy_type_id"`
}

// helper for rmr that handles retries and sleep
func (e *A1TestXapp) sendMsgRetry(params *xapp.RMRParams) {

	// TODO! Use rts! take in a switchable flag

	for { // just keep trying until it works
		if e.rmrReady { // we must wait for ready, else SendMsg will blow with a nullptr
			if ok := xapp.Rmr.SendMsg(params); !ok {
				xapp.Logger.Info("Query failed to send...")
			} else {
				return
			}
		} else {
			xapp.Logger.Info("rmr not ready...")
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
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

	if DELAY > 0 {
		xapp.Logger.Info("Xapp is sleeping...")
		time.Sleep(time.Duration(DELAY) * time.Second) // so much work to replicate python's time.sleep(5)...
	}

	e.sendMsgRetry(params)

	xapp.Logger.Info("Policy response sent!")
}

func (e *A1TestXapp) sendQuery() {
	// form the query
	res := &policyQuery{
		1006001,
	}
	outgoing, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	params := &xapp.RMRParams{
		Mtype:   20012,
		Payload: outgoing,
	}

	for { // TODO, WHY WE DO MULTIPLE TIMES?
		e.sendMsgRetry(params)
		xapp.Logger.Info("Query sent successfully")
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func (e *A1TestXapp) messageLoop() {
	for {
		xapp.Logger.Info("Waiting for message..")

		msg := <-e.msgChan

		xapp.Logger.Info("Message received!")
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
	xapp.Logger.SetMdc(HANDLER_ID, "0.1.0")

	// Register various callback functions for application management
	xapp.SetReadyCB(func(d interface{}) { e.rmrReady = true }, true)

	// start message loop. We cannot wait for e.rmrReady here since that doesn't get populated until Run() runs.
	go e.messageLoop()

	if DO_QUERY {
		// we are in the query tester; kick off a loop that does that until it works
		go e.sendQuery()
	}

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

	DELAY = 0
	if d, ok := os.LookupEnv("TEST_RCV_SEC_DELAY"); ok {
		DELAY, _ = strconv.Atoi(d)
	}

	HANDLER_ID = "test_receiver"
	if hid, ok := os.LookupEnv("HANDLER_ID"); ok {
		HANDLER_ID = hid
	}

	DO_QUERY = false
	if _, ok := os.LookupEnv("DO_QUERY"); ok {
		DO_QUERY = true
	}

	NewA1TestXapp(true, false).Run()
}
