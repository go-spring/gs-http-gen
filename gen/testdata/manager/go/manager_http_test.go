package go_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/go-spring/gs-http-gen/gen/testdata/manager/go/proto"
)

type MyManagerServer struct{}

func (m *MyManagerServer) GetManager(ctx context.Context, req *proto.ManagerReq) *proto.GetManagerResp {
	data := proto.NewManager()
	data.SetName("Jim")
	res := proto.NewGetManagerResp()
	res.SetData(data)
	return res
}

func (m *MyManagerServer) CreateManager(ctx context.Context, req *proto.CreateManagerReq) *proto.CreateManagerResp {
	return nil
}

func (m *MyManagerServer) UpdateManager(ctx context.Context, req *proto.UpdateManagerReq) *proto.UpdateManagerResp {
	return nil
}

func (m *MyManagerServer) DeleteManager(ctx context.Context, req *proto.ManagerReq) *proto.DeleteManagerResp {
	return nil
}

func (m *MyManagerServer) ListManagersByPage(ctx context.Context, req *proto.ListManagersByPageReq) *proto.ListManagersByPageResp {
	return nil
}

func (m *MyManagerServer) Stream(ctx context.Context, req *proto.StreamReq, resp chan<- *proto.StreamResp) {
	for i := 0; i < 5; i++ {
		resp <- &proto.StreamResp{
			Id: strconv.Itoa(i),
			Payload: proto.Payload{
				FieldType: proto.PayloadType_TextData,
				TextData:  "123",
			},
		}
	}
}

func TestManager(t *testing.T) {
	mux := http.NewServeMux()
	proto.InitRouter(mux, &MyManagerServer{})
	svr := &http.Server{
		Addr:    ":9191",
		Handler: mux,
	}
	go func() {
		fmt.Println(svr.ListenAndServe())
	}()
	time.Sleep(time.Millisecond * 300)

	resp, err := http.Get("http://localhost:9191/managers/123")
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	fmt.Println(string(b))
	svr.Shutdown(t.Context())
}

func TestStream(t *testing.T) {
	mux := http.NewServeMux()
	proto.InitRouter(mux, &MyManagerServer{})
	svr := &http.Server{
		Addr:    ":9191",
		Handler: mux,
	}
	go func() {
		fmt.Println(svr.ListenAndServe())
	}()
	time.Sleep(time.Millisecond * 300)

	resp, err := http.Get("http://localhost:9191/stream")
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	fmt.Print(string(b))
	svr.Shutdown(t.Context())
}
