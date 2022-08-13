package router

import (
	"appjam/domain"
	"appjam/response"
	"appjam/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Rank(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	row, err := util.DB.Query("SELECT user_name, post_count FROM \"user\" ORDER BY post_count ASC;")

	if err != nil {
		response.GlobalErr("selecting err", err, 500, w)
		return
	}

	var Ranking []domain.UserRank
	for row.Next() {
		var Rank domain.UserRank
		err := row.Scan(&Rank.UserName, &Rank.PostCount)
		if err != nil {
			continue
		}
		Ranking = append(Ranking, Rank)
	}

	resData, _ := json.Marshal(response.Res{
		Data: Ranking,
		Err:  false,
	})
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(resData))
}
