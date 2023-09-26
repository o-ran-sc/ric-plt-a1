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
	"os"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"github.com/spf13/viper"
)

type Configuration struct {
	LogLevel          string
	Name              string
	MaxSize           int
	ThreadType        int
	LowLatency        bool
	FastAck           bool
	MaxRetryOnFailure int
	Port              int
}

func ParseConfiguration() *Configuration {
	viper.SetConfigType("yaml")
	viper.SetConfigName("configuration")
	configFile := os.Getenv("A1_CONFIG_FILE")
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		a1.Logger.Info("#configuration.ParseConfiguration - failed to read configuration file: %s\n", err)
	}

	config := Configuration{}
	config.LogLevel = viper.GetString("log-level")
	viper.SetDefault("log-level", "info")
	config.Name = viper.GetString("NAME")
	viper.SetDefault("NAME", "")
	config.MaxSize = viper.GetInt("MAX_SIZE")
	viper.SetDefault("MAX_SIZE", 65534)
	config.ThreadType = viper.GetInt("THREAD_TYPE")
	viper.SetDefault("THREAD_TYPE", 0)
	config.LowLatency = viper.GetBool("LOW_LATENCY")
	viper.SetDefault("LOW_LATENCY", false)
	config.FastAck = viper.GetBool("FAST_ACK")
	viper.SetDefault("FAST_ACK", false)
	config.MaxRetryOnFailure = viper.GetInt("MAX_RETRY_ON_FAILURE")
	viper.SetDefault("MAX_RETRY_ON_FAILURE", 1)
	config.Port = viper.GetInt("PORT")
	viper.SetDefault("PORT", 4562)
	return &config
}
