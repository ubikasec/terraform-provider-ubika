// Code generated by protoc-gen-api. DO NOT EDIT.
package crypto

import (
	grpc "google.golang.org/grpc"
)

// Client aggregate all gRPC services' clients
type Client interface {
}

// client implements Client
type client struct {
}

func NewGRPCClient(conn *grpc.ClientConn) Client {
	return &client{}
}

// compile-time assertion
var _ Client = &client{}
