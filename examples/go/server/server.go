/*
 * Copyright 2025 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"strconv"
	"time"

	"examples/proto"

	"github.com/go-spring/stdlib/httpsvr"
	"github.com/go-spring/stdlib/ptrutil"
)

type ManagerServer struct{}

func (m *ManagerServer) GetManager(ctx context.Context, req *proto.ManagerReq) *proto.GetManagerResp {
	return &proto.GetManagerResp{
		Data: &proto.Manager{
			Name:  ptrutil.New("Jim"),
			Level: ptrutil.New(proto.ManagerLevelAsString(proto.ManagerLevel_JUNIOR)),
		},
	}
}

func (m *ManagerServer) CreateManager(ctx context.Context, req *proto.CreateManagerReq) *proto.CreateManagerResp {
	return nil
}

func (m *ManagerServer) UpdateManager(ctx context.Context, req *proto.UpdateManagerReq) map[string]any {
	return nil
}

func (m *ManagerServer) DeleteManager(ctx context.Context, req *proto.ManagerReq) *proto.DeleteManagerResp {
	return nil
}

func (m *ManagerServer) ListManagersByPage(ctx context.Context, req *proto.ListManagersByPageReq) *proto.ListManagersByPageResp {
	return nil
}

func (m *ManagerServer) Assistant(ctx context.Context, req *proto.AssistantReq, resp chan<- *httpsvr.Event[*proto.AssistantResp]) {
	for i := 0; i < 5; i++ {
		event := httpsvr.NewEvent[*proto.AssistantResp]().
			ID(strconv.Itoa(i)).
			Event("message").
			Data(&proto.AssistantResp{
				Id: ptrutil.New(strconv.Itoa(i)),
				Payload: ptrutil.New(proto.Payload{
					FieldType:  proto.PayloadTypeAsString(proto.PayloadType_PayloadOf1),
					PayloadOf1: ptrutil.New(proto.PayloadOf1{}),
				}),
				Image: []byte("000111222333444555666777888999000"),
			})
		resp <- event
		time.Sleep(time.Second)
	}
}

func (m *ManagerServer) AssistantV2(ctx context.Context, req *proto.AssistantReq, resp chan<- *httpsvr.Event[string]) {
	for i := 0; i < 5; i++ {
		resp <- httpsvr.NewEvent[string]().
			ID(strconv.Itoa(i)).
			Data("123456")
		time.Sleep(time.Second)
	}
}
