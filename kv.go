package main

import "sync"


type KV struct {
	data map[string][]byte
	mu sync.Mutex
}

func NewKV() *KV {
	return &KV{
		data: map[string][]byte{},
	}
}

func (kv *KV) Set(key, value string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = []byte(value)

	return nil
}

func (kv *KV) Get(key string) ([]byte, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	v, ok := kv.data[key]

	return v, ok
}


