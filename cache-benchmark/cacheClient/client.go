package cacheClient

//命令结构体
/*
Name:命令名称
Key:键
Value:值
Error:错误
 */
type Cmd struct {
	Name  string
	Key   string
	Value string
	Error error
}
//客户端接口
type Client interface {
	Run(*Cmd)
	PipelinedRun([]*Cmd)
}


//根据相应的类型创建指定的客户端
func New(typ, server string) Client {
	if typ == "redis" {
		return newRedisClient(server)
	}
	if typ == "http" {
		return newHTTPClient(server)
	}
	if typ == "tcp" {
		return newTCPClient(server)
	}
	panic("unknown client type " + typ)
}
