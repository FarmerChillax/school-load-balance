package storage

import (
	"balance/discover"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// setting
const (
	server        = "192.168.2.122"
	port          = 6379
	password      = "farmer233"
	db            = 5
	REDIS_KEY     = "proxies"
	SCORE_MAX     = 50
	SCORE_DEFAULT = 10
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

func add(addrs discover.Addrs) (err error) {
	for _, addr := range addrs {
		buf, err := addr.MarshalBinary()
		if err != nil {
			return nil
		}
		err = rdb.ZAdd(ctx, REDIS_KEY, &redis.Z{
			Score:  SCORE_DEFAULT,
			Member: buf,
		}).Err()

		if err != nil {
			return err
		}
	}

	return nil
}

// 降低目标地址分数
func decrease(addr discover.Addr) (err error) {
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
func max(addr discover.Addr) (err error) {
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
func exists(addr discover.Addr) bool {
	buf, err := addr.MarshalBinary()
	if err != nil {
		return false
	}
	err = rdb.ZScore(ctx, REDIS_KEY, string(buf)).Err()
	return err == nil
}

// get batch of jwglxtS
func batch(cursor uint64, match string, count int64) (res []string, retCursor uint64, err error) {
	res, retCursor, err = rdb.ZScan(ctx, REDIS_KEY, cursor, match, count).Result()
	if err != nil {
		return []string{}, 0, err
	}
	return res, retCursor, err
}

// format batch result value
func GetBatch(cursor uint64, match string, count int64) (res discover.Addrs, err error) {
	var addr discover.Addr
	var ret []string
	for {
		ret, cursor, err = batch(cursor, match, count)
		if err != nil {
			return res, err
		}
		for i := 0; i < len(ret); i += 2 {
			addr.UnmarshalBinary([]byte(ret[i]))
			res = append(res, addr)
		}
		if cursor <= 0 {
			break
		}
	}
	return res, nil
	// for i := 0; i < len(res); i += 2 {
	// 	err := addr.UnmarshalBinary([]byte(res[i]))
	// 	if err != nil {
	// 		return discover.Addrs, retCursor,
	// 	}
	// 	fmt.Println()
	// }
	// return ret, retCursor, err
}
