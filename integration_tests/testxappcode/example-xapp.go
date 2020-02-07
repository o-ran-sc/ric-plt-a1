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

var delay int        // used for the delay receiver
var handlerID string // used for the delay receiver too
var doQuery bool     // used for the query receiver

type a1Receiver struct {
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

func (e *a1Receiver) sendMsgRetry(params *xapp.RMRParams) {
	// helper for rmr that handles retries and sleep
	retries := 0
	for { // just keep trying until it works
		if e.rmrReady { // we must wait for ready, else SendMsg will blow with a nullptr
			if ok := xapp.Rmr.SendMsg(params); ok {
				xapp.Logger.Info("Msg successfully sent after %d retries!", retries)
				return
			} else {
				retries++
			}
			//xapp.Logger.Info("Query failed to send...")
		} else {
			xapp.Logger.Info("rmr not ready...")
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}

func (e *a1Receiver) handlePolicyReq(msg *xapp.RMRParams) {

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

	/*
		WARNING:
		we want to use rts here. However, the current go xapp framework does not support
		1) using rts after changing the message type
		2) using rts after changing the payload where the payload does not fit in the original message (this can happen as the request from a1 may not actually have a payload!)
		So, for now, we do NOT use rts here. This means we need routing entries in the helm chart.
		If the go xapp framework improves here, the routing entries can be removed.

		UPDATE: rts is just plain broken
	*/
	params := &xapp.RMRParams{
		Mtype:   20011,
		Payload: outgoing,
	}

	if delay > 0 {
		xapp.Logger.Info("Xapp is sleeping...")
		time.Sleep(time.Duration(delay) * time.Second) // so much work to replicate python's time.sleep(5)...
	}

	e.sendMsgRetry(params)

	xapp.Logger.Info("Policy response sent!")
}

func (e *a1Receiver) sendQuery() {
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

	for {
		/* We do this in a loop here, because even when the query first works, it could be the case that
		        a1 does not even have the type yet, or there are no instances yet. In this integration test,
				we basically just keep pounding away so that eventually a1 returns the list this int test is looking for.
				A real xapp would NOT call the query in a loop like this.
		*/
		e.sendMsgRetry(params)
		xapp.Logger.Info("Query sent successfully")
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func (e *a1Receiver) messageLoop() {
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

// Consume: This named function is a required callback for e to use the xapp interface. it is called on all received rmr messages.
func (e *a1Receiver) Consume(rp *xapp.RMRParams) (err error) {
	e.msgChan <- rp
	return
}

func (e *a1Receiver) Run() {
	// Set MDC (read: name visible in the logs)
	xapp.Logger.SetMdc(handlerID, "0.1.0")

	// Register various callback functions for application management
	xapp.SetReadyCB(func(d interface{}) { e.rmrReady = true }, true)

	// start message loop. We cannot wait for e.rmrReady here since that doesn't get populated until Run() runs.
	go e.messageLoop()

	if doQuery {
		// we are in the query tester; kick off a loop that does that until it works
		go e.sendQuery()
	}

	xapp.Run(e)
}

func newA1Receiver(appReady, rmrReady bool) *a1Receiver {
	return &a1Receiver{
		msgChan:  make(chan *xapp.RMRParams),
		rmrReady: rmrReady,
		appReady: appReady,
	}
}

func main() {

	delay = 0
	if d, ok := os.LookupEnv("TEST_RCV_SEC_DELAY"); ok {
		delay, _ = strconv.Atoi(d)
	}

	handlerID = "test_receiver"
	if hid, ok := os.LookupEnv("HANDLER_ID"); ok {
		handlerID = hid
	}

	doQuery = false
	if _, ok := os.LookupEnv("DO_QUERY"); ok {
		doQuery = true
	}

	newA1Receiver(true, false).Run()
}
