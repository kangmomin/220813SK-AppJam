package main

import (
	"appjam/router"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()

	// 자잘한 오류와 예외처리가 잘 안돼있지만 기능적으론 큰 문제가 없음.
	r.GET("/posts", router.GetPosts)              // test success
	r.POST("/posts", router.WritePost)            // test success
	r.GET("/post/:postId", router.PostDetail)     // test success
	r.DELETE("/posts/:postId", router.DeletePost) // 기능만은 구현 완료

	r.POST("/login", router.Login)    // test success
	r.POST("/logout", router.Logout)  // test success
	r.POST("/sign-up", router.SignUp) // test success

	http.ListenAndServe(":8080", r)
}
