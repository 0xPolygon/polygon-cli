package dbbench

import (
	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/bloom"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"runtime"
	"sync"
)

type (
	PebbleDBWrapper struct {
		handle *pebble.DB
		wo     *pebble.WriteOptions
		sync.Mutex
	}
	WrappedPebbleIterator struct {
		*pebble.Iterator
		*sync.Mutex
	}
)

func NewWrappedPebbleDB() (*PebbleDBWrapper, error) {
	memTableLimit := 2
	memTableSize := *cacheSize * 1024 * 1024 / 2 / memTableLimit
	opt := &pebble.Options{
		Cache:                       pebble.NewCache(int64(*cacheSize * 1024 * 1024)),
		MemTableSize:                uint64(memTableSize),
		MemTableStopWritesThreshold: memTableLimit,
		MaxConcurrentCompactions:    func() int { return runtime.NumCPU() },
		Levels: []pebble.LevelOptions{
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
		},
		ReadOnly: *readOnly || *fullScan,
	}
	p, err := pebble.Open(*dbPath, opt)
	if err != nil {
		return nil, err
	}
	db := new(PebbleDBWrapper)
	db.handle = p
	db.wo = &pebble.WriteOptions{Sync: *syncWrites}
	return db, err
}

func (p *PebbleDBWrapper) Close() error {
	return p.handle.Close()
}
func (p *PebbleDBWrapper) Compact() error {
	// this is a hack to ideally get a key that's larger than most of the other keys
	return p.handle.Compact([]byte{0}, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, true)
}
func (p *PebbleDBWrapper) NewIterator() iterator.Iterator {
	io := pebble.IterOptions{
		LowerBound:                nil,
		UpperBound:                nil,
		TableFilter:               nil,
		PointKeyFilters:           nil,
		RangeKeyFilters:           nil,
		KeyTypes:                  0,
		RangeKeyMasking:           pebble.RangeKeyMasking{},
		OnlyReadGuaranteedDurable: false,
		UseL6Filters:              false,
	}
	iter, _ := p.handle.NewIter(&io)
	wrappedIter := WrappedPebbleIterator{iter, &p.Mutex}
	return &wrappedIter
}
func (w *WrappedPebbleIterator) Seek(key []byte) bool {
	// SeekGE has a different name but has the same logic as the IteratorSeeker `Seek` method
	return w.SeekGE(key)
}
func (w *WrappedPebbleIterator) SetReleaser(releaser util.Releaser) {
}
func (w *WrappedPebbleIterator) Release() {
	w.Iterator.Close()
}
func (p *PebbleDBWrapper) Get(key []byte) ([]byte, error) {
	resp, closer, err := p.handle.Get(key)
	if err != nil {
		return nil, err
	}
	closer.Close()
	return resp, nil
}
func (p *PebbleDBWrapper) Put(key []byte, value []byte) error {
	return p.handle.Set(key, value, p.wo)
}
