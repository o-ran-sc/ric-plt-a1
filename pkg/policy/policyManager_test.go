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
        "errors"
	"os"
	"strconv"
	"testing"
        "fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
        "gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
        a1InstanceMetadataPrefix        = "a1.policy_inst_metadata."
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
	policyInstanceID := "123456"
	var status string
	status = "OK"
	instancehandlerKey := a1HandlerPrefix + strconv.FormatInt(20001, 10) + "." + policyInstanceID
	instancearr := []interface{}{instancehandlerKey, status}
	sdlInst.On("Set", "A1m_ns", instancearr).Return(nil)
	errresp := pm.SetPolicyInstanceStatus(policyTypeId, policyInstanceID, status)
	assert.NoError(t, errresp)
	sdlInst.AssertExpectations(t)
}

func TestSetPolicyInstanceFail(t *testing.T) {
        var policyTypeId int
        policyTypeId = 0
        policyInstanceID := ""
        var status string
        status = "NOK"
        sdlInst.On("Set", "A1m_ns", mock.Anything).Return(errors.New("Some Error"))
        errresp := pm.SetPolicyInstanceStatus(policyTypeId, policyInstanceID, status)
        a1.Logger.Debug("err from set test  : %+v", errresp)
        assert.Error(t, errresp)
        sdlInst.AssertExpectations(t)
}

func TestGetPolicyInstance(t *testing.T) {

        var policyTypeId models.PolicyTypeID
        policyTypeId = 20001
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = "123456"
        policyString := "testval"
        typekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
        a1.Logger.Debug("policyString String : %+v", policyString)
        a1.Logger.Debug("typekey from get test  : %+v", typekey)
        var keys [1]string
        keys[0] = typekey
        testmp1 := map[string]interface{}{typekey: string(policyString)}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(testmp1, nil).Once()
        instancekey := a1InstancePrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
        a1.Logger.Debug("policyString String : %+v", policyString)
        a1.Logger.Debug("instancekey from get test  : %+v", instancekey)
        var instancekeys [1]string
        instancekeys[0] = instancekey
        testmp2 := map[string]interface{}{instancekey: string(policyString)}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", instancekeys[:]).Return(testmp2, nil).Once()
        resp, err := pm.GetPolicyInstance(policyTypeId, policyInstanceID)
        a1.Logger.Debug("resp from get test  : %+v", resp)
        a1.Logger.Debug("err from get test  : %+v", err)
        //sdlInst.AssertExpectations(t)
	assert.NoError(t, err)
        assert.NotNil(t, resp)
	sdlInst.AssertExpectations(t)
}

func TestGetPolicyInstanceFail1(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 20009
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = "123"
        policyString := "testval"
        typekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
        a1.Logger.Debug("policyString String : %+v", policyString)
        a1.Logger.Debug("typekey from get test  : %+v", typekey)
        var keys [1]string
        keys[0] = typekey
        testmp1 := map[string]interface{}{typekey: string(policyString)}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(testmp1,errors.New("Some Error")).Once()
        resp, err := pm.GetPolicyInstance(policyTypeId, policyInstanceID)
        a1.Logger.Debug("resp from get test  : %+v", resp)
        a1.Logger.Debug("err from get test  : %+v", err)
        //sdlInst.AssertExpectations(t)
        assert.Nil(t, resp)
        assert.NotNil(t, err)
        sdlInst.AssertExpectations(t)

}
func TestGetPolicyInstanceFail2(t *testing.T) {
        var policyTypeId models.PolicyTypeID
        policyTypeId = 20009
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = "123"
        policyString := "testval"
        typekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
        a1.Logger.Debug("policyString String : %+v", policyString)
        a1.Logger.Debug("typekey from get test  : %+v", typekey)
        var keys [1]string
        keys[0] = typekey
        testmp1 := map[string]interface{}{typekey: string(policyString)}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(testmp1,nil).Once()
        resp, err := pm.GetPolicyInstance(policyTypeId, policyInstanceID)
        a1.Logger.Debug("resp from get test  : %+v", resp)
        a1.Logger.Debug("err from get test  : %+v", err)
        //sdlInst.AssertExpectations(t)
        assert.Nil(t, resp)
        assert.NotNil(t, err)
        sdlInst.AssertExpectations(t)
}

func TestGetPolicyInstanceFail3(t *testing.T) {

        var policyTypeId models.PolicyTypeID
        policyTypeId = 20001
        var policyInstanceID models.PolicyInstanceID
        policyInstanceID = "123456"
        policyString := "testval"
        typekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
        a1.Logger.Debug("policyString String : %+v", policyString)
        a1.Logger.Debug("typekey from get test  : %+v", typekey)
        var keys [1]string
        keys[0] = typekey
        testmp1 := map[string]interface{}{typekey: string(policyString)}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(testmp1, nil).Once()
        instancekey := a1InstancePrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
        a1.Logger.Debug("policyString String : %+v", policyString)
        a1.Logger.Debug("instancekey from get test  : %+v", instancekey)
        var instancekeys [1]string
        instancekeys[0] = instancekey
        testmp2 := map[string]interface{}{instancekey: string(policyString)}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", instancekeys[:]).Return(testmp2,errors.New("Some Error")).Once()
        resp, err := pm.GetPolicyInstance(policyTypeId, policyInstanceID)
        a1.Logger.Debug("resp from get test  : %+v", resp)
        a1.Logger.Debug("err from get test  : %+v", err)
        //sdlInst.AssertExpectations(t)
        assert.NotNil(t, err)
        sdlInst.AssertExpectations(t)
}

func TestGetPolicyInstanceStatus(t *testing.T) {
	var policyTypeId int
	policyTypeId = 20001
	policyInstanceID := "123456"
        policyString := "testval"
        instancekey := a1HandlerPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + policyInstanceID
        a1.Logger.Debug("policyString String : %+v", policyString)
        a1.Logger.Debug("key from get test  : %+v", instancekey)
        var keys [1]string
        keys[0] = instancekey
        instancearr := []interface{}{instancekey, "OK"}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(instancearr, nil).Once()
        resp, errresp := pm.GetPolicyInstanceStatus(policyTypeId, policyInstanceID)
        a1.Logger.Debug("resp from get test  : %+v", resp)
        a1.Logger.Debug("errresp from get test  : %+v", errresp)
        assert.Nil(t, errresp)
        assert.NotNil(t, resp)
        sdlInst.AssertExpectations(t)
}

func TestGetPolicyInstanceStatusFail(t *testing.T) {
        var policyTypeId int
        policyTypeId = 200
        policyInstanceID := "123"
        policyString := "testval"
        instancekey := a1HandlerPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + policyInstanceID
        a1.Logger.Debug("policyString String : %+v", policyString)
        a1.Logger.Debug("key from get test  : %+v", instancekey)
        var keys [1]string
        keys[0] = instancekey
        instancearr := []interface{}{instancekey, "NOK"}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(instancearr, errors.New("Some Error")).Once()
        resp, errresp := pm.GetPolicyInstanceStatus(policyTypeId, policyInstanceID)
        a1.Logger.Debug("resp from get test  : %+v", resp)
        a1.Logger.Debug("errresp from get test  : %+v", errresp)
        assert.NotNil(t, errresp)
        sdlInst.AssertExpectations(t)
}

func TestSendPolicyStatusNotification(t *testing.T) {
	var policyTypeId int
	policyTypeId = 20001
	policyInstanceID := "123456"
        notificationDestinationkey := a1NotificationDestinationPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + fmt.Sprint(policyInstanceID)
        urlString := "http://www.abc.com"
        a1.Logger.Debug("Test: Calling key from get test  : %+v", notificationDestinationkey)
        keys := [1]string{notificationDestinationkey}
	var status string
        status = "OK"
        notDes := map[string]interface{}{notificationDestinationkey: string(urlString)}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(notDes, nil).Once()
        instancekey := a1HandlerPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + policyInstanceID
        var instancekeys [1]string
        instancekeys[0] = instancekey
        instancearr := []interface{}{instancekey, "OK"}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", instancekeys[:]).Return(instancearr, nil).Once()
        err := pm.SendPolicyStatusNotification(policyTypeId,policyInstanceID,notificationDestinationkey,status)
        assert.Nil(t, err)
        sdlInst.AssertExpectations(t)
}

func TestSendPolicyStatusNotificationFail1(t *testing.T) {
        var policyTypeId int
        policyTypeId = 20000
        policyInstanceID := "12345"
        notificationDestinationkey := a1NotificationDestinationPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + fmt.Sprint(policyInstanceID)
        urlString := "www.abc.com"
        a1.Logger.Debug("Test: Calling key from get test  : %+v", notificationDestinationkey)
        keys := [1]string{notificationDestinationkey}
        var status string
        status = "OK"
        notDes := map[string]interface{}{notificationDestinationkey: string(urlString)}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(notDes, errors.New("Some Error")).Once()
        err := pm.SendPolicyStatusNotification(policyTypeId,policyInstanceID,notificationDestinationkey,status)
        assert.NotNil(t, err)
        sdlInst.AssertExpectations(t)
}

func TestSendPolicyStatusNotificationFail2(t *testing.T) {
        var policyTypeId int
        policyTypeId = 20001
        policyInstanceID := "123456"
        notificationDestinationkey := a1NotificationDestinationPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + fmt.Sprint(policyInstanceID)
        urlString := "http://www.abc.com"
        a1.Logger.Debug("Test: Calling key from get test  : %+v", notificationDestinationkey)
        keys := [1]string{notificationDestinationkey}
        var status string
        status = "OK"
        notDes := map[string]interface{}{notificationDestinationkey: string(urlString)}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", keys[:]).Return(notDes,nil).Once()
        instancekey := a1HandlerPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + policyInstanceID
        var instancekeys [1]string
        instancekeys[0] = instancekey
        instancearr := []interface{}{instancekey, "NOK"}
        //Setup Expectations
        sdlInst.On("Get", "A1m_ns", instancekeys[:]).Return(instancearr, errors.New("Some Error")).Once()
        err := pm.SendPolicyStatusNotification(policyTypeId,policyInstanceID,notificationDestinationkey,status)
        assert.NotNil(t, err)
        sdlInst.AssertExpectations(t)
}

func TestGetAllPolicyIntances(t *testing.T) {
	var policyTypeId int
	policyTypeId = 20005
        keys := []string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",}
	sdlInst.On("GetAll", "A1m_ns").Return(keys,nil).Once()
	resp, err := pm.GetAllPolicyInstance(policyTypeId)
        fmt.Println("GETALL is ",resp)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp))
}

func TestGetAllPolicyIntancesFail(t *testing.T) {
        var policyTypeId int
        policyTypeId = 20009
        keys := []string{"a1.policy_instance.1006001.qos","a1.policy_instance.20005.123456","a1.policy_instance.20005.234567","a1.policy_type.1006001","a1.policy_type.20000","a1.policy_inst_metadata.1006001.qos",}
	sdlInst.On("GetAll", "A1m_ns").Return(keys,nil).Once()
        resp, err := pm.GetAllPolicyInstance(policyTypeId)
        fmt.Println("GETALL is ",resp)
        assert.Error(t, err)
        assert.Equal(t, 0, len(resp))
}


func TestGetAllPolicyIntancesFail2(t *testing.T) {
        keys := []string{"",}
        sdlInst.On("GetAll", "A1m_ns").Return(keys, errors.New("Some Error")).Once()
        resp, err := pm.GetAllPolicyInstance(0)
        fmt.Println("GETALL is ",resp)
        fmt.Println("GETALL is ",err)
        assert.Error(t, err)
        assert.Equal(t, 0, len(resp))
}


func (s *SdlMock) Set(ns string, pairs ...interface{}) error {
	args := s.MethodCalled("Set", ns, pairs)
	return args.Error(0)
}


func (s *SdlMock) Get(ns string, keys []string) (map[string]interface{}, error) {
        a1.Logger.Debug("Mock:Get Called ")
        args := s.MethodCalled("Get", ns, keys)
        a1.Logger.Debug("Mock: ns :%+v", args.Get(0))
        var policySchemaString string
        var key string
        a1.Logger.Debug("Mock: Key[0] is:%+v", keys[0])
        if keys[0] == "a1.policy_instance.20001.123456" {
                policySchemaString = "testval"
                key = a1InstancePrefix + strconv.FormatInt(20001, 10) + "." + "123456"
        } else if keys[0] == "a1.policy_type.20001" {
                policySchemaString = "testval"
                key = a1PolicyPrefix + strconv.FormatInt((int64(20001)), 10)
        } else if keys[0] == "a1.policy_notification_destination.20001.123456" {
                policySchemaString = "http://www.xyz.com"
                key = a1NotificationDestinationPrefix + strconv.FormatInt((int64(20001)), 10) + "." + "123456"
        } else if keys[0] == "a1.policy_handler.20001.123456" {
                policySchemaString = "testval"
                key = a1HandlerPrefix + strconv.FormatInt(20001, 10) + "." + "123456"
        } else if keys[0] == "a1.policy_inst_metadata.20001.123456" {
                policySchemaString = "testval"
                key = a1InstanceMetadataPrefix + strconv.FormatInt(20001, 10) + "." + "123456"
        } else if keys[0] == "a1.policy_notification_destination.20000.12345" {
                policySchemaString = "www.xyz.com"
                key = a1NotificationDestinationPrefix + strconv.FormatInt((int64(20000)), 10) + "." + "12345"
        }
        a1.Logger.Debug(" Mock: policy SchemaString %+v", policySchemaString)
        a1.Logger.Debug(" Mock: key for policy type %+v", key)
        mp := map[string]interface{}{key: string(policySchemaString)}
        a1.Logger.Debug("Mock: Get Called and mp return %+v ", mp)
        return mp, args.Error(1)
}

func (s *SdlMock) GetAll(ns string) ([]string, error) {
	args := s.MethodCalled("GetAll", ns)
	return args.Get(0).([]string), args.Error(1)
}

func (s *SdlMock) RemoveAll(ns string) error {
        //args := s.MethodCalled("RemoveAll", ns)
        return nil
}

