package proto

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

type MyManagerServer struct{}

func (m *MyManagerServer) GetManager(ctx context.Context, req *ManagerReq) *ResponseManager {
	res := NewResponseManager()
	res.SetData(&Manager{
		Name: "Jim",
	})
	return res
}

func TestManager(t *testing.T) {
	mux := http.NewServeMux()
	InitRouter(mux, &MyManagerServer{})
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
