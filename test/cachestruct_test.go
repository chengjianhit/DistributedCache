package test

import (
	"CacheAdmin/cache/lrustruct"
	"fmt"
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	cache := lrustruct.NewCache(int64(0), nil)
	cache.Set("name", String("chengjian"))
	value, succ := cache.Get("name")

	fmt.Printf("value is %s  result %t", value, succ)

}

func TestAdd(t *testing.T) {
	cache := lrustruct.NewCache(int64(0), nil)
	cache.Set("key1", String("cheng"))
	cache.Set("key2", String("jian"))
	size := cache.Len()
	println("size is %s", size)
}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := lrustruct.NewCache(int64(cap), nil)
	lru.Set(k1, String(v1))
	lru.Set(k2, String(v2))
	lru.Set(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value lrustruct.Value) {
		fmt.Printf("begin Onvicted, the key is %s, the value is %s \n", key, value)
		keys = append(keys, key)
	}

	cache := lrustruct.NewCache(int64(10), callback)
	cache.Set("key1", String("123456"))
	cache.Set("k2", String("k2"))
	cache.Set("k3", String("k3"))
	cache.Set("k4", String("k4"))

	expect := []string{"key1", "k2"}

	equal := reflect.DeepEqual(expect, keys)
	if equal {
		println("the keys equal expect")
	} else {
		println("the keys not equal expect")

	}

}
