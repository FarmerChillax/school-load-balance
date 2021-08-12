package storage

import (
	"balance/network"
	"encoding/json"
	"log"
	"net/http"
)

func writeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var segmentRes network.SegmentResult
		err := decoder.Decode(&segmentRes)
		log.Println(segmentRes.Host, segmentRes.Ports)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// 写数据库
		err = writeDB(segmentRes.Host, segmentRes.Ports)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func pingHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		respData := pingDB()
		rw.Write([]byte(respData))
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func RegisterHandlers() {
	http.HandleFunc("/ping", pingHandler)

	http.HandleFunc("/write", writeHandler)

}

// var ctx = context.Background()

// func ExampleClient() {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "192.168.2.122:6379",
// 		Password: "farmer233", // no password set
// 		DB:       5,           // use default DB
// 	})
// 	fmt.Println(rdb.Ping(ctx))
// 	key := "farmer"

// 	farmerVal, err := rdb.Get(ctx, key).Result()
// 	if err == redis.Nil {
// 		fmt.Printf("key: %s not exist\n", key)
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Printf("Key: farmer; value: %v\n", farmerVal)
// 	}

// 	err = rdb.Set(ctx, key, fmt.Sprintf("%s%d", farmerVal, 233), 0).Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	val, err := rdb.Get(ctx, "key").Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("key", val)

// 	val2, err := rdb.Get(ctx, "key2").Result()
// 	if err == redis.Nil {
// 		fmt.Println("key2 does not exist")
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println("key2", val2)
// 	}
// 	// Output: key value
// 	// key2 does not exist
// }

// func main() {
// 	ExampleClient()
// }
