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
	"github.com/go-spring/gs-http-gen/lib/httpsvr"
	"github.com/go-spring/gs-http-gen/lib/pathidl"
	"github.com/lvan100/golib/ptrutil"
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

type MyManagerServer struct{}

var _ proto.ManagerService = (*MyManagerServer)(nil)

func (m *MyManagerServer) GetManager(ctx context.Context, req *proto.ManagerReq) *proto.GetManagerResp {
	return &proto.GetManagerResp{
		Data: &proto.Manager{
			Name:  ptrutil.New("Jim"),
			Level: ptrutil.New(proto.ManagerLevelAsString(proto.ManagerLevel_JUNIOR)),
		},
	}
}

func (m *MyManagerServer) CreateManager(ctx context.Context, req *proto.CreateManagerReq) *proto.CreateManagerResp {
	return nil
}

func (m *MyManagerServer) UpdateManager(ctx context.Context, req *proto.UpdateManagerReq) map[string]any {
	return nil
}

func (m *MyManagerServer) DeleteManager(ctx context.Context, req *proto.ManagerReq) *proto.DeleteManagerResp {
	return nil
}

func (m *MyManagerServer) ListManagersByPage(ctx context.Context, req *proto.ListManagersByPageReq) *proto.ListManagersByPageResp {
	return nil
}

func (m *MyManagerServer) Assistant(ctx context.Context, req *proto.AssistantReq, resp chan<- *httpsvr.Event[*proto.AssistantResp]) {
	for i := 0; i < 5; i++ {
		event := httpsvr.NewEvent[*proto.AssistantResp]().
			ID(strconv.Itoa(i)).
			Event("message").
			Data(&proto.AssistantResp{
				Id: ptrutil.New(strconv.Itoa(i)),
				Payload: ptrutil.New(proto.Payload{
					FieldType: proto.PayloadTypeAsString(proto.PayloadType_Payload_1),
					Payload1:  ptrutil.New(proto.Payload_1{}),
				}),
				Image: []byte("000111222333444555666777888999000"),
			})
		resp <- event
		time.Sleep(time.Second)
	}
}

func (m *MyManagerServer) AssistantV2(ctx context.Context, req *proto.AssistantReq, resp chan<- *httpsvr.Event[string]) {
	for i := 0; i < 5; i++ {
		resp <- httpsvr.NewEvent[string]().
			ID(strconv.Itoa(i)).
			Data("123456")
		time.Sleep(time.Second)
	}
}

func TestManager(t *testing.T) {
	svr := NewGinServer(":9191")
	for _, r := range proto.Routers(&MyManagerServer{}) {
		svr.HandleFunc(r.Method, r.Pattern, r.Handler)
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
	svr := httpsvr.NewSimpleServer(":9191")
	for _, r := range proto.Routers(&MyManagerServer{}) {
		svr.Route(r)
	}
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

func TestStreamV2(t *testing.T) {
	svr := httpsvr.NewSimpleServer(":9191")
	for _, r := range proto.Routers(&MyManagerServer{}) {
		svr.Route(r)
	}
	go func() {
		fmt.Println(svr.ListenAndServe())
	}()
	time.Sleep(time.Millisecond * 300)

	resp, err := http.Get("http://localhost:9191/assistantV2")
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

/*

type Item struct {
	ID int64 `json:"id"`
}

type Object struct {
	Item *Item  `json:"item"`
	Text string `json:"text"`
}

type HelloRequest struct {
	HelloRequestBody
	Int             int               `json:"int" query:"int"`
	String          string            `json:"string" query:"string"`
	IntPtr          *int              `json:"int_ptr" query:"int_ptr"`
	StringPtr       *string           `json:"string_ptr" query:"string_ptr"`
	IntSlice        []int             `json:"int_slice" query:"int_slice"`
	StringSlice     []string          `json:"string_slice" query:"string_slice"`
	ByteSlice       []byte            `json:"byte_slice" query:"byte_slice"`
	Object          *Object           `json:"object" query:"object"`
	ObjectSlice     []Object          `json:"object_slice" query:"object_slice"`
	IntStringMap    map[int]string    `json:"int_string_map" query:"int_string_map"`
	StringObjectMap map[string]Object `json:"string_object_map" query:"string_object_map"`
}

func (req *HelloRequest) QueryString() (string, error) {
	m := url.Values{}

	// Encode scalar values using the key format (e.g. a=1)
	m.Add("int", strconv.FormatInt(int64(req.Int), 10))
	m.Add("string", req.String)
	if req.IntPtr != nil {
		m.Add("int_ptr", strconv.FormatInt(int64(*req.IntPtr), 10))
	}
	if req.StringPtr != nil {
		m.Add("string_ptr", *req.StringPtr)
	}

	// Encode arrays using the repeated key format (e.g. a=1&a=2)
	for _, v := range req.IntSlice {
		m.Add("int_slice", strconv.FormatInt(int64(v), 10))
	}
	// Encode arrays using the repeated key format (e.g. a=1&a=2)
	for _, v := range req.StringSlice {
		m.Add("string_slice", v)
	}

	// Encode byte slices using base64 encoding (e.g., a=YWJj)
	if req.ByteSlice != nil {
		m.Add("byte_slice", base64.StdEncoding.EncodeToString(req.ByteSlice))
	}

	// Encode an array of objects using repeated keys with JSON values
	// e.g. items={"id":1,"name":"A"}&items={"id":2,"name":"B"}
	for _, v := range req.ObjectSlice {
		b, err := jsonflow.Marshal(v)
		if err != nil {
			return "", err
		}
		m.Add("object_slice", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"})
	if req.Object != nil {
		b, err := jsonflow.Marshal(req.Object)
		if err != nil {
			return "", err
		}
		m.Add("object", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"})
	if req.StringObjectMap != nil {
		b, err := jsonflow.Marshal(req.StringObjectMap)
		if err != nil {
			return "", err
		}
		m.Add("string_object_map", string(b))
	}

	// Encode maps or structs as JSON strings (e.g. data={"id":1,"name":"Alice"})
	if req.IntStringMap != nil {
		b, err := jsonflow.Marshal(req.IntStringMap)
		if err != nil {
			return "", err
		}
		m.Add("int_string_map", string(b))
	}

	return m.Encode(), nil
}

type HelloRequestBody struct{}

*/
