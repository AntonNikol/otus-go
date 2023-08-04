package main

import (
	"fmt"
	hw04lrucache "github.com/AntonNikol/hw04_lru_cache"
)

func main() {
	cache := hw04lrucache.NewCache(3)
	var a hw04lrucache.Key = "1"
	cache.Set(a, 10)

	fmt.Printf("cache %+v", cache)
}
