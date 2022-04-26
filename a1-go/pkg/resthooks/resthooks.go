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

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v2"
)

const (
	a1PolicyPrefix = "a1.policy_type."
	a1MediatorNs   = "A1m_ns"
)

var typeAlreadyError = errors.New("Policy Type already exists")
var InstanceAlreadyError = errors.New("Policy Instance already exists")
var typeMismatchError = errors.New("Policytype Mismatch")
var invalidJsonSchema = errors.New("Invalid Json ")
var policyInstanceNotFoundError = errors.New("Policy Instance Not Found")
var policyTypeNotFoundError = errors.New("Policy Type Not Found")

func (rh *Resthook) IsPolicyTypePresent(err error) bool {
	return err == policyTypeNotFoundError
}

func (rh *Resthook) IsPolicyInstancePresent(err error) bool {
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
	return createResthook(sdlgo.NewSyncStorage())
}

func createResthook(sdlInst iSdl) *Resthook {
	return &Resthook{
		db: sdlInst,
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

	a1.Logger.Debug("policytype map : ", valmap)

	if len(valmap) == 0 {
		a1.Logger.Error("policy type Not Present for policyid : %v", policyTypeId)
		return policytypeschema
	}

	if err != nil {
		a1.Logger.Error("error in retrieving policy type. err: %v", err)
		return policytypeschema
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

func validate(yamlText interface{}, schemaText string) bool {
	var m interface{}
	err := yaml.Unmarshal([]byte(yamlText.(string)), &m)
	if err != nil {
		a1.Logger.Error("Unmarshal error : %+v", err)
	}
	m, err = toStringKeys(m)
	if err != nil {
		a1.Logger.Error("Conversion to string error : %+v", err)
		return false
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(schemaText)); err != nil {
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

func (rh *Resthook) StorePolicyInstance(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID, httpBody interface{}) (string, error) {
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
		//return negative
	}
	// creation_timestamp := time.Now() // will be needed for rmr to notify the creation of instance

	instancekey := a1PolicyPrefix + strconv.FormatInt((int64(policyTypeId)), 10) + "." + string(policyInstanceID)
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

func (rh *Resthook) CreatePolicyInstance(policyTypeId models.PolicyTypeID, policyInstanceID models.PolicyInstanceID, httpBody interface{}) error {
	a1.Logger.Debug("CreatePolicyInstance function")
	//  validate the PUT against the schema
	var policyTypeSchema *models.PolicyTypeSchema
	policyTypeSchema = rh.GetPolicyType(policyTypeId)
	out, err := json.Marshal(policyTypeSchema.CreateSchema)
	if err != nil {
		a1.Logger.Error("Json Marshal error : %+v", err)
	}
	a1.Logger.Debug("schema to validate %+v", string(out))
	a1.Logger.Debug("httpbody to validate %+v", httpBody)

	isvalid := validate(httpBody, string(out))
	if isvalid {
		operation, err := rh.StorePolicyInstance(policyTypeId, policyInstanceID, httpBody)
		if err != nil {
			a1.Logger.Error("error :%+v", err)
			return err
		}
		a1.Logger.Debug("policy instance :%+v", operation)
	} else {
		a1.Logger.Error("%+v", invalidJsonSchema)
		return invalidJsonSchema
	}

	return nil
}
