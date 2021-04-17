package common

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
	"time"
	"context"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// SetUser save id and current time
func SetUser(w http.ResponseWriter, r *http.Request, uid int) {
	id := strconv.Itoa(uid)
	now := time.Now().Format("2006-01-02 15:04:05")

	value, _ := json.Marshal([]string{id,now,now})

	rand := StringRand(20)

	rdb := redis.NewClient(&redis.Options{
        Addr:     RedisConnect,
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    err := rdb.Set(ctx, rand, string(value), 60000000000 * 60 * 24 * 3).Err()
    if err != nil {
        panic(err)
    }

    fmt.Println("key", string(value))
	fmt.Println("rand", rand)

	cookie := &http.Cookie{
		Name:     "ss",
		Value:    rand,
		MaxAge:   86400 * 3,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

// GetUser check value exsist not expired
func GetUser(w http.ResponseWriter, r *http.Request) int {
	cookie, err := r.Cookie("ss")
	if err != nil {
		return 0
	}
	rdb := redis.NewClient(&redis.Options{
        Addr:     RedisConnect,
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    js, err := rdb.Get(ctx, cookie.Value).Result()
    if err != nil {
        panic(err)
    }

	log.Print("redis json: ", js)
	var d = []string{}
    if err := json.Unmarshal([]byte(js), &d); err != nil {
        panic(err)
    }
    fmt.Println(d[1])
	stampTime, err := time.Parse("2006-01-02 15:04:05", d[1])
	if err != nil {
		log.Print("time parse: ", err)
	}
	addDays := stampTime.AddDate(0, 0, 30)
	fmt.Println(addDays)
	if time.Now().After(addDays) {
		log.Print(d[0] + "session expired")
		return 0
	}
	stampTime, err = time.Parse("2006-01-02 15:04:05", d[2])
	addDays = stampTime.AddDate(0, 0, 2)
	fmt.Println(addDays)
	if time.Now().After(addDays) { // regenerate session id
		rdb.Del(ctx, cookie.Value)
		rand := StringRand(20)
		now := time.Now().Format("2006-01-02 15:04:05")
		value, _ := json.Marshal([]string{d[0],d[1],now})
		err := rdb.Set(ctx, rand, value, 60000000000 * 60 * 24 * 3).Err()
		if err != nil {
			panic(err)
		}
		fmt.Println("key", string(value))
		fmt.Println("rand", rand)	
		cookie := &http.Cookie{
			Name:     "ss",
			Value:    rand,
			MaxAge:   86400 * 3,
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	}
	usrID, err := strconv.Atoi(d[0])
	if err != nil {
		panic(err)
	}
	return usrID
}
