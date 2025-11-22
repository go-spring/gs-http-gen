package go_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-spring/gs-http-gen/gen/testdata/manager/go/proto"
	"github.com/go-spring/gs-http-gen/lib/httputil"
	"github.com/go-spring/gs-http-gen/lib/pathidl"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// GinServer defines the interface that service must implement.
type GinServer struct {
	*http.Server
	engine *gin.Engine
}

// NewGinServer creates a new GinServer instance.
func NewGinServer(addr string) *GinServer {
	engine := gin.New()
	svr := &http.Server{Addr: addr, Handler: engine.Handler()}
	return &GinServer{Server: svr, engine: engine}
}

// ToGinPath converts a pathidl.Path to a Gin compatible path.
func ToGinPath(pattern string) string {
	path, _ := pathidl.Parse(pattern)
	var sb strings.Builder
	for _, s := range path {
		sb.WriteString("/")
		switch s.Type {
		case pathidl.Static:
			sb.WriteString(s.Value)
		case pathidl.Param:
			sb.WriteString(":")
			sb.WriteString(s.Value)
		case pathidl.Wildcard:
			sb.WriteString("*")
			sb.WriteString(s.Value)
		}
	}
	return sb.String()
}

// HandleFunc registers a new route for the given HTTP method and pattern.
func (s *GinServer) HandleFunc(method string, pattern string, handler http.HandlerFunc) {
	pattern = ToGinPath(pattern)
	switch method {
	case "GET":
		s.engine.GET(pattern, gin.WrapF(handler))
	case "POST":
		s.engine.POST(pattern, gin.WrapF(handler))
	case "PUT":
		s.engine.PUT(pattern, gin.WrapF(handler))
	case "DELETE":
		s.engine.DELETE(pattern, gin.WrapF(handler))
	case "HEAD":
		s.engine.HEAD(pattern, gin.WrapF(handler))
	default:
		panic(fmt.Sprintf("unsupported method: %s", method))
	}
}

// HttpServer defines the interface that service must implement.
type HttpServer struct {
	*http.Server
	mux *http.ServeMux
}

// NewHttpServer creates a new HttpServer instance.
func NewHttpServer(addr string) *HttpServer {
	mux := http.NewServeMux()
	svr := &http.Server{Addr: addr, Handler: mux}
	return &HttpServer{Server: svr, mux: mux}
}

// HandleFunc registers a new route for the given HTTP method and pattern.
func (s *HttpServer) HandleFunc(method string, pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(strings.TrimSpace(method+" "+pattern), handler)
}

type MyManagerServer struct{}

var _ proto.ManagerServer = (*MyManagerServer)(nil)

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

func (m *MyManagerServer) Assistant(ctx context.Context, req *proto.AssistantReq, resp chan<- *proto.SSEEvent[*proto.AssistantResp]) {
	for i := 0; i < 5; i++ {
		event := proto.NewSSEEvent[*proto.AssistantResp]().ID("1").Event("message").Data(
			&proto.AssistantResp{
				Id: httputil.Ptr(strconv.Itoa(i)),
				Payload: httputil.Ptr(proto.Payload{
					FieldType: httputil.Ptr(proto.PayloadTypeAsString(proto.PayloadType_Payload_1)),
					Payload1:  httputil.Ptr(proto.Payload_1{}),
				}),
				Image: []byte("000111222333444555666777888999000"),
			})
		resp <- event
		time.Sleep(time.Second)
	}
}

func TestManager(t *testing.T) {
	svr := NewGinServer(":9191")
	proto.SetupRouter(svr, &MyManagerServer{})
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
	svr := NewHttpServer(":9191")
	proto.SetupRouter(svr, &MyManagerServer{})
	go func() {
		fmt.Println(svr.ListenAndServe())
	}()
	time.Sleep(time.Millisecond * 300)

	resp, err := http.Get("http://localhost:9191/assistant")
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
