package router

import (
	"appjam/domain"
	"appjam/response"
	"appjam/util"
	"strconv"

	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetPosts(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	var postList []domain.PostList
	data, err := util.DB.Query("SELECT post_id, u.user_name, title, created FROM post INNER JOIN \"user\" u ON u.user_id = post.owner_id ") // 추후 page searching도 만들어야함.
	if err != nil {
		if err == sql.ErrNoRows {
			response.GlobalErr("no data", nil, 404, w)
		} else {
			response.GlobalErr("select error", err, 400, w)
		}
		return
	}

	for data.Next() {
		var post domain.PostList
		data.Scan(&post.PostId, &post.UserName, &post.Title, &post.Created)

		postList = append(postList, post)
	}

	resData, _ := json.Marshal(response.Res{
		Data: postList,
		Err:  false,
	})

	fmt.Fprint(w, string(resData))
}

func PostDetail(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	postId := p.ByName("postId")
	if _, err := strconv.Atoi(postId); len(postId) < 1 || err != nil {
		response.GlobalErr("not enough params", nil, 500, w)
		return
	}

	var postDetail domain.Post
	err := util.DB.QueryRow("SELECT u.user_name, p.title, p.created FROM post p INNER JOIN \"user\" u ON u.user_id=p.owner_id WHERE post_id=$1;", postId).
		Scan(&postDetail.OwnerName, &postDetail.Title, &postDetail.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			response.GlobalErr("no data", nil, 404, w)
		} else {
			response.GlobalErr("select error", err, http.StatusBadRequest, w)
		}
		return
	}

	resData, _ := json.Marshal(response.Res{
		Data: postDetail,
		Err:  false,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}

func DeletePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var userId interface{}
	if userId = util.LoginCheck(r); userId == nil {
		response.GlobalErr("need login", nil, 401, w)
		return
	}

	postId := p.ByName("postId")
	if numPostId, err := strconv.Atoi(postId); err != nil || numPostId < 1 {
		response.GlobalErr("wrong post_id", nil, 400, w)
		return
	}
	_, err := util.DB.Exec(`DELETE from post WHERE owner_id=$1 AND post_id=$2`, userId, postId)
	if err != nil {
		response.GlobalErr("cannot delete post", err, 500, w)
		return
	}

	resData, _ := json.Marshal(response.Res{
		Data: "delete success",
		Err:  false,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}

func WritePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userId interface{}

	if userId = util.LoginCheck(r); userId == nil {
		response.GlobalErr("need login", nil, http.StatusForbidden, w)
		return
	}

	// post data
	var pd domain.WritePost

	err := json.NewDecoder(r.Body).Decode(&pd)
	if err != nil {
		response.GlobalErr("body data wrong", err, 400, w)
		return
	}

	postId, err := util.DB.Exec(`INSERT INTO public.post(
		title, description, owner_id) VALUES ($1, $2, $3);`,
		pd.Title, pd.Description, userId)

	if err != nil {
		response.GlobalErr("inserting err", err, 500, w)
		return
	}
	resData, _ := json.Marshal(response.Res{
		Data: postId,
		Err:  false,
	})
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(resData))
}
