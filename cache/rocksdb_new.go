package cache

// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import "runtime"

//初始化RocksDB
func newRocksdbCache() *rocksdbCache {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)
	var e *C.char
	//设置存储路径
	db := C.rocksdb_open(options, C.CString("/mnt/rocksdb"), &e)
	if e != nil {
		panic(C.GoString(e))
	}
	C.rocksdb_options_destroy(options)
	//ch缓冲区数量为5000，每次写入的pair少于5000就不会阻塞
	c := make(chan *pair, 5000)
	wo := C.rocksdb_writeoptions_create()
	//开辟一个协程执行写入
	go write_func(db, c, wo)
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), wo, e, c}
}
