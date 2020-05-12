package conf

//Conf ...
type Conf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}

//KafkaConf ...
type KafkaConf struct {
	Address []string `ini:"hosts"`
}

//EtcdConf ...
type EtcdConf struct {
	Address []string `ini:"address"`
	Timeout int      `ini:"timeout"`
}

//TaillogConf ...
// type TaillogConf struct {
// 	Path string `ini:"path"`
// }
