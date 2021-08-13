package storage

import (
	"balance/network"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// setting
const (
	server    = "192.168.2.122"
	port      = 6379
	password  = "farmer233"
	db        = 5
	REDIS_KEY = "ipList"
	SCORE_MAX = 50
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

func add(addrs network.Addrs) (err error) {
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

	return nil
}

// 降低目标地址分数
func decrease(addr network.Addr) (err error) {
	address := fmt.Sprintf("%s://%s:%d", addr.Protocol, addr.Host, addr.Port)
	err = rdb.ZIncrBy(ctx, REDIS_KEY, -1, address).Err()
	if err != nil {
		return err
	}
	score := rdb.ZScore(ctx, REDIS_KEY, address)
	if score.Val() <= 0.00 {
		log.Printf("%v current score %v, remove.\n", address, score.Val())
		err := rdb.ZRem(ctx, REDIS_KEY, address).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// set jwglxt to max score
func max(addr network.Addr) (err error) {
	address := fmt.Sprintf("%s://%s:%d", addr.Protocol, addr.Host, addr.Port)
	err = rdb.ZAdd(ctx, REDIS_KEY, &redis.Z{
		Score:  SCORE_MAX,
		Member: address,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

// /////// api service func ///////

// get count of jwglxt
func count() int64 {
	count, err := rdb.ZCard(ctx, REDIS_KEY).Result()
	if err != nil {
		return -1
	}
	return count
}

// 检测地址是否存在
func exists(addr network.Addr) bool {
	address := fmt.Sprintf("%s://%s:%d", addr.Protocol, addr.Host, addr.Port)
	err := rdb.ZScore(ctx, REDIS_KEY, address).Err()
	if err == nil {
		return true
	}
	return false
}

// get batch of jwglxt
func batch(cursor uint64, match string, count int64) (res []string, retCursor uint64, err error) {
	res, retCursor, err = rdb.ZScan(ctx, REDIS_KEY, cursor, match, count).Result()
	if err != nil {
		return []string{}, 0, err
	}
	return res, retCursor, err
}

func GetBatch(cursor uint64, match string, count int64) (map[string]string, uint64, error) {
	ret := make(map[string]string)
	res, retCursor, err := batch(cursor, match, count)
	for i := 0; i < len(res); i += 2 {
		ret[res[i]] = res[i+1]
	}
	return ret, retCursor, err
}
