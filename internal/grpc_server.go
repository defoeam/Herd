package keyvaluestore

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/credentials"

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
func (s *GRPCServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.KeyValue, error) {
	value, ok := s.kv.Get(req.Key)
	if !ok {
		return nil, fmt.Errorf("key not found: %s", req.Key)
	}

	return &pb.KeyValue{
		Key:   req.Key,
		Value: value,
	}, nil
}

// GetAll returns all items in the key-value store.
func (s *GRPCServer) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
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
func (s *GRPCServer) GetKeys(ctx context.Context, req *pb.GetKeysRequest) (*pb.GetKeysResponse, error) {
	keys := s.kv.GetKeys()
	return &pb.GetKeysResponse{
		Keys: keys,
	}, nil
}

// GetValues returns all values in the key-value store.
func (s *GRPCServer) GetValues(ctx context.Context, req *pb.GetValuesRequest) (*pb.GetValuesResponse, error) {
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
func (s *GRPCServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.kv.Set(req.Key, req.Value)

	return &pb.SetResponse{
		Item: &pb.KeyValue{
			Key:   req.Key,
			Value: req.Value,
		},
	}, nil
}

// Delete deletes an item in the key-value store by key.
func (s *GRPCServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	value, ok := s.kv.Clear(req.Key)
	if !ok {
		return nil, fmt.Errorf("key not found: %s", req.Key)
	}

	return &pb.DeleteResponse{
		DeletedItem: &pb.KeyValue{
			Key:   req.Key,
			Value: value,
		},
	}, nil
}

// DeleteAll deletes all items in the key-value store.
func (s *GRPCServer) DeleteAll(ctx context.Context, req *pb.DeleteAllRequest) (*pb.DeleteAllResponse, error) {
	if err := s.kv.ClearAll(); err != nil {
		return nil, fmt.Errorf("failed to clear all items: %v", err)
	}
	return &pb.DeleteAllResponse{}, nil
}

// StartGRPCServer starts a gRPC server on port 50051.
// If enableLogging is true, it initializes logging to the specified file with a rotation interval of 1 hour.
//
// Parameters:
//   - enableLogging: A boolean flag to enable or disable logging.
//
// Returns:
//   - error: An error if the server fails to start, initialize logging, or listen on the specified port.
func StartGRPCServer(enableLogging bool) error {
	server := NewGRPCServer()
	if enableLogging {
		if err := server.kv.InitLogging("/app/log/transaction.log", 1*time.Hour); err != nil {
			return fmt.Errorf("failed to initialize logging: %v", err)
		}
	}

	// Load TLS credentials
	creds, err := credentials.NewServerTLSFromFile("/app/certs/server.crt", "/app/certs/server.key")
	if err != nil {
		return fmt.Errorf("failed to load TLS credentials: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	pb.RegisterKeyValueServiceServer(s, server)
	log.Printf("Starting secured gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
