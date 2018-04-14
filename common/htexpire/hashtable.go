package htexpire

import (
	"hash/fnv"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// HashTable ..
type HashTable struct {
	locker *sync.Mutex
	data   map[int]time.Time
}

// New ..
func New() *HashTable {
	ht := &HashTable{
		data:   make(map[int]time.Time),
		locker: new(sync.Mutex),
	}
	go ht.runLoop()
	return ht
}

// Size ..
func (ht *HashTable) Size() int {
	return len(ht.data)
}

// Set ..
func (ht *HashTable) Set(key string) error {
	ht.locker.Lock()
	defer ht.locker.Unlock()

	if ht.Size() == maxInt-1 {
		ht.check()
	}
	var position int
	var err error
	if position, err = ht.genHash(key); err != nil {
		return err
	}
	ht.data[position] = time.Now()
	return nil
}

func (ht *HashTable) check() {
	for {
		ht.runLoop()
		if ht.Size() < maxInt-1 {
			return
		}
	}
}

// Get ..
func (ht *HashTable) Get(key string) (bool, error) {
	var position int
	var err error
	if position, err = ht.genHash(key); err != nil {
		return false, err
	}
	var val time.Time
	var ok bool
	if val, ok = ht.data[position]; ok {
		if int(time.Since(val).Seconds()) > viper.GetInt("Crawler.UniReqTimeout") {
			if err = ht.Remove(key); err != nil {
				return true, err
			}
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

// Remove ..
func (ht *HashTable) Remove(key string) error {
	ht.locker.Lock()
	defer ht.locker.Unlock()
	var position int
	var err error
	if position, err = ht.genHash(key); err != nil {
		return err
	}
	var b bool
	if b, err = ht.Get(key); err != nil {
		return err
	} else if b {
		ht.remove(position)
	}
	return nil
}

func (ht *HashTable) remove(key int) {
	delete(ht.data, key)
}

const (
	minUint uint = 0
	maxUint      = ^minUint
	maxInt       = int(maxUint >> 1)
)

func (ht *HashTable) genHash(s string) (int, error) {
	hash := fnv.New64()
	_, err := hash.Write([]byte(s))

	return int(hash.Sum64() % uint64(maxInt)), err
}

func (ht *HashTable) runLoop() {
	for key, val := range ht.data {
		if int(time.Since(val).Seconds()) > viper.GetInt("Crawler.UniReqTimeout") {
			ht.remove(key)
		}
	}
	time.Sleep(time.Second * 10)
}
