package conf

type LogTransfer struct {
	KafkaConf `ini:"kafka"`
	ESConf    `ini:"es"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type ESConf struct {
	Address string `ini:"address"`
}
