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
        "errors"
        "fmt"
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
	}, nil).Once()
	rmrSenderInst = new(RmrSenderMock)
	a1.Init()
	rh = createResthook(sdlInst, rmrSenderInst)
	code := m.Run()
	os.Exit(code)
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

func TestHealthFail(t *testing.T) {
        keys := []string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",}
        sdlInst.On("GetAll", "A1m_ns").Return(keys,errors.New("Some Error")).Once()
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
        keys := []string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",}
        sdlInst.On("GetAll", "A1m_ns").Return(keys,nil).Once()
	resp := rh.GetAllPolicyType()
	assert.Equal(t, 2, len(resp))
}

func TestGetAllPolicyTypeFail(t *testing.T) {
        keys := []string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",}
        sdlInst.On("GetAll", "A1m_ns").Return(keys,errors.New("Some Error")).Once()
        resp := rh.GetAllPolicyType()
        assert.Equal(t, 0, len(resp))
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
	resp,err := rh.GetPolicyType(policyTypeId)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
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
	sdlInst.On("SetIfNotExists", a1MediatorNs, key, string(data)).Return(true, nil).Once()
	errresp := rh.CreatePolicyType(policyTypeId, policyTypeSchema)
	//Data Assertion
	assert.Nil(t, errresp)
	//Mock Assertion :Behavioral
	sdlInst.AssertExpectations(t)
}


func TestCreatePolicyType2(t *testing.T) {
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
	sdlInst.On("SetIfNotExists", a1MediatorNs, key, string(data)).Return(false, nil).Once()
        resp := rh.CreatePolicyType(policyTypeId, policyTypeSchema)
        //Data Assertion
        assert.NotNil(t, resp)
        //Mock Assertion :Behavioral
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

func TestGetPolicyInstanceFail(t *testing.T) {

        var policyTypeId models.PolicyTypeID
        policyTypeId = 0
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = ""
        httpBody := `{
                "enforce":true,
                "window_length":20,
           "blocking_rate":20,
                "trigger_threshold":10
                }`
        instancekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
        a1.Logger.Debug("httpBody String : %+v", httpBody)
        a1.Logger.Debug("key   : %+v", instancekey)
        var keys [1]string
        keys[0] = instancekey
        testmp1 := map[string]interface{}{instancekey: string(httpBody)}
        //Setup Expectations
        sdlInst.On("Get", a1MediatorNs, keys[:]).Return(testmp1, errors.New("Some Error"))

        resp, err := rh.GetPolicyInstance(policyTypeId, policyInstanceID)
        a1.Logger.Error("err : %+v", err)
        assert.Error(t, err)
        assert.Equal(t, 0, len(resp))
        sdlInst.AssertExpectations(t)
}

func TestGetAllPolicyIntances(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 20005
        keys := []string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",}
        sdlInst.On("GetAll", "A1m_ns").Return(keys,nil).Once()
        resp, err := rh.GetAllPolicyInstance(policyTypeId)
        fmt.Println("GETALL is ",resp)
        assert.NoError(t, err)
        assert.Equal(t, 2, len(resp))
}


func TestDeletePolicyType(t *testing.T) {

	policyTypeId := models.PolicyTypeID(20001)
	key := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
	var keys [1]string
	keys[0] = key

	//Setup Expectations
        sdlInst.On("GetAll", "A1m_ns").Return([]string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",},nil).Once()
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
	sdlInst.On("Get", a1MediatorNs, keys[:]).Return(httpBody,nil)
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

	rmrSenderInst.On("RmrSendToXapp", httpBodyString, 20010, int(policyTypeId)).Return(true)
	notificationDestinationkey := a1NotificationDestinationPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
	var notificationDestinationkeys [1]string
	notificationDestinationkeys[0] = notificationDestinationkey
	sdlInst.On("Remove", a1MediatorNs, notificationDestinationkeys[:]).Return(nil)
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
	rmrSenderInst.On("RmrSendToXapp", httpBodyString, 20017, -1).Return(true)
	errresp := rh.DataDelivery(instancedata)

	assert.Nil(t, errresp)
	sdlInst.AssertExpectations(t)
}
/*
func TestDataDeliveryFail(t *testing.T) {

        httpBody := `{
                "job":"1",
                "payload":"payload"
                }`
        var instancedata interface{}

        json.Unmarshal([]byte(httpBody), &instancedata)
        a1.Logger.Debug("Marshaled data : %+v", (instancedata))
        httpBodyString := `{"ei_job_id1":"1","payload":"payload"}`
        rmrSenderInst.On("RmrSendToXapp", httpBodyString, 20017, -1).Return(false)
        errresp := rh.DataDelivery(instancedata)

        assert.NotNil(t, errresp)
        sdlInst.AssertExpectations(t)
}
*/

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

func TestGetMetaDataFail(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 0
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = ""
        instanceMetadataKey := a1InstanceMetadataPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
        a1.Logger.Debug("key : %+v", instanceMetadataKey)
        var keys [1]string
        keys[0] = instanceMetadataKey
        policySchemaString := `{
                "created_at":"2022-11-02 10:30:20",
                "instance_status":"NOT IN EFFECT"
                }`
        sdlInst.On("Get", a1MediatorNs, keys[:]).Return(map[string]interface{}{instanceMetadataKey: policySchemaString}, errors.New("Some Error"))
        _, errresp := rh.getMetaData(policyTypeId, policyInstanceID)
        assert.NotNil(t, errresp)
}



func TestGetAllPolicyIntancesFail1(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 20009
        keys := []string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",}
        sdlInst.On("GetAll", "A1m_ns").Return(keys,nil).Once()
        resp, err := rh.GetAllPolicyInstance(policyTypeId)
        fmt.Println("GETALL is ",resp)
        fmt.Println("GETALL is ",err)
        assert.Equal(t, 0, len(resp))
}

func TestGetAllPolicyIntancesFail2(t *testing.T) {
        keys := []string{"",}
        sdlInst.On("GetAll", "A1m_ns").Return(keys, errors.New("Some Error")).Once()
        resp, err := rh.GetAllPolicyInstance(0)
        fmt.Println("GETALL is ",resp)
        fmt.Println("GETALL is ",err)
        assert.Error(t, err)
        assert.Equal(t, 0, len(resp))
}



func TestGetPolicyTypeFail(t *testing.T) {

        policyTypeId := models.PolicyTypeID(0)

        var policyTypeSchema models.PolicyTypeSchema
        name := "admission_control_policy_mine"
        policyTypeSchema.Name = &name
        policytypeid := int64(0)
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
        sdlInst.On("Get", a1MediatorNs, keys[:]).Return(map[string]interface{}{key: policyTypeSchema}, errors.New("Some Error")).Once()
        resp,err := rh.GetPolicyType(policyTypeId)
        assert.Nil(t, resp)
        assert.NotNil(t, err)
        //sdlInst.AssertExpectations(t)

}

func TestCreatePolicyTypeFail(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 0
        var policyTypeSchema models.PolicyTypeSchema
        name := "admission_control_policy_mine"
        policyTypeSchema.Name = &name
        policytypeid := int64(0)
        policyTypeSchema.PolicyTypeID = &policytypeid
        description := "various parameters to control admission of dual connection"
        policyTypeSchema.Description = &description
        policyTypeSchema.CreateSchema = `{"$schema": "http://json-schema.org/draft-07/schema#","type":"object","properties": {"enforce": {"type":"boolean","default":"true",},"window_length": {"type":        "integer","default":1,"minimum":1,"maximum":60,"description": "Sliding window length (in minutes)",},
"blocking_rate": {"type":"number","default":10,"minimum":1,"maximum":100,"description": "% Connections to block",},"additionalProperties": false,},}`

        data, err := policyTypeSchema.MarshalBinary()
        a1.Logger.Debug("error : %+v ", err)
        a1.Logger.Debug("data : %+v ", data)
        key := a1PolicyPrefix + strconv.FormatInt(0, 10)
        a1.Logger.Debug("key : %+v ", key)
        //Setup Expectations
        sdlInst.On("SetIfNotExists", a1MediatorNs, key, string(data)).Return(false, errors.New("Some Error")).Once()
        errresp := rh.CreatePolicyType(policyTypeId, policyTypeSchema)
        //Data Assertion
        assert.NotNil(t, errresp)
        //Mock Assertion :Behavioral
        //sdlInst.AssertExpectations(t)
}

func TestValidateInvalidJson(t *testing.T) {
         
         schemaString := `{"$schema": "http://json-schema.org/draft-07/schema#","type":"object","properties": {"enforce": {"type":"boolean","default":"true",},"window_length": {"type":        "integer","default":1,"minimum":1,"maximum":60,"description": "Sliding window length (in minutes)",},
"blocking_rate": {"type":"number","default":10,"minimum":1,"maximum":100,"description": "% Connections to block",},"additionalProperties": false,},}`
         httpBodyString := `{"blocking_rate"::20,"enforce"::true,"trigger_threshold":10,"window_length":20}`
         resp := validate(httpBodyString , schemaString) 
         assert.Equal(t, false, resp)
}

func TestValidateInvalidSchema(t *testing.T) {
         schemaString := `{"$schema": "http://json-schema.org/draft-07/sc"}`
         httpBodyString := `{"blocking_rate":20, 20:true, 30:10,"window_length":20}`
         resp := validate(httpBodyString , schemaString)
         assert.Equal(t, false, resp)
}

func TestValidateIntKeys(t *testing.T) {
         schemaString := `{"$schema": "http://json-schema.org/draft-07/sc"}`
         httpBodyString := `{20:20}`
         resp := validate(httpBodyString , schemaString)
         assert.Equal(t, false, resp)
}

func TestStorePolicyInstanceFail(t *testing.T) {
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = ""
        var policyTypeId models.PolicyTypeID
        policyTypeId = 0
        var httpBody = `{"enforce":true,"window_length":20,"blocking_rate":20,"trigger_threshold":10}`
        key := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
        var keys [1]string
        keys[0] = key
        //Setup Expectations
        sdlInst.On("Get", a1MediatorNs, keys[:]).Return(map[string]interface{}{key: httpBody}, errors.New("Some Error")).Once()
        notificationDestination :="abc"
        _,err := rh.storePolicyInstance(policyTypeId,policyInstanceID ,httpBody ,notificationDestination )
        assert.NotNil(t, err)
}

func TestStorePolicyInstanceMetadataFail(t *testing.T) {
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = ""
        var policyTypeId models.PolicyTypeID
        policyTypeId = 0
	sdlInst.On("Set", "A1m_ns", mock.Anything).Return(errors.New("Some Error")).Once()
        resp,err := rh.storePolicyInstanceMetadata(policyTypeId,policyInstanceID)
        assert.NotNil(t, err)
        assert.Equal(t, false, resp)
}

func TestStoreDeletedPolicyInstanceMetadataFail(t *testing.T) {
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = ""
        var policyTypeId models.PolicyTypeID
        policyTypeId = 0
        sdlInst.On("Set", "A1m_ns", mock.Anything).Return(errors.New("Some Error")).Once()
        err := rh.storeDeletedPolicyInstanceMetadata(policyTypeId,policyInstanceID,"")
        assert.NotNil(t, err)
        
}




func TestToString(t *testing.T) {
         var instance_map []interface{}
	 instancekey := "testkey"
         data := "testdata"
         instance_map = append(instance_map, instancekey, data)
         resp,err := toStringKeys(instance_map)
         assert.Nil(t, err)
         assert.NotNil(t, resp)
}

func TestToStringFail(t *testing.T) {
         m := make(map[interface{}]interface{})
         m[12] = "bar"
         _,err := toStringKeys(m)
         assert.NotNil(t, err)
}
/*
func TestToStringFail2(t *testing.T) {
         var keys [1]string
         keys[0] = ""
         _,err := toStringKeys(keys[:])
         assert.NotNil(t, err)
}
*/
func TestTypeValidityFail(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 123
        httpBody := `{
                "enforce":true,
                "window_length":20,
           "blocking_rate":20,
                "trigger_threshold":10
                }`
        instancekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
        a1.Logger.Debug("httpBody String : %+v", httpBody)
        a1.Logger.Debug("key   : %+v", instancekey)
        var keys [1]string
        keys[0] = instancekey
        testmp1 := map[string]interface{}{instancekey: string(httpBody)}
        //Setup Expectations
        sdlInst.On("Get", a1MediatorNs, keys[:]).Return(testmp1, errors.New("Some Error")).Once()
        err:= rh.typeValidity(policyTypeId)
        a1.Logger.Debug("err from get test  : %+v", err)
        assert.NotNil(t, err)
}

func TestDeletePolicyTypeFail(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 20009
        keys := []string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",}
        sdlInst.On("GetAll", "A1m_ns").Return(keys,errors.New("Some Error")).Once()
        err := rh.DeletePolicyType(policyTypeId)
        assert.NotNil(t, err)
}       

func TestDeletePolicyTypeFail2(t *testing.T) {

        policyTypeId := models.PolicyTypeID(20005)
        key := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
        var keys [1]string
        keys[0] = key

        //Setup Expectations
        sdlInst.On("GetAll", "A1m_ns").Return([]string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",},nil).Once()
        sdlInst.On("Remove", a1MediatorNs, keys[:]).Return(errors.New("Some Error")).Once()

        errresp := rh.DeletePolicyType(policyTypeId)

        assert.NotNil(t, errresp)
}

func TestDeletePolicyTypeFail3(t *testing.T) {

        policyTypeId := models.PolicyTypeID(20000)
        key := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
        var keys [1]string
        keys[0] = key

        //Setup Expectations
        sdlInst.On("GetAll", "A1m_ns").Return([]string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",},nil).Once()
        sdlInst.On("Remove", a1MediatorNs, keys[:]).Return(errors.New("Some Error")).Once()
        errresp := rh.DeletePolicyType(policyTypeId)
        assert.NotNil(t, errresp)
}

func TestGetPolicyInstanceStatusFail(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 0
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = ""
        httpBody := `{
                "created_at":"0001-01-01T00:00:00.000Z",
                "instance_status":"NOT IN EFFECT"
                }`
        instancekey := a1InstanceMetadataPrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
        a1.Logger.Debug("httpBody String : %+v", httpBody)
        a1.Logger.Debug("key   : %+v", instancekey)
        var keys [1]string
        keys[0] = instancekey
        sdlInst.On("Get", a1MediatorNs, keys[:]).Return(httpBody,nil)
        instancekey = a1HandlerPrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
        var instancekeys [1]string
        instancekeys[0] = instancekey
        instancearr := []interface{}{instancekey, "OK"}
        sdlInst.On("Get", a1MediatorNs, instancekeys[:]).Return(instancearr, errors.New("Some Error"))
        _, errresp := rh.GetPolicyInstanceStatus(policyTypeId, policyInstanceID)
        assert.NotNil(t, errresp)
}

func TestDeleteInstanceDataFail(t *testing.T) {
	var policyTypeId models.PolicyTypeID
	policyTypeId = 0
	var policyInstanceID models.PolicyInstanceID
	policyInstanceID = ""
        instancekey := a1InstancePrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
	var keys [1]string
	keys[0] = instancekey

        //Setup Expectations
        sdlInst.On("Remove", a1MediatorNs, keys[:]).Return(errors.New("Some Error"))

        errresp := rh.deleteInstancedata(policyTypeId,policyInstanceID)

        assert.NotNil(t, errresp)
}


func TestDeleteNotificationDestinationFail(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 0
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = ""
        notificationDestinationkey  := a1NotificationDestinationPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
        var keys [1]string
        keys[0] = notificationDestinationkey 

        //Setup Expectations
        sdlInst.On("Remove", a1MediatorNs, keys[:]).Return(errors.New("Some Error"))

        errresp := rh.deleteNotificationDestination(policyTypeId,policyInstanceID)

        assert.NotNil(t, errresp)
}

func TestDeleteMetadataFail(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 0
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = ""
        instanceMetadataKey := a1InstanceMetadataPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
        var keys [1]string
        keys[0] = instanceMetadataKey

        //Setup Expectations
        sdlInst.On("Remove", a1MediatorNs, keys[:]).Return(errors.New("Some Error"))

        errresp := rh.deleteMetadata(policyTypeId,policyInstanceID)

        assert.NotNil(t, errresp)
}


func TestGetPolicyInstanceStFail(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 200
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = "123"
        policyString := "testval"
        instancekey := a1HandlerPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
        a1.Logger.Debug("policyString String : %+v", policyString)
        a1.Logger.Debug("key from get test  : %+v", instancekey)
        var keys [1]string
        keys[0] = instancekey
        instancearr := []interface{}{instancekey, "NOK"}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(instancearr, errors.New("Some Error")).Once()
        resp, errresp := rh.getPolicyInstanceStatus(policyTypeId, policyInstanceID)
        a1.Logger.Debug("resp from get test  : %+v", resp)
        a1.Logger.Debug("errresp from get test  : %+v", errresp)
        assert.NotNil(t, errresp)
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
	sdlInst.On("SetIfNotExists", a1MediatorNs, instancekey, string(httpBody), string(data)).Return(true, nil)
        notificationDestinationkey := a1NotificationDestinationPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
        notificationDestination := "https://www.abc.com"
	notificationarr := []interface{}{notificationDestinationkey, string(notificationDestination)}
	sdlInst.On("Set", "A1m_ns", notificationarr).Return(nil)
          
	rmrSenderInst.On("RmrSendToXapp", "httpBodyString", 20010, int(policyTypeId)).Return(true)

	errresp := rh.CreatePolicyInstance(policyTypeId, policyInstanceID, instancedata, notificationDestination)

	assert.Nil(t, errresp)
}

func TestCreatePolicyTypeInstance2(t *testing.T) {
	var policyInstanceID models.PolicyInstanceID
	policyInstanceID = "123"
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
        sdlInst.On("Get", "A1m_ns", mock.Anything).Return(instancearr, nil).Once()

	metadatainstancekey := a1InstanceMetadataPrefix + strconv.FormatInt(20001, 10) + "." + string(policyInstanceID)
	creation_timestamp := time.Now()
	var metadatajson []interface{}
	metadatajson = append(metadatajson, map[string]string{"created_at": creation_timestamp.Format("2006-01-02 15:04:05"), "has_been_deleted": "False"})
	metadata, _ := json.Marshal(metadatajson)
	a1.Logger.Debug("Marshaled Metadata : %+v", string(metadata))
	a1.Logger.Debug("metadatainstancekey   : %+v", metadatainstancekey)
	metadatainstancearr := []interface{}{metadatainstancekey, string(metadata)}
	sdlInst.On("Set", "A1m_ns", metadatainstancearr).Return(nil)
	sdlInst.On("SetIfNotExists", a1MediatorNs, instancekey, string(httpBody), string(data)).Return(true, nil)
        notificationDestinationkey := a1NotificationDestinationPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
        notificationDestination := "https://www.abc.com"
	notificationarr := []interface{}{notificationDestinationkey, string(notificationDestination)}
	sdlInst.On("Set", "A1m_ns", notificationarr).Return(nil)
          
	rmrSenderInst.On("RmrSendToXapp", "httpBodyString", 20010, int(policyTypeId)).Return(true)

	errresp := rh.CreatePolicyInstance(policyTypeId, policyInstanceID, instancedata, notificationDestination)

	assert.Nil(t, errresp)
}

func TestCreatePolicyTypeInstanceFail(t *testing.T) {
	var policyInstanceID models.PolicyInstanceID
	policyInstanceID = "123456"
	var policyTypeId models.PolicyTypeID
	policyTypeId = 20001
        errresp := rh.CreatePolicyInstance(policyTypeId, policyInstanceID, "", "")
        assert.NotNil(t, errresp)
}

type SdlMock struct {
	mock.Mock
}

func (s *SdlMock) GetAll(ns string) ([]string, error) {
        args := s.MethodCalled("GetAll", ns)
        return args.Get(0).([]string), args.Error(1)
}

func (s *SdlMock) Get(ns string, keys []string) (map[string]interface{}, error) {
	a1.Logger.Debug("Get Called ")
	args := s.MethodCalled("Get", ns, keys)
	a1.Logger.Debug("ns :%+v", args.Get(0))
	var policySchemaString string
	var key string
	if keys[0] == "a1.policy_instance.20001.123456" {
	        policySchemaString = `{"enforce":true,"window_length":20,"blocking_rate":20,"trigger_threshold":10}`
		key = a1InstancePrefix + strconv.FormatInt(20001, 10) + "." + "123456"
	} else if keys[0] == "a1.policy_type.20001" {
		policySchemaString = `{"create_schema":{"$schema":"http://json-schema.org/draft-07/schema#","properties":{"additionalProperties":false,"blocking_rate":{"default":10,"description":"% Connections to block","maximum":1001,"minimum":1,"type":"number"},"enforce":{"default":"true","type":"boolean"},"window_length":{"default":1,"description":"Sliding window length (in minutes)","maximum":60,"minimum":1,"type":"integer"}},"type":"object"},"description":"various parameters to control admission of dual connection","name":"admission_control_policy_mine","policy_type_id":20001}`
		key = a1PolicyPrefix + strconv.FormatInt((20001), 10)
	} else if keys[0] == "a1.policy_inst_metadata.20001.123456" {
		policySchemaString = `{
			"created_at":"2022-11-02 10:30:20",
			"instance_status":"NOT IN EFFECT"
			}`
		key = a1InstanceMetadataPrefix + strconv.FormatInt(20001, 10) + "." + "123456"
	}
	a1.Logger.Debug(" policy SchemaString %+v", policySchemaString)
	policyTypeSchema, _ := json.Marshal((policySchemaString))
	a1.Logger.Debug(" policyTypeSchema %+v", string(policyTypeSchema))

	a1.Logger.Debug(" key for policy type %+v", key)
	mp := map[string]interface{}{key: string(policySchemaString)}
	a1.Logger.Debug("Get Called and mp return %+v ", mp)
	return mp, args.Error(1)
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

func (rmr *RmrSenderMock) RmrSendToXapp(httpBodyString string, mtype int, subid int) bool {
	if httpBodyString == `{"blocking_rate":20,"enforce":true,"trigger_threshold":10,"window_length":20}` {
		args := rmr.MethodCalled("RmrSendToXapp", httpBodyString, mtype, subid)
		return args.Bool(0)
	} else if httpBodyString == `{"ei_job_id":"1","payload":"payload"}` {
		args := rmr.MethodCalled("RmrSendToXapp", httpBodyString, mtype, subid)
		return args.Bool(0)
	}
	return true
}

func (s *SdlMock) Remove(ns string, keys []string) error {
	args := s.MethodCalled("Remove", ns, keys)
	return args.Error(0)
}
