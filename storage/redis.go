package storage

import (
	"balance/network"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// setting
const (
	server    = "192.168.2.122"
	port      = 6379
	password  = "farmer233"
	db        = 5
	REDIS_KEY = "ipList"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%d", server, port),
	Password: password,
	DB:       db,
})

var ctx = context.Background()

func init() {
	fmt.Println(rdb.Ping(ctx))
}

func pingDB() string {
	res := fmt.Sprintf("%v", rdb.Ping(ctx))
	return res
}

func writeDB(addrs network.Addrs) (err error) {
	for _, addr := range addrs {
		address := fmt.Sprintf("%s://%s:%d", addr.Protocol, addr.Host, addr.Port)
		err = rdb.ZAdd(ctx, REDIS_KEY, &redis.Z{
			Score:  10,
			Member: address,
		}).Err()
		if err != nil {
			return err
		}
	}
	// err = rdb.Set(ctx, key, ports, 0).Err()
	// for _, port := range ports {
	// 	address := fmt.Sprintf("%s:%d", host, port)
	// 	err = rdb.ZAdd(ctx, REDIS_KEY, &redis.Z{
	// 		Score:  10,
	// 		Member: address,
	// 	}).Err()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
