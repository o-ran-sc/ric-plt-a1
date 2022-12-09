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
	"os"
	"strconv"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SdlMock struct {
	mock.Mock
}

var sdlInst *SdlMock
var pm *PolicyManager

func TestMain(m *testing.M) {
	sdlInst = new(SdlMock)
	a1.Init()
	pm = createPolicyManager(sdlInst)
	code := m.Run()
	os.Exit(code)
}
func TestSetPolicyInstance(t *testing.T) {
	var policyTypeId int
	policyTypeId = 20001
	var policyInstanceID int
	policyInstanceID = 123456
	var status string
	status = "OK"
	instancehandlerKey := a1HandlerPrefix + strconv.FormatInt(20001, 10) + "." + strconv.FormatInt(int64(policyInstanceID), 10)
	instancearr := []interface{}{instancehandlerKey, status}
	sdlInst.On("Set", "A1m_ns", instancearr).Return(nil)
	errresp := pm.SetPolicyInstanceStatus(policyTypeId, policyInstanceID, status)
	assert.NoError(t, errresp)
	sdlInst.AssertExpectations(t)
}

func TestGetAllPolicyIntances(t *testing.T) {
	var policyTypeId int
	policyTypeId = 20005
	sdlInst.On("GetAll", "A1m_ns").Return([]string{"a1.policy_instance.1006001.qos",
		"a1.policy_instance.20005.123456",
		"a1.policy_instance.20005.234567",
		"a1.policy_type.1006001",
		"a1.policy_type.20000",
		"a1.policy_inst_metadata.1006001.qos",
	}, nil)
	resp, err := pm.GetAllPolicyInstance(policyTypeId)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp))
}

func (s *SdlMock) Set(ns string, pairs ...interface{}) error {
	args := s.MethodCalled("Set", ns, pairs)
	return args.Error(0)
}

func (s *SdlMock) Get(ns string, keys []string) (map[string]interface{}, error) {
	a1.Logger.Error("Get Called ")
	return map[string]interface{}{}, nil
}

func (s *SdlMock) GetAll(ns string) ([]string, error) {
	args := s.MethodCalled("GetAll", ns)
	return args.Get(0).([]string), nil
}
