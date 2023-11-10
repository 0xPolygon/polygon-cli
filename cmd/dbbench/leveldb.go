package dbbench

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type (
	LevelDBWrapper struct {
		ro     *opt.ReadOptions
		wo     *opt.WriteOptions
		handle *leveldb.DB
	}
)

func NewWrappedLevelDB() (*LevelDBWrapper, error) {
	db, err := leveldb.OpenFile(*dbPath, &opt.Options{
		Filter:                 filter.NewBloomFilter(10),
		DisableSeeksCompaction: true,
		OpenFilesCacheCapacity: *openFilesCacheCapacity,
		BlockCacheCapacity:     *cacheSize / 2 * opt.MiB,
		WriteBuffer:            *cacheSize / 4 * opt.MiB,
		// if we've disabled writes, or we're doing a full scan, we should open the database in read only mode
		ReadOnly: *readOnly || *fullScan,
	})
	if err != nil {
		return nil, err
	}

	wo := &opt.WriteOptions{
		NoWriteMerge: *noWriteMerge,
		Sync:         *syncWrites,
	}
	ro := &opt.ReadOptions{
		DontFillCache: *dontFillCache,
	}
	if *readStrict {
		ro.Strict = opt.StrictAll
	} else {
		ro.Strict = opt.DefaultStrict
	}
	if *nilReadOptions {
		ro = nil
	}
	wrapper := new(LevelDBWrapper)
	wrapper.handle = db
	wrapper.wo = wo
	wrapper.ro = ro
	return wrapper, nil
}
func (l *LevelDBWrapper) Close() error {
	return l.handle.Close()
}
func (l *LevelDBWrapper) Compact() error {
	return l.handle.CompactRange(util.Range{Start: nil, Limit: nil})
}
func (l *LevelDBWrapper) NewIterator() iterator.Iterator {
	return l.handle.NewIterator(nil, nil)
}
func (l *LevelDBWrapper) Get(key []byte) ([]byte, error) {
	return l.handle.Get(key, l.ro)
}
func (l *LevelDBWrapper) Put(key []byte, value []byte) error {
	return l.handle.Put(key, value, l.wo)
}
