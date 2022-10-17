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
	a1PolicyRequest = 20010
	a1SourceName    = "service-ricplt-a1mediator-http"
)

type RmrSender struct {
	rmrclient *xapp.RMRClient
}

type IRmrSender interface {
	RmrSendToXapp(httpBodyString string) bool
}

func NewRMRSender() IRmrSender {
	RMRclient := xapp.NewRMRClientWithParams(&xapp.RMRClientParams{
		StatDesc: "",
		RmrData: xapp.PortData{
			Name:              "",
			MaxSize:           65534,
			ThreadType:        0,
			LowLatency:        false,
			FastAck:           false,
			MaxRetryOnFailure: 1,
		},
	})
	return &RmrSender{
		rmrclient: RMRclient,
	}
}

func (rmr *RmrSender) RmrSendToXapp(httpBodyString string) bool {

	params := &xapp.RMRParams{}
	params.Mtype = a1PolicyRequest
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
