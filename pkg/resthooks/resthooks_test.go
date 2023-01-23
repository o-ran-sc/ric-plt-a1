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
	"time"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
        "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RmrSenderMock struct {
	mock.Mock
}

var rh *Resthook
var sdlInst *SdlMock
var rmrSenderInst *RmrSenderMock

func TestMain(m *testing.M) {
	sdlInst = new(SdlMock)

	sdlInst.On("GetAll", "A1m_ns").Return([]string{"a1.policy_instance.1006001.qos",
		"a1.policy_instance.20005.123456",
		"a1.policy_instance.20005.234567",
		"a1.policy_type.1006001",
		"a1.policy_type.20000",
		"a1.policy_inst_metadata.1006001.qos",
	}, nil)
	rmrSenderInst = new(RmrSenderMock)
	a1.Init()
	rh = createResthook(sdlInst, rmrSenderInst)
	code := m.Run()
	os.Exit(code)
}

func (s *SdlMock) GetAll(ns string) ([]string, error) {
	args := s.MethodCalled("GetAll", ns)
	return args.Get(0).([]string), nil
}


func TestHealth(t *testing.T) {
	resp := rh.GetA1Health()
	if resp == true {
		a1.Logger.Debug("A1 is healthy ")
		assert.Equal(t, true, resp)
	} else {
		a1.Logger.Debug("A1 is unhealthy")
		assert.Equal(t, false, resp)
	}
}

func TestGetAllPolicyType(t *testing.T) {
	resp := rh.GetAllPolicyType()
	assert.Equal(t, 2, len(resp))
}

func TestGetPolicyType(t *testing.T) {

	policyTypeId := models.PolicyTypeID(20001)

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
	var keys [1]string
	keys[0] = key
	//Setup Expectations
	sdlInst.On("Get", a1MediatorNs, keys[:]).Return(map[string]interface{}{key: policyTypeSchema}, nil)
	resp := rh.GetPolicyType(policyTypeId)
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
	var policyInstanceID models.PolicyInstanceID
	policyInstanceID = "123456"
	var httpBody = `{"enforce":true,"window_length":20,"blocking_rate":20,"trigger_threshold":10}`
	instancekey := a1InstancePrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
	var policyTypeId models.PolicyTypeID
	policyTypeId = 20001

	var instancedata map[string]interface{}

	json.Unmarshal([]byte(httpBody), &instancedata)

	data, _ := json.Marshal(instancedata)
	a1.Logger.Debug("Marshaled data : %+v", string(data))
	a1.Logger.Debug("instancekey   : %+v", instancekey)
	instancearr := []interface{}{instancekey, string(data)}
	sdlInst.On("Set", "A1m_ns", instancearr).Return(nil)

	metadatainstancekey := a1InstanceMetadataPrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
	creation_timestamp := time.Now()
	var metadatajson []interface{}
	metadatajson = append(metadatajson, map[string]string{"created_at": creation_timestamp.Format("2006-01-02 15:04:05"), "has_been_deleted": "False"})
	metadata, _ := json.Marshal(metadatajson)
	a1.Logger.Debug("Marshaled Metadata : %+v", string(metadata))
	a1.Logger.Debug("metadatainstancekey   : %+v", metadatainstancekey)
	metadatainstancearr := []interface{}{metadatainstancekey, string(metadata)}
	sdlInst.On("Set", "A1m_ns", metadatainstancearr).Return(nil)
	rmrSenderInst.On("RmrSendToXapp", "httpBodyString", 20010).Return(true)

	errresp := rh.CreatePolicyInstance(policyTypeId, policyInstanceID, instancedata)

	assert.Nil(t, errresp)
	sdlInst.AssertExpectations(t)
}

func TestGetPolicyInstance(t *testing.T) {

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
	instancekey := a1InstancePrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
	a1.Logger.Debug("httpBody String : %+v", httpBody)
	a1.Logger.Debug("key   : %+v", instancekey)
	var keys [1]string
	keys[0] = instancekey
	//Setup Expectations
	sdlInst.On("Get", a1MediatorNs, keys[:]).Return(httpBody, nil)

	resp, err := rh.GetPolicyInstance(policyTypeId, policyInstanceID)
	a1.Logger.Error("err : %+v", err)
	assert.NotNil(t, resp)

	sdlInst.AssertExpectations(t)
}

func TestGetAllPolicyIntances(t *testing.T) {
	var policyTypeId models.PolicyTypeID
	policyTypeId = 20005
	resp, err := rh.GetAllPolicyInstance(policyTypeId)
	a1.Logger.Error("err : %+v", err)
	assert.Equal(t, 2, len(resp))
}

func TestDeletePolicyType(t *testing.T) {

	policyTypeId := models.PolicyTypeID(20001)
	key := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
	var keys [1]string
	keys[0] = key

	//Setup Expectations
	sdlInst.On("Remove", a1MediatorNs, keys[:]).Return(nil)

	errresp := rh.DeletePolicyType(policyTypeId)

	assert.Nil(t, errresp)
	sdlInst.AssertExpectations(t)
}

func TestGetPolicyInstanceStatus(t *testing.T) {
	var policyTypeId models.PolicyTypeID
	policyTypeId = 20001
	var policyInstanceID models.PolicyInstanceID
	policyInstanceID = "123456"
	httpBody := `{
		"created_at":"0001-01-01T00:00:00.000Z",
		"instance_status":"NOT IN EFFECT"
		}`
	instancekey := a1InstanceMetadataPrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
	a1.Logger.Debug("httpBody String : %+v", httpBody)
	a1.Logger.Debug("key   : %+v", instancekey)
	var keys [1]string
	keys[0] = instancekey
	sdlInst.On("Get", a1MediatorNs, keys[:]).Return(httpBody)
	instancekey = a1HandlerPrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
	var instancekeys [1]string
	instancekeys[0] = instancekey
	instancearr := []interface{}{instancekey, "OK"}
	sdlInst.On("Get", a1MediatorNs, instancekeys[:]).Return(instancearr, nil)
	resp, errresp := rh.GetPolicyInstanceStatus(policyTypeId, policyInstanceID)

	assert.Nil(t, errresp)
	assert.NotNil(t, resp)
	sdlInst.AssertExpectations(t)
}

func TestDeletePolicyInstance(t *testing.T) {
	var policyTypeId models.PolicyTypeID
	policyTypeId = 20001
	var policyInstanceID models.PolicyInstanceID
	policyInstanceID = "123456"
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
	var policytypekeys [1]string
	policytypekeys[0] = key

	sdlInst.On("Get", a1MediatorNs, policytypekeys[:]).Return(map[string]interface{}{key: policyTypeSchema}, nil)

	httpBody := `{
		"enforce":true,
		"window_length":20,
	   "blocking_rate":20,
		"trigger_threshold":10
		}`
	instancekey := a1InstancePrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
	var instancekeys [1]string
	instancekeys[0] = instancekey

	sdlInst.On("Get", a1MediatorNs, instancekeys[:]).Return(httpBody, nil)

	var instanceMetadataKeys [1]string
	instanceMetadataKey := a1InstanceMetadataPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
	instanceMetadataKeys[0] = instanceMetadataKey
	httpBody = `{
		"created_at":"2022-11-02 10:30:20",
			"instance_status":"NOT IN EFFECT"
		}`

	sdlInst.On("Get", a1MediatorNs, instanceMetadataKeys[:]).Return(httpBody, nil)

	sdlInst.On("Remove", a1MediatorNs, instanceMetadataKeys[:]).Return(nil)

	var metadatainstancekeys [1]string
	metadatainstancekeys[0] = instancekey

	sdlInst.On("Remove", a1MediatorNs, metadatainstancekeys[:]).Return(nil)

	metadatainstancekey := a1InstanceMetadataPrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
	deleted_timestamp := time.Now()
	var metadatajson interface{}
	metadatajson = map[string]string{"created_at": "2022-11-02 10:30:20", "deleted_at": deleted_timestamp.Format("2006-01-02 15:04:05"), "has_been_deleted": "True"}
	metadata, _ := json.Marshal(metadatajson)
	metadatainstancearr := []interface{}{metadatainstancekey, string(metadata)}

	sdlInst.On("Set", "A1m_ns", metadatainstancearr).Return(nil)

	httpBodyString := `{"operation":"DELETE","payload":"","policy_instance_id":"123456","policy_type_id":"20001"}`

	rmrSenderInst.On("RmrSendToXapp", httpBodyString, 20010).Return(true)

	errresp := rh.DeletePolicyInstance(policyTypeId, policyInstanceID)

	assert.Nil(t, errresp)
	sdlInst.AssertExpectations(t)
}
func TestDataDelivery(t *testing.T) {

	httpBody := `{
		"job":"1",
		"payload":"payload"
		}
		`
	var instancedata interface{}

	json.Unmarshal([]byte(httpBody), &instancedata)
	a1.Logger.Debug("Marshaled data : %+v", (instancedata))
	httpBodyString := `{"ei_job_id":"1","payload":"payload"}`
	rmrSenderInst.On("RmrSendToXapp", httpBodyString, 20017).Return(true)
	errresp := rh.DataDelivery(instancedata)

	assert.Nil(t, errresp)
	sdlInst.AssertExpectations(t)
}

func TestGetMetaData(t *testing.T) {
	var policyTypeId models.PolicyTypeID
	policyTypeId = 20001
	var policyInstanceID models.PolicyInstanceID
	policyInstanceID = "123456"
	instanceMetadataKey := a1InstanceMetadataPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
	a1.Logger.Debug("key : %+v", instanceMetadataKey)

	var keys [1]string
	keys[0] = instanceMetadataKey

	policySchemaString := `{
		"created_at":"2022-11-02 10:30:20",
		"instance_status":"NOT IN EFFECT"
		}`

	sdlInst.On("Get", a1MediatorNs, keys[:]).Return(map[string]interface{}{instanceMetadataKey: policySchemaString}, nil)

	resp, errresp := rh.getMetaData(policyTypeId, policyInstanceID)

	assert.Nil(t, errresp)
	assert.NotNil(t, resp)
	sdlInst.AssertExpectations(t)
}

type SdlMock struct {
	mock.Mock
}


func (s *SdlMock) Get(ns string, keys []string) (map[string]interface{}, error) {
	a1.Logger.Debug("Get Called ")
	args := s.MethodCalled("Get", ns, keys)
	a1.Logger.Debug("ns :%+v", args.Get(0))
	policytypeid := int64(20001)
	policyInstanceID := "123456"
	var policySchemaString string
	var key string
	if keys[0] == "a1.policy_instance.20001.123456" {
		policySchemaString = `{
			"enforce":true,
			"window_length":20,
		   "blocking_rate":20,
			"trigger_threshold":10
			}`
		key = a1InstancePrefix + strconv.FormatInt(policytypeid, 10) + "." + string(policyInstanceID)
	} else if keys[0] == "a1.policy_type.20001" {
		policySchemaString = `{"create_schema":{"$schema":"http://json-schema.org/draft-07/schema#","properties":{"additionalProperties":false,"blocking_rate":{"default":10,"description":"% Connections to block","maximum":1001,"minimum":1,"type":"number"},"enforce":{"default":"true","type":"boolean"},"window_length":{"default":1,"description":"Sliding window length (in minutes)","maximum":60,"minimum":1,"type":"integer"}},"type":"object"},"description":"various parameters to control admission of dual connection","name":"admission_control_policy_mine","policy_type_id":20001}`
		key = a1PolicyPrefix + strconv.FormatInt((policytypeid), 10)
	} else if keys[0] == "a1.policy_inst_metadata.20001.123456" {
		policySchemaString = `{
			"created_at":"2022-11-02 10:30:20",
			"instance_status":"NOT IN EFFECT"
			}`
		key = a1InstanceMetadataPrefix + strconv.FormatInt(policytypeid, 10) + "." + string(policyInstanceID)
	}
	a1.Logger.Debug(" policy SchemaString %+v", policySchemaString)
	policyTypeSchema, _ := json.Marshal((policySchemaString))
	a1.Logger.Debug(" policyTypeSchema %+v", string(policyTypeSchema))

	a1.Logger.Debug(" key for policy type %+v", key)
	mp := map[string]interface{}{key: string(policySchemaString)}
	a1.Logger.Debug("Get Called and mp return %+v ", mp)
	return mp, nil
}
func (s *SdlMock) SetIfNotExists(ns string, key string, data interface{}) (bool, error) {
	args := s.MethodCalled("SetIfNotExists", ns, key, data)
	return args.Bool(0), args.Error(1)
}

func (s *SdlMock) Set(ns string, pairs ...interface{}) error {
	args := s.MethodCalled("Set", ns, pairs)
	return args.Error(0)
}
func (s *SdlMock) SetIf(ns string, key string, oldData, newData interface{}) (bool, error) {
	args := s.MethodCalled("SetIfNotExists", ns, key, oldData, newData)
	return args.Bool(0), args.Error(1)
}

func (rmr *RmrSenderMock) RmrSendToXapp(httpBodyString string, mtype int) bool {
	if httpBodyString == `{"blocking_rate":20,"enforce":true,"trigger_threshold":10,"window_length":20}` {
		args := rmr.MethodCalled("RmrSendToXapp", httpBodyString, mtype)
		return args.Bool(0)
	} else if httpBodyString == `{"ei_job_id":"1","payload":"payload"}` {
		args := rmr.MethodCalled("RmrSendToXapp", httpBodyString, mtype)
		return args.Bool(0)
	}
	return true
}

func (s *SdlMock) Remove(ns string, keys []string) error {
	args := s.MethodCalled("Remove", ns, keys)
	return args.Error(0)
}
