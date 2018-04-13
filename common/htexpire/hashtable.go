package htexpire

import (
	"hash/fnv"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type HashTable struct {
	locker *sync.Mutex
	data   map[int]time.Time
}

func New() *HashTable {
	ht := &HashTable{
		data:   make(map[int]time.Time),
		locker: new(sync.Mutex),
	}
	go ht.runLoop()
	return ht
}

func (ht *HashTable) Size() int {
	return len(ht.data)
}

func (ht *HashTable) Set(key string) {
	ht.locker.Lock()
	defer ht.locker.Unlock()

	if ht.Size() == maxInt-1 {
		ht.check()
	}
	ht.data[ht.genHash(key)] = time.Now()
}

func (ht *HashTable) check() {
	for {
		ht.runLoop()
		if ht.Size() < maxInt-1 {
			return
		}
	}
}

func (ht *HashTable) Get(key string) bool {
	if _, ok := ht.data[ht.genHash(key)]; ok {
		return true
	}
	return false
}

func (ht *HashTable) Remove(key string) {
	ht.locker.Lock()
	defer ht.locker.Unlock()

	if ht.Get(key) {
		ht.remove(ht.genHash(key))
	}
}

func (ht *HashTable) remove(key int) {
	delete(ht.data, key)
}

const (
	minUint uint = 0
	maxUint      = ^minUint
	maxInt       = int(maxUint >> 1)
	minInt       = ^maxInt
)

func (ht *HashTable) genHash(s string) int {
	hash := fnv.New64()
	hash.Write([]byte(s))

	return int(hash.Sum64() % uint64(maxInt))
}

func (ht *HashTable) runLoop() {
	for key, val := range ht.data {
		if int(time.Since(val).Seconds()) > viper.GetInt("Crawler.UniReqTimeout") {
			ht.remove(key)
		}
	}
	time.Sleep(time.Second * 10)
}
