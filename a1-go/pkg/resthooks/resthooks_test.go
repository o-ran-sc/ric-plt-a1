/*
==================================================================================
  Copyright (c) 2021 Samsung

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
package resthooks

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var rh *Resthook
var sdlInst *SdlMock

func TestMain(m *testing.M) {
	sdlInst = new(SdlMock)

	sdlInst.On("GetAll", "A1m_ns").Return([]string{"a1.policy_instance.1006001.qos",
		"a1.policy_type.1006001",
		"a1.policy_type.20000",
		"a1.policy_inst_metadata.1006001.qos",
	}, nil)

	a1.Init()
	rh = createResthook(sdlInst)
	code := m.Run()
	os.Exit(code)
}

func TestGetAllPolicyType(t *testing.T) {
	resp := rh.GetAllPolicyType()
	assert.Equal(t, 2, len(resp))
}

func TestGetPolicyType(t *testing.T) {

	policyTypeId := models.PolicyTypeID(20001)

	resp := rh.GetPolicyType(policyTypeId)

	var policyTypeSchema models.PolicyTypeSchema
	name := "admission_control_policy_mine"
	policyTypeSchema.Name = &name
	policytypeid := int64(20001)
	policyTypeSchema.PolicyTypeID = &policytypeid
	description := "various parameters to control admission of dual connection"
	policyTypeSchema.Description = &description
	schema := `{"$schema": "http://json-schema.org/draft-07/schema#","type":"object","properties": {"enforce": {"type":"boolean","default":"true",},"window_length": {"type":        "integer","default":1,"minimum":1,"maximum":60,"description": "Sliding window length (in minutes)",},
"blocking_rate": {"type":"number","default":10,"minimum":1,"maximum":100,"description": "% Connections to block",},"additionalProperties": false,},}`
	policyTypeSchema.CreateSchema = schema
	key := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)

	//Setup Expectations
	sdlInst.On("Get", a1MediatorNs, policyTypeId).Return(map[string]interface{}{key: policyTypeSchema}, nil)

	assert.NotNil(t, resp)

	sdlInst.AssertExpectations(t)

}

func TestCreatePolicyType(t *testing.T) {
	var policyTypeId models.PolicyTypeID
	policyTypeId = 20001
	var policyTypeSchema models.PolicyTypeSchema
	name := "admission_control_policy_mine"
	policyTypeSchema.Name = &name
	policytypeid := int64(20001)
	policyTypeSchema.PolicyTypeID = &policytypeid
	description := "various parameters to control admission of dual connection"
	policyTypeSchema.Description = &description
	policyTypeSchema.CreateSchema = `{"$schema": "http://json-schema.org/draft-07/schema#","type":"object","properties": {"enforce": {"type":"boolean","default":"true",},"window_length": {"type":        "integer","default":1,"minimum":1,"maximum":60,"description": "Sliding window length (in minutes)",},
"blocking_rate": {"type":"number","default":10,"minimum":1,"maximum":100,"description": "% Connections to block",},"additionalProperties": false,},}`

	data, err := policyTypeSchema.MarshalBinary()
	a1.Logger.Debug("error : %+v ", err)
	a1.Logger.Debug("data : %+v ", data)
	key := a1PolicyPrefix + strconv.FormatInt(20001, 10)
	a1.Logger.Debug("key : %+v ", key)
	//Setup Expectations
	sdlInst.On("SetIfNotExists", a1MediatorNs, key, string(data)).Return(true, nil)

	errresp := rh.CreatePolicyType(policyTypeId, policyTypeSchema)
	//Data Assertion
	assert.Nil(t, errresp)
	//Mock Assertion :Behavioral
	sdlInst.AssertExpectations(t)
}

func TestCreatePolicyTypeInstance(t *testing.T) {
	var policyTypeId models.PolicyTypeID
	policyTypeId = 20001
	var policyInstanceID models.PolicyInstanceID
	policyInstanceID = "123456"
	httpBody := `{
		"enforce":true,
		"window_length":20,
	   "blocking_rate":20,
		"trigger_threshold":10
		}`
	instancekey := a1PolicyPrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
	data, _ := json.Marshal(httpBody)
	a1.Logger.Debug("Marshaled String : %+v", string(data))
	a1.Logger.Debug("key   : %+v", instancekey)

	instancearr := []interface{}{instancekey, string(data)}
	sdlInst.On("Set", "A1m_ns", instancearr).Return("CREATE", nil)
	errresp := rh.CreatePolicyInstance(policyTypeId, policyInstanceID, httpBody)

	assert.Nil(t, errresp)
	sdlInst.AssertExpectations(t)
}

type SdlMock struct {
	mock.Mock
}

func (s *SdlMock) GetAll(ns string) ([]string, error) {
	args := s.MethodCalled("GetAll", ns)
	return args.Get(0).([]string), nil
}

func (s *SdlMock) Get(ns string, keys []string) (map[string]interface{}, error) {
	a1.Logger.Debug("Get Called ")
	args := s.MethodCalled("Get", ns, keys)
	a1.Logger.Debug("keys :%+v", args.Get(1))
	policytypeid := int64(20001)

	policyTypeSchemaString := `{"name":"admission_control_policy_mine","description":"various parameters to control admission of dual connection","policy_type_id": 20001,"create_schema":{"$schema": "http://json-schema.org/draft-07/schema#","type":    "object","properties": {"enforce": {"type":    "boolean","default": "true"},"window_length": {"type":"integer","default":     1,"minimum":     1,"maximum":     60,"description": "Sliding window length (in minutes)"},"blocking_rate": {"type":        "number","default":     10,"minimum":     1,"maximum":     1001,"description": "% Connections to block"},"additionalProperties": false}}}`

	a1.Logger.Error(" policyTypeSchemaString %+v", policyTypeSchemaString)
	policyTypeSchema, _ := json.Marshal((policyTypeSchemaString))
	// a1.Logger.Error(" policyTypeSchema error %+v",  err)
	a1.Logger.Error(" policyTypeSchema %+v", string(policyTypeSchema))
	var p models.PolicyTypeSchema
	_ = json.Unmarshal([]byte(string(policyTypeSchemaString)), &p)
	a1.Logger.Error("unmarshalled  policyTypeSchema %+v", p.CreateSchema)
	key := a1PolicyPrefix + strconv.FormatInt((policytypeid), 10)
	a1.Logger.Error(" key for policy type %+v", key)
	mp := map[string]interface{}{key: string(policyTypeSchema)}
	a1.Logger.Error("Get Called and mp return %+v ", mp)
	return mp, nil
}

func (s *SdlMock) SetIfNotExists(ns string, key string, data interface{}) (bool, error) {
	args := s.MethodCalled("SetIfNotExists", ns, key, data)
	return args.Bool(0), args.Error(1)
}

func (s *SdlMock) Set(ns string, pairs ...interface{}) error {
	args := s.MethodCalled("Set", ns, pairs)
	return args.Error(1)
}
func (s *SdlMock) SetIf(ns string, key string, oldData, newData interface{}) (bool, error) {
	args := s.MethodCalled("SetIfNotExists", ns, key, oldData, newData)
	return args.Bool(0), args.Error(1)
}
