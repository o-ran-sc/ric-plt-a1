package rmr

type Config struct {
    Name              string `mapstructure:"NAME"`
    MaxSize           int `mapstructure:"MAX_SIZE"`
    ThreadType        int `mapstructure:"THREAD_TYPE"`
    LowLatency        bool `mapstructure:"LOW_LATENCY"`
    FastAck           bool `mapstructure:"FAST_ACK"`
    MaxRetryOnFailure int `mapstructure:"MAX_RETRY_ON_FAILURE"`
    Port              int `mapstructure:"PORT"`
}

