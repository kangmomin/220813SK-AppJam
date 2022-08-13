package response

type Res struct {
	Data any  `json:"data"`
	Err  bool `json:"error"`
}
