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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/rmr"
	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v2"
)

const (
	a1PolicyPrefix           = "a1.policy_type."
	a1MediatorNs             = "A1m_ns"
	a1InstancePrefix         = "a1.policy_instance."
	a1InstanceMetadataPrefix = "a1.policy_inst_metadata."
	a1HandlerPrefix          = "a1.policy_handler."
)

var typeAlreadyError = errors.New("Policy Type already exists")
var InstanceAlreadyError = errors.New("Policy Instance already exists")
var typeMismatchError = errors.New("Policytype Mismatch")
var invalidJsonSchema = errors.New("Invalid Json ")
var policyInstanceNotFoundError = errors.New("Policy Instance Not Found")
var policyTypeNotFoundError = errors.New("Policy Type Not Found")
var policyTypeCanNotBeDeletedError = errors.New("tried to delete a type that isn't empty")

func (rh *Resthook) CanPolicyTypeBeDeleted(err error) bool {
	return err == policyTypeCanNotBeDeletedError
}

func (rh *Resthook) IsPolicyTypePresent(err error) bool {
	return err == policyTypeNotFoundError
}

func (rh *Resthook) IsPolicyInstanceNotFound(err error) bool {
	return err == policyInstanceNotFoundError
}

func (rh *Resthook) IsTypeAlready(err error) bool {
	return err == typeAlreadyError
}
func (rh *Resthook) IsInstanceAlready(err error) bool {
	return err == InstanceAlreadyError
}
func (rh *Resthook) IsTypeMismatch(err error) bool {
	return err == typeMismatchError
}

func (rh *Resthook) IsValidJson(err error) bool {
	return err == invalidJsonSchema
}
func NewResthook() *Resthook {
	return createResthook(sdlgo.NewSyncStorage(), rmr.NewRMRSender())
}

func createResthook(sdlInst iSdl, rmrSenderInst rmr.IRmrSender) *Resthook {
	return &Resthook{
		db:             sdlInst,
		iRmrSenderInst: rmrSenderInst,
	}
}

func (rh *Resthook) GetAllPolicyType() []models.PolicyTypeID {

	var policyTypeIDs []models.PolicyTypeID

	keys, err := rh.db.GetAll("A1m_ns")

	if err != nil {
		a1.Logger.Error("error in retrieving policy. err: %v", err)
		return policyTypeIDs
	}
	a1.Logger.Debug("keys : %+v", keys)

	for _, key := range keys {
		if strings.HasPrefix(strings.TrimLeft(key, " "), a1PolicyPrefix) {
			pti := strings.Split(strings.Trim(key, " "), a1PolicyPrefix)[1]
			ptii, _ := strconv.ParseInt(pti, 10, 64)
			policyTypeIDs = append(policyTypeIDs, models.PolicyTypeID(ptii))
		}
	}

	a1.Logger.Debug("return : %+v", policyTypeIDs)
	return policyTypeIDs
}

func (rh *Resthook) GetPolicyType(policyTypeId models.PolicyTypeID) *models.PolicyTypeSchema {
	a1.Logger.Debug("GetPolicyType1")

	var policytypeschema *models.PolicyTypeSchema
	var keys [1]string

	key := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
	keys[0] = key

	a1.Logger.Debug("key : %+v", key)

	valmap, err := rh.db.Get(a1MediatorNs, keys[:])

	a1.Logger.Debug("policytype map : %+v", valmap)

	if len(valmap) == 0 {
		a1.Logger.Error("policy type Not Present for policyid : %v", policyTypeId)
		return policytypeschema
	}

	if err != nil {
		a1.Logger.Error("error in retrieving policy type. err: %v", err)
		return nil
	}

	if valmap[key] == nil {
		a1.Logger.Error("policy type Not Present for policyid : %v", policyTypeId)
		return policytypeschema
	}

	a1.Logger.Debug("keysmap : %+v", valmap[key])

	var item models.PolicyTypeSchema
	valStr := fmt.Sprint(valmap[key])

	a1.Logger.Debug("Policy type for %+v :  %+v", key, valStr)
	valkey := "`" + valStr + "`"
	valToUnmarshall, err := strconv.Unquote(valkey)
	if err != nil {
		a1.Logger.Error("unquote error : %+v", err)
		return nil
	}

	a1.Logger.Debug("Policy type for %+v :  %+v", key, string(valToUnmarshall))

	errunm := json.Unmarshal([]byte(valToUnmarshall), &item)

	a1.Logger.Debug(" Unmarshalled json : %+v", (errunm))
	a1.Logger.Debug("Policy type Name :  %v", (item.Name))

	return &item
}

func (rh *Resthook) CreatePolicyType(policyTypeId models.PolicyTypeID, httprequest models.PolicyTypeSchema) error {
	a1.Logger.Debug("CreatePolicyType function")
	if policyTypeId != models.PolicyTypeID(*httprequest.PolicyTypeID) {
		//error message
		a1.Logger.Debug("Policytype Mismatch")
		return typeMismatchError
	}
	key := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
	a1.Logger.Debug("key %+v ", key)
	if data, err := httprequest.MarshalBinary(); err == nil {
		a1.Logger.Debug("Marshaled String : %+v", string(data))
		success, err1 := rh.db.SetIfNotExists(a1MediatorNs, key, string(data))
		a1.Logger.Info("success:%+v", success)
		if err1 != nil {
			a1.Logger.Error("error :%+v", err1)
			return err1
		}
		if !success {
			a1.Logger.Debug("Policy type %+v already exist", policyTypeId)
			return typeAlreadyError
		}
	}
	return nil
}

func toStringKeys(val interface{}) (interface{}, error) {
	var err error
	switch val := val.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			k, ok := k.(string)
			if !ok {
				return nil, errors.New("found non-string key")
			}
			m[k], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case []interface{}:
		var l = make([]interface{}, len(val))
		for i, v := range val {
			l[i], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return val, nil
	}
}

func validate(httpBodyString string, schemaString string) bool {
	var m interface{}
	err := yaml.Unmarshal([]byte(httpBodyString), &m)
	if err != nil {
		a1.Logger.Error("Unmarshal error : %+v", err)
	}
	m, err = toStringKeys(m)
	if err != nil {
		a1.Logger.Error("Conversion to string error : %+v", err)
		return false
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(schemaString)); err != nil {
		a1.Logger.Error("string reader error : %+v", err)
		return false
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		a1.Logger.Error("schema json compile error : %+v", err)
		return false
	}
	if err := schema.Validate(m); err != nil {
		a1.Logger.Error("schema validation error : %+v", err)
		return false
	}
	a1.Logger.Debug("validation successfull")
	return true
}

func (rh *Resthook) storePolicyInstance(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID, httpBody interface{}) (string, error) {
	var keys [1]string
	operation := "CREATE"
	typekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
	keys[0] = typekey

	a1.Logger.Debug("key1 : %+v", typekey)

	valmap, err := rh.db.Get(a1MediatorNs, keys[:])
	if err != nil {
		a1.Logger.Error("policy type error : %+v", err)
	}
	a1.Logger.Debug("policytype map : %+v", valmap)
	if valmap[typekey] == nil {
		a1.Logger.Error("policy type Not Present for policyid : %v", policyTypeId)
		return operation, policyTypeNotFoundError
	}
	// TODO : rmr creation_timestamp := time.Now() // will be needed for rmr to notify the creation of instance

	instancekey := a1InstancePrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
	keys[0] = typekey
	instanceMap, err := rh.db.Get(a1MediatorNs, keys[:])
	if err != nil {
		a1.Logger.Error("policy type error : %v", err)
	}
	a1.Logger.Debug("policyinstancetype map : %+v", instanceMap)

	if instanceMap[instancekey] != nil {
		operation = "UPDATE"
		a1.Logger.Debug("UPDATE")
		data, _ := json.Marshal(httpBody)
		a1.Logger.Debug("Marshaled String : %+v", string(data))
		a1.Logger.Debug("key   : %+v", instancekey)
		success, err1 := rh.db.SetIf(a1MediatorNs, instancekey, instanceMap[instancekey], string(data))
		if err1 != nil {
			a1.Logger.Error("error2 :%+v", err1)
			return operation, err1
		}
		if !success {
			a1.Logger.Debug("Policy instance %+v already exist", policyInstanceID)
			return operation, InstanceAlreadyError
		}
	} else {
		data, _ := json.Marshal(httpBody)
		a1.Logger.Debug("Marshaled String : %+v", string(data))
		a1.Logger.Debug("key   : %+v", instancekey)

		var instance_map []interface{}
		instance_map = append(instance_map, instancekey, string(data))
		a1.Logger.Debug("policyinstancetype map : %+v", instance_map[1])
		a1.Logger.Debug("policyinstancetype to create : %+v", instance_map)

		err1 := rh.db.Set(a1MediatorNs, instancekey, string(data))
		if err1 != nil {
			a1.Logger.Error("error1 :%+v", err1)
			return operation, err1
		}
	}
	a1.Logger.Debug("Policy Instance created ")
	return operation, nil
}

func (rh *Resthook) storePolicyInstanceMetadata(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID) (bool, error) {

	creation_timestamp := time.Now()
	instanceMetadataKey := a1InstanceMetadataPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)

	a1.Logger.Debug("key : %+v", instanceMetadataKey)

	var metadatajson []interface{}
	metadatajson = append(metadatajson, map[string]string{"created_at": creation_timestamp.Format("2006-01-02 15:04:05"), "has_been_deleted": "False"})
	metadata, _ := json.Marshal(metadatajson)

	a1.Logger.Debug("policyinstanceMetaData to create : %+v", string(metadata))

	err := rh.db.Set(a1MediatorNs, instanceMetadataKey, string(metadata))

	if err != nil {
		a1.Logger.Error("error :%+v", err)
		return false, err
	}

	a1.Logger.Debug("Policy Instance Meta Data created at :%+v", creation_timestamp)

	return true, nil
}

func (rh *Resthook) CreatePolicyInstance(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID, httpBody interface{}) error {
	a1.Logger.Debug("CreatePolicyInstance function")
	//  validate the PUT against the schema
	var policyTypeSchema *models.PolicyTypeSchema
	policyTypeSchema = rh.GetPolicyType(policyTypeId)
	schemaStr, err := json.Marshal(policyTypeSchema.CreateSchema)
	if err != nil {
		a1.Logger.Error("Json Marshal error : %+v", err)
		return err
	}
	a1.Logger.Debug("schema to validate %+v", string(schemaStr))
	a1.Logger.Debug("httpbody to validate %+v", httpBody)
	schemaString := fmt.Sprint(string(schemaStr))
	httpBodyMarshal, err := json.Marshal(httpBody)
	httpBodyString := string((httpBodyMarshal))
	a1.Logger.Debug("schema to validate sprint  %+v", (schemaString))
	a1.Logger.Debug("httpbody to validate sprint %+v", httpBodyString)
	isvalid := validate(httpBodyString, schemaString)
	if isvalid {
		var operation string
		operation, err = rh.storePolicyInstance(policyTypeId, policyInstanceID, httpBody)
		if err != nil {
			a1.Logger.Error("error :%+v", err)
			return err
		}
		a1.Logger.Debug("policy instance :%+v", operation)
		iscreated, errmetadata := rh.storePolicyInstanceMetadata(policyTypeId, policyInstanceID)
		if errmetadata != nil {
			a1.Logger.Error("error :%+v", errmetadata)
			return errmetadata
		}
		if iscreated {
			a1.Logger.Debug("policy instance metadata created")
		}
		isSent := rh.iRmrSenderInst.RmrSendToXapp(httpBodyString)
		if isSent {
			a1.Logger.Debug("rmrSendToXapp : message sent")
		} else {
			a1.Logger.Debug("rmrSendToXapp : message not sent")
		}

	} else {
		a1.Logger.Error("%+v", invalidJsonSchema)
		return invalidJsonSchema
	}

	return nil
}

func (rh *Resthook) GetPolicyInstance(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID) (interface{}, error) {
	a1.Logger.Debug("GetPolicyInstance1")

	var keys [1]string

	typekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
	keys[0] = typekey

	a1.Logger.Debug("key1 : %+v", typekey)

	valmap, err := rh.db.Get(a1MediatorNs, keys[:])
	if len(valmap) == 0 {
		a1.Logger.Debug("policy type Not Present for policyid : %v", policyTypeId)
		return "{}", policyTypeNotFoundError
	}

	if err != nil {
		a1.Logger.Error("error in retrieving policy type. err: %v", err)
		return "{}", err
	}

	if valmap[typekey] == nil {
		a1.Logger.Debug("policy type Not Present for policyid : %v", policyTypeId)
		return "{}", policyTypeNotFoundError
	}

	a1.Logger.Debug("keysmap : %+v", valmap[typekey])

	instancekey := a1InstancePrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
	a1.Logger.Debug("key2 : %+v", instancekey)
	keys[0] = instancekey
	instanceMap, err := rh.db.Get(a1MediatorNs, keys[:])
	if err != nil {
		a1.Logger.Error("policy instance error : %v", err)
	}
	a1.Logger.Debug("policyinstancetype map : %+v", instanceMap)

	if instanceMap[instancekey] == nil {
		a1.Logger.Debug("policy instance Not Present for policyinstaneid : %v", policyInstanceID)
		return "{}", policyInstanceNotFoundError
	}

	valStr := fmt.Sprint(instanceMap[instancekey])
	return valStr, nil
}

func (rh *Resthook) GetAllPolicyInstance(policyTypeId models.PolicyTypeID) ([]models.PolicyInstanceID, error) {
	a1.Logger.Debug("GetAllPolicyInstance")
	var policyTypeInstances = []models.PolicyInstanceID{}

	keys, err := rh.db.GetAll("A1m_ns")

	if err != nil {
		a1.Logger.Error("error in retrieving policy. err: %v", err)
		return policyTypeInstances, err
	}
	a1.Logger.Debug("keys : %+v", keys)
	typekey := a1InstancePrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "."

	for _, key := range keys {
		if strings.HasPrefix(strings.TrimLeft(key, " "), typekey) {
			pti := strings.Split(strings.Trim(key, " "), typekey)[1]
			a1.Logger.Debug("pti %+v", pti)
			policyTypeInstances = append(policyTypeInstances, models.PolicyInstanceID(pti))
		}
	}

	if len(policyTypeInstances) == 0 {
		a1.Logger.Debug("policy instance Not Present  ")
	}

	a1.Logger.Debug("return : %+v", policyTypeInstances)
	return policyTypeInstances, nil
}

func (rh *Resthook) DeletePolicyType(policyTypeId models.PolicyTypeID) error {
	a1.Logger.Debug("DeletePolicyType")

	policyinstances, err := rh.GetAllPolicyInstance(policyTypeId)
	if err != nil {
		a1.Logger.Error("error in retrieving policy. err: %v", err)
		return err
	}

	var keys [1]string
	key := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
	keys[0] = key
	if len(policyinstances) == 0 {
		err := rh.db.Remove(a1MediatorNs, keys[:])
		if err != nil {
			a1.Logger.Error("error in deleting policy type err: %v", err)
			return err
		}
	} else {
		a1.Logger.Error("tried to delete a type that isn't empty")
		return policyTypeCanNotBeDeletedError
	}
	return nil
}

func (rh *Resthook) typeValidity(policyTypeId models.PolicyTypeID) error {
	var keys [1]string

	typekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10)
	keys[0] = typekey

	a1.Logger.Debug("key1 : %+v", typekey)
	valmap, err := rh.db.Get(a1MediatorNs, keys[:])
	if err != nil {
		a1.Logger.Error("error in retrieving policytype err: %v", err)
		return err
	}
	if len(valmap) == 0 {
		a1.Logger.Error("policy type Not Present for policyid : %v", policyTypeId)
		return policyTypeNotFoundError
	}
}

func (rh *Resthook) instanceValidity(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID) error {
	err := rh.typeValidity(policyTypeId)
	if err != nil {
		return err
	}
	policyTypeInstances, err := rh.GetPolicyInstance(policyTypeId, policyInstanceID)
	if err != nil {
		a1.Logger.Error("policy instance error : %v", err)
		return err
	}
	if len(policyTypeInstances.(string)) == 0 {
		a1.Logger.Debug("policy instance Not Present  ")
		return policyInstanceNotFoundError
	}
}

func (rh *Resthook) getMetaData(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID) (map[string]interface{}, error) {
	instanceMetadataKey := a1InstanceMetadataPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
	a1.Logger.Debug("instanceMetadata key : %+v", instanceMetadataKey)
	var keys [1]string
	keys[0] = instanceMetadataKey
	instanceMetadataMap, err := rh.db.Get(a1MediatorNs, keys[:])
	if err != nil {
		a1.Logger.Error("policy instance error : %v", err)
	}
	a1.Logger.Debug("instanceMetadata map : %+v", instanceMetadataMap)
	if instanceMetadataMap[instanceMetadataKey] == nil {
		a1.Logger.Error("policy instance Not Present for policyinstaneid : %v", policyInstanceID)
		return map[string]interface{}{}, policyInstanceNotFoundError
	}
	return instanceMetadataMap, nil
}

func (rh *Resthook) GetPolicyInstanceStatus(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID) (*a1_mediator.A1ControllerGetPolicyInstanceStatusOKBody, error) {
	err := rh.instanceValidity(policyTypeId, policyInstanceID)
	if err != nil && err == policyInstanceNotFoundError || err == policyTypeNotFoundError {
		policyInstanceStatus.InstanceStatus = "NOT IN EFFECT"
		return &policyInstanceStatus, err
	}
	policyInstanceStatus := a1_mediator.A1ControllerGetPolicyInstanceStatusOKBody{}
	metadata, err := rh.getMetaData(policyTypeId, policyInstanceID)
	a1.Logger.Debug(" metadata %v", metadata)
	if err != nil {
		a1.Logger.Error("policy instance error : %v", err)
		policyInstanceStatus.InstanceStatus = "NOT IN EFFECT"
		return &policyInstanceStatus, err
	}
	jsonbody, err := json.Marshal(metadata)
	if err != nil {
		a1.Logger.Error("marshal error : %v", err)
	}

	if err := json.Unmarshal(jsonbody, &policyInstanceStatus); err != nil {
		a1.Logger.Error("unmarshal error : %v", err)
	}
	if policyInstanceStatus.HasBeenDeleted == false {
		policyInstanceStatus.InstanceStatus = "IN EFFECT"
	} else {
		policyInstanceStatus.InstanceStatus = "NOT IN EFFECT"
	}
	return &policyInstanceStatus, nil
}
