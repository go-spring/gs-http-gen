package client

type PingReq struct {
	PingReqBody
	Name string `query:"name"`
}

type PingReqBody struct {
	Name string `json:"name"`
}

type PingResp struct {
	Data string `json:"data"`
}
