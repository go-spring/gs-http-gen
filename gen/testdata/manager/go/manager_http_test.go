package go_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/go-spring/gs-http-gen/gen/testdata/manager/go/proto"
)

type MyManagerServer struct{}

func (m *MyManagerServer) GetManager(ctx context.Context, req *proto.ManagerReq) *proto.ResponseManager {
	res := proto.NewResponseManager()
	res.SetData(&proto.Manager{
		Name: "Jim",
	})
	return res
}

func TestManager(t *testing.T) {
	mux := http.NewServeMux()
	proto.InitRouter(mux, &MyManagerServer{})
	go func() {
		fmt.Println(http.ListenAndServe(":9191", mux))
	}()
	time.Sleep(time.Millisecond * 300)

	resp, err := http.Get("http://localhost:9191/user")
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	fmt.Println(string(b))
}
