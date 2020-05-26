package conf

//Conf ...
type Conf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}

//KafkaConf ...
type KafkaConf struct {
	Address []string `ini:"hosts"`
	LogSize int      `ini:"logchansize"`
}

//EtcdConf ...
type EtcdConf struct {
	Address []string `ini:"address"`
	Timeout int      `ini:"timeout"`
	Logkey  string   `ini:"logkey"`
}

//TaillogConf ...
// type TaillogConf struct {
// 	Path string `ini:"path"`
// }
