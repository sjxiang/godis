package main

import "sync"


type KV struct {
	data map[string][]byte
	mu   sync.Mutex
}

func NewKV() *KV {
	return &KV{
		data: map[string][]byte{},
	}
}

func (kv *KV) Set(key, value []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[string(key)] = []byte(value)
	return nil
}

func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	value, ok :=  kv.data[string(key)]
	return value, ok
}


