package cache

//节点再平衡实现功能
type Scanner interface {
	Scan() bool
	Key() string
	Value() []byte
	Close()
}
