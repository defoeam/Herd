package keyvaluestore

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/defoeam/herd/api/proto"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedKeyValueServiceServer
	kv *KeyValueStore
}

// NewGRPCServer creates a new gRPC server with an empty key-value store.
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{
		kv: NewKeyValueStore(),
	}
}

// Get returns an item in the key-value store by key.
func (s *GRPCServer) Get(_ context.Context, req *pb.GetRequest) (*pb.KeyValue, error) {
	value, ok := s.kv.Get(req.GetKey())
	if !ok {
		return nil, fmt.Errorf("key not found: %s", req.GetKey())
	}

	return &pb.KeyValue{
		Key:   req.GetKey(),
		Value: value,
	}, nil
}

// GetAll returns all items in the key-value store.
func (s *GRPCServer) GetAll(_ context.Context, _ *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	data := s.kv.GetAll()
	items := make([]*pb.KeyValue, 0, len(data))

	for k, v := range data {
		items = append(items, &pb.KeyValue{
			Key:   k,
			Value: v,
		})
	}

	return &pb.GetAllResponse{Items: items}, nil
}

// GetKeys returns all keys in the key-value store.
func (s *GRPCServer) GetKeys(_ context.Context, _ *pb.GetKeysRequest) (*pb.GetKeysResponse, error) {
	keys := s.kv.GetKeys()
	return &pb.GetKeysResponse{
		Keys: keys,
	}, nil
}

// GetValues returns all values in the key-value store.
func (s *GRPCServer) GetValues(_ context.Context, _ *pb.GetValuesRequest) (*pb.GetValuesResponse, error) {
	values := s.kv.GetValues()
	byteValues := make([][]byte, len(values))
	for i, v := range values {
		byteValues[i] = []byte(v)
	}
	return &pb.GetValuesResponse{
		Values: byteValues,
	}, nil
}

// Set sets an item in the key-value store by key and value.
func (s *GRPCServer) Set(_ context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.kv.Set(req.GetKey(), req.GetValue())

	return &pb.SetResponse{
		Item: &pb.KeyValue{
			Key:   req.GetKey(),
			Value: req.GetValue(),
		},
	}, nil
}

// Delete deletes an item in the key-value store by key.
func (s *GRPCServer) Delete(_ context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	value, ok := s.kv.Delete(req.GetKey())
	if !ok {
		return nil, fmt.Errorf("key not found: %s", req.GetKey())
	}

	return &pb.DeleteResponse{
		DeletedItem: &pb.KeyValue{
			Key:   req.GetKey(),
			Value: value,
		},
	}, nil
}

// DeleteAll deletes all items in the key-value store.
func (s *GRPCServer) DeleteAll(_ context.Context, _ *pb.DeleteAllRequest) (*pb.DeleteAllResponse, error) {
	if err := s.kv.DeleteALL(); err != nil {
		return nil, fmt.Errorf("failed to clear all items: %w", err)
	}
	return &pb.DeleteAllResponse{}, nil
}

// StartGRPCServer starts a gRPC server on port 50051.
// If enableLogging is true, it initializes logging to the specified file with a rotation interval of 1 hour.
func StartGRPCServer(enableLogging bool) error {
	server := NewGRPCServer()
	if enableLogging {
		if err := server.kv.InitLogging("/app/log/transaction.log", 1*time.Hour); err != nil {
			return fmt.Errorf("failed to initialize logging: %w", err)
		}
	}

	s := grpc.NewServer()

	lis, listenErr := net.Listen("tcp", "127.0.0.1:50051")
	if listenErr != nil {
		return fmt.Errorf("failed to listen: %w", listenErr)
	}

	pb.RegisterKeyValueServiceServer(s, server)
	log.Printf("Starting unsecured gRPC server on :50051")
	if serveErr := s.Serve(lis); serveErr != nil {
		return fmt.Errorf("failed to serve: %w", serveErr)
	}
	return nil
}
