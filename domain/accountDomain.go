package domain

type Login struct {
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
}

type ConfirmLoginData struct {
	Password   string `json:"comfirm_password"`
	Salt       string `json:"confirm_salt"`
	DecodeSalt []byte `json:"decode_password"`
}

type SignUp struct {
	UserName string `json:"user_name"`
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
}

func (s SignUp) IsValidLen() (isOk bool) {
	if len(s.LoginId) < 2 ||
		len(s.Password) < 2 ||
		len(s.UserName) < 2 {
		return false
	}
	return true
}
