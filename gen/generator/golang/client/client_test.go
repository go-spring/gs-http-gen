package client

import (
	"context"
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient(nil)
	resp, p, err := c.Ping(context.Background(), &PingReq{})
	if err != nil {
		return
	}
}
