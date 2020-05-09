package conf

//Conf ...
type Conf struct {
	KafkaConf   `ini:"kafka"`
	TaillogConf `ini:"taillog"`
}

type KafkaConf struct {
	Address []string `ini:"hosts"`
	Topic   string   `ini:"topic"`
}
type TaillogConf struct {
	Path string `ini:"path"`
}
