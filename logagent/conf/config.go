package conf

type Appconf struct {
	KafkaConf   `ini:"kafka"`
	EtcdConf    `ini:"etcd"`
	TaillogConf `ini:"taillog"`
}

type KafkaConf struct {
	Address string `ini:"address"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Key     string `ini:"key"`
	TaskMax int    `ini:"taskMax"`
}

/*  ----- unused↓ ------ */
type TaillogConf struct {
	FileName string `ini:"filename"`
	ChanMax  int    `ini:"chanMax"`
}
