package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/lvan100/golib/httputil"
)

func TestClient(t *testing.T) {
	c := &Client{&ClientImpl{
		Client: &httputil.DefaultClient{
			Client: http.DefaultClient,
			Scheme: "http",
			Host:   "127.0.0.1:9090",
		},
	}}
	go func() {
		http.HandleFunc("/v1/ping", func(w http.ResponseWriter, r *http.Request) {
			resp := &PingResp{
				Data: "hello world",
			}
			json.NewEncoder(w).Encode(resp)
		})
		http.ListenAndServe(":9090", nil)
	}()
	time.Sleep(time.Millisecond * 50)
	_, s, err := c.Ping(context.Background(), &PingReq{
		PingReqBody: PingReqBody{
			Name: "abc",
		},
		Name: "123",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}
