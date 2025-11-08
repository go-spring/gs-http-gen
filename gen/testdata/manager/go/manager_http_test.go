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
	"github.com/go-spring/gs-http-gen/lib/httputil"
)

type MyManagerServer struct{}

func (m *MyManagerServer) GetManager(ctx context.Context, req *proto.ManagerReq) *proto.GetManagerResp {
	return &proto.GetManagerResp{
		Data: &proto.Manager{
			Name:  httputil.Ptr("Jim"),
			Level: httputil.Ptr(proto.ManagerLevelAsString(proto.ManagerLevel_JUNIOR)),
		},
	}
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
			Id: httputil.Ptr(strconv.Itoa(i)),
			Payload: httputil.Ptr(proto.Payload{
				FieldType: httputil.Ptr(proto.PayloadTypeAsString(proto.PayloadType_text_data)),
				TextData:  httputil.Ptr("123"),
			}),
		}
		time.Sleep(time.Second)
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
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1025))
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	fmt.Print(string(b))
	svr.Shutdown(t.Context())
}
