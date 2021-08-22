package storage

import (
	"balance/discover"
	"balance/utils"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

// setting
const (
	SCORE_MAX     = 50
	SCORE_MIN     = 0
	SCORE_DEFAULT = 10
)

// var rdb = redis.NewClient(&redis.Options{
// 	Addr:     fmt.Sprintf("%s:%d", server, port),
// 	Password: password,
// 	DB:       db,
// })

var rdb *redis.Client
var redisConfig utils.Redis

var ctx = context.Background()

func init() {
	redisConfig = utils.Config.Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Server, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Panicf("redis server connect timeout, err: %v", err)
		os.Exit(0)
	}
	fmt.Println("redis server connect success")
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
		err = rdb.ZAdd(ctx, redisConfig.Key, &redis.Z{
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
	member, err := addr.MarshalBinary()
	if err != nil {
		return err
	}
	err = rdb.ZIncrBy(ctx, redisConfig.Key, -1, string(member)).Err()
	if err != nil {
		return err
	}
	score := rdb.ZScore(ctx, redisConfig.Key, string(member))
	if score.Val() <= 0.00 {
		log.Printf("%v current score %v, remove.\n", addr, score.Val())
		err := rdb.ZRem(ctx, redisConfig.Key, member).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// set jwglxt to max score
func max(addr discover.Addr) (err error) {
	member, err := addr.MarshalBinary()
	if err != nil {
		return err
	}
	err = rdb.ZAdd(ctx, redisConfig.Key, &redis.Z{
		Score:  SCORE_MAX,
		Member: member,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

// /////// api service func ///////

// get count of jwglxt
func count() int64 {
	count, err := rdb.ZCard(ctx, redisConfig.Key).Result()
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
	err = rdb.ZScore(ctx, redisConfig.Key, string(buf)).Err()
	return err == nil
}

// get batch of jwglxtS
func batch(cursor uint64, match string, count int64) (res []string, retCursor uint64, err error) {
	res, retCursor, err = rdb.ZScan(ctx, redisConfig.Key, cursor, match, count).Result()
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
}

// get all proxies
func all() (proxies []string, err error) {
	proxies, err = rdb.ZRangeByScore(ctx, redisConfig.Key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", SCORE_MIN),
		Max: fmt.Sprintf("%d", SCORE_MAX),
	}).Result()
	if err != nil {
		return proxies, err
	}

	return proxies, err
}
