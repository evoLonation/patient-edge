package config

type Conf struct {
	DB DBConf
}
type DBConf struct {
	DataSource string
}

type RpcClientConf struct {
}
