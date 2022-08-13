package router

import (
	"appjam/domain"
	"appjam/response"
	"appjam/util"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/argon2"
)

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if userId := util.LoginCheck(r); userId != nil {
		response.GlobalErr("already login", nil, 400, w)
		return
	}

	var loginData domain.Login
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		response.GlobalErr("data isn't json", err, 400, w)
		return
	}

	var confirmData domain.ConfirmLoginData
	var userId int
	err = util.DB.QueryRow("SELECT user_id, password, salt FROM \"user\" WHERE login_id=$1", loginData.LoginId).
		Scan(&userId, &confirmData.Password, &confirmData.Salt)

	if err != nil {
		response.GlobalErr("id error", err, http.StatusUnauthorized, w)
		return
	}

	confirmData.DecodeSalt, _ = hex.DecodeString(confirmData.Salt)

	encodedPwd := hex.EncodeToString(
		argon2.IDKey([]byte(loginData.Password), confirmData.DecodeSalt,
			argonConfig.Time, argonConfig.Memory, argonConfig.Thread, argonConfig.KeyLen))

	if encodedPwd != confirmData.Password {
		response.GlobalErr("password error", err, http.StatusUnauthorized, w)
		return
	}

	sessionID := [8]byte{}
	rand.Read(sessionID[:])

	_, err = util.Rdb.Set(ctx, hex.EncodeToString(sessionID[:]), userId, 0).Result()
	if err != nil {
		response.GlobalErr("generate session error", err, 500, w)
		return
	}

	data, _ := util.Rdb.Get(ctx, hex.EncodeToString(sessionID[:])).Result()
	fmt.Println(data)
	fmt.Println(hex.EncodeToString(sessionID[:]))

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    hex.EncodeToString(sessionID[:]),
		HttpOnly: true,
	})
	resData, _ := json.Marshal(response.Res{
		Data: "login sucess",
		Err:  false,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}

func SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var signUpData domain.SignUp
	err := json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		response.GlobalErr("data isn't json", err, 400, w)
		return
	}

	if !signUpData.IsValidLen() {
		response.GlobalErr("data is not enough", nil, 400, w)
		return
	}

	salt := make([]byte, 32)
	rand.Read(salt)
	encryptedPwd := argon2.IDKey([]byte(signUpData.Password), salt, argonConfig.Time, argonConfig.Memory, argonConfig.Thread, argonConfig.KeyLen)

	_, err = util.DB.Exec("INSERT INTO public.user (user_name, login_id, password, salt) VALUES ($1, $2, $3, $4);",
		signUpData.UserName, signUpData.LoginId, hex.EncodeToString(encryptedPwd), hex.EncodeToString(salt))

	if err != nil {
		response.GlobalErr("cannot sign up", err, 400, w)
		return
	}

	resData, _ := json.Marshal(response.Res{
		Data: "success",
		Err:  false,
	})

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(resData))
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if userId := util.LoginCheck(r); userId == nil {
		response.GlobalErr("didn't login", nil, 400, w)
		return
	}
	c, err := r.Cookie("sessionID")
	if err != nil {
		response.GlobalErr("cannot logout", err, 500, w)
		return
	}
	c.MaxAge = -1
	http.SetCookie(w, c)

	resData, _ := json.Marshal(response.Res{
		Data: "logout success",
		Err:  false,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}
