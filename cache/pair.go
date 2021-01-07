package cache

//存放从客户端发来的键值对，存入channel中
type pair struct {
	k string
	v []byte
}
