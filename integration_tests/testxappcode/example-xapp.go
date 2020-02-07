/*
==================================================================================
  Copyright (c) 2019 AT&T Intellectual Property.
  Copyright (c) 2019 Nokia

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
	//	"net/http"
)

// This could be defined in types.go
type ExampleXapp struct {
	msgChan chan *xapp.RMRParams
	//	stats    map[string]xapp.Counter
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

//func (e *ExampleXapp) handleRICExampleMessage(ranName string, r *xapp.RMRParams) {
//	// Just update metrics and echo the message back (update the message type)
//	e.stats["RICExampleMessageRx"].Inc()
//
//	r.Mtype = r.Mtype + 1
//	if ok := xapp.Rmr.SendMsg(r); !ok {
//		xapp.Logger.Info("Rmr.SendMsg failed ...")
//	}
//}

func (e *ExampleXapp) handlePolicyReq(msg *xapp.RMRParams) {

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

	fmt.Println("outgoings", string(outgoing))
	//params := &xapp.RMRParams{}

}

func (e *ExampleXapp) messageLoop() {
	for {
		fmt.Println("Waiting for message..")

		msg := <-e.msgChan
		defer xapp.Rmr.Free(msg.Mbuf)

		switch msg.Mtype {
		case 20010:
			e.handlePolicyReq(msg)
		default:
			fmt.Println("WHAAAT")
			fmt.Println(msg.Mtype)
		}
	}
}

func (e *ExampleXapp) Consume(rp *xapp.RMRParams) (err error) {
	e.msgChan <- rp
	return
}

//func (u *ExampleXapp) TestRestHandler(w http.ResponseWriter, r *http.Request) {
//	xapp.Logger.Info("TestRestHandler called!")
//}

//func (u *ExampleXapp) ConfigChangeHandler(f string) {
//	xapp.Logger.Info("Config file changed, do something meaningful!")
//}

//func (u *ExampleXapp) StatusCB() bool {
//	xapp.Logger.Info("Status callback called, do something meaningful!")
//	return true
//}

func (e *ExampleXapp) Run() {
	// Set MDC (read: name visible in the logs)
	xapp.Logger.SetMdc("test receiver xapp", "0.1.0")

	// Register various callback functions for application management
	xapp.SetReadyCB(func(d interface{}) { e.rmrReady = true }, true)
	//xapp.AddConfigChangeListener(e.ConfigChangeHandler)
	//xapp.Resource.InjectStatusCb(e.StatusCB)

	// Inject own REST handler for testing purpose
	//xapp.Resource.InjectRoute("/ric/v1/testing", e.TestRestHandler, "POST")

	go e.messageLoop()

	xapp.Logger.Info("About to run")

	xapp.Run(e)
}

//func GetMetricsOpts() []xapp.CounterOpts {
//	return []xapp.CounterOpts{
//		{Name: "RICIndicationsRx", Help: "The total number of RIC inidcation events received"},
//		{Name: "RICExampleMessageRx", Help: "The total number of RIC example messages received"},
//	}
//}

func NewExampleXapp(appReady, rmrReady bool) *ExampleXapp {
	//metrics := GetMetricsOpts()
	return &ExampleXapp{
		//		stats:    xapp.Metric.RegisterCounterGroup(metrics, "ExampleXapp"),
		msgChan:  make(chan *xapp.RMRParams),
		rmrReady: rmrReady,
		appReady: appReady,
	}
}

func main() {
	NewExampleXapp(true, false).Run()
}
