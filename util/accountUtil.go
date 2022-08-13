package util

import (
	"context"
	"log"
	"net/http"
)

func LoginCheck(r *http.Request) interface{} {
	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Println(nil)
		}
		return nil
	}

	data, err := Rdb.Get(context.Background(), sessionID.Value).Result()
	if err != nil {
		log.Println(err)
		return nil
	}

	return data
}
