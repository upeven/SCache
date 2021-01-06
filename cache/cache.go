package cache

//声明缓存服务支持的操作
type Cache interface {
	//设置缓存
	Set(string, []byte) error
	//获取缓存值
	Get(string) ([]byte, error)
	//删除缓存
	Del(string) error
	//获取缓存服务器状态
	GetStat() Stat
	NewScanner() Scanner
}
