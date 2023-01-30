package rmr
import (
     "github.com/spf13/viper"
)

type Config struct {
    Name              string `mapstructure:"NAME"`
    MaxSize           string `mapstructure:"MAX_SIZE"`
    ThreadType        string `mapstructure:"THREAD_TYPE"`
    LowLatency        string `mapstructure:"LOW_LATENCY"`
    FastAck           string `mapstructure:"FAST_ACK"`
    MaxRetryOnFailure string `mapstructure:"MAX_RETRY_ON_FAILURE"`
    Port              string `mapstructure:"PORT"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")

    viper.AutomaticEnv()

    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    err = viper.Unmarshal(&config)
    return
}
