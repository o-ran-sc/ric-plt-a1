/*
==================================================================================
  Copyright (c) 2023 Samsung

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
package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigurationSuccess(t *testing.T) {
	config := ParseConfiguration()
	assert.Equal(t, "debug", config.Logging.LogLevel)
	assert.Equal(t, "", config.Name)
	assert.Equal(t, 65536, config.Rmr.MaxMsgSize)

	assert.Equal(t, 0, config.ThreadType)
	assert.Equal(t, false, config.LowLatency)
	assert.Equal(t, false, config.FastAck)
	assert.Equal(t, 1, config.MaxRetryOnFailure)
	assert.Equal(t, 4561, config.Port)
}
