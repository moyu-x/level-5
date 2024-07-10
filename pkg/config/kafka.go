package config

type KafkaConfig struct {
	FilePath string
	Round    int
	Mode     string
	Data     string
	Topic    string
	GroupId  string
}
