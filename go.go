package main

import (
	"github.com/go-redis/redis"
	"fmt"
)

func main() {
	rDb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	k := "weibo.com-host.list"
	lenInt64 := rDb.LLen(k)
	lenE := lenInt64.Val()
	listRs := rDb.LRange(k, 0, lenE)
	listRsE := listRs.Val()
	for i := 0; i < int(lenE); i++ {
		fmt.Println(listRsE[i])
	}
}
