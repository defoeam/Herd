package keyvaluestore

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/defoeam/herd/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCServer struct {
	proto.UnimplementedKeyValueServiceServer
	kv *KeyValueStore
}

// NewGRPCServer creates a new gRPC server with an empty key-value store.
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{
		kv: NewKeyValueStore(),
	}
}

// Get returns an item in the key-value store by key.
func (s *GRPCServer) Get(_ context.Context, req *proto.GetRequest) (*proto.KeyValue, error) {
	value, ok := s.kv.Get(req.GetKey())
	if !ok {
		return nil, fmt.Errorf("key not found: %s", req.GetKey())
	}

	return &proto.KeyValue{
		Key:   req.GetKey(),
		Value: value,
	}, nil
}

// GetAll returns all items in the key-value store.
func (s *GRPCServer) GetAll(_ context.Context, _ *proto.GetAllRequest) (*proto.GetAllResponse, error) {
	data := s.kv.GetAll()
	items := make([]*proto.KeyValue, 0, len(data))

	for k, v := range data {
		items = append(items, &proto.KeyValue{
			Key:   k,
			Value: v,
		})
	}

	return &proto.GetAllResponse{Items: items}, nil
}

// GetKeys returns all keys in the key-value store.
func (s *GRPCServer) GetKeys(_ context.Context, _ *proto.GetKeysRequest) (*proto.GetKeysResponse, error) {
	keys := s.kv.GetKeys()
	return &proto.GetKeysResponse{
		Keys: keys,
	}, nil
}

// GetValues returns all values in the key-value store.
func (s *GRPCServer) GetValues(_ context.Context, _ *proto.GetValuesRequest) (*proto.GetValuesResponse, error) {
	values := s.kv.GetValues()
	byteValues := make([][]byte, len(values))
	for i, v := range values {
		byteValues[i] = []byte(v)
	}
	return &proto.GetValuesResponse{
		Values: byteValues,
	}, nil
}

// Set sets an item in the key-value store by key and value.
func (s *GRPCServer) Set(_ context.Context, req *proto.SetRequest) (*proto.SetResponse, error) {
	s.kv.Set(req.GetKey(), req.GetValue())

	return &proto.SetResponse{
		Item: &proto.KeyValue{
			Key:   req.GetKey(),
			Value: req.GetValue(),
		},
	}, nil
}

// Delete deletes an item in the key-value store by key.
func (s *GRPCServer) Delete(_ context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	value, ok := s.kv.Delete(req.GetKey())
	if !ok {
		return nil, fmt.Errorf("key not found: %s", req.GetKey())
	}

	return &proto.DeleteResponse{
		DeletedItem: &proto.KeyValue{
			Key:   req.GetKey(),
			Value: value,
		},
	}, nil
}

// DeleteAll deletes all items in the key-value store.
func (s *GRPCServer) DeleteAll(_ context.Context, _ *proto.DeleteAllRequest) (*proto.DeleteAllResponse, error) {
	if err := s.kv.DeleteALL(); err != nil {
		return nil, fmt.Errorf("failed to clear all items: %w", err)
	}
	return &proto.DeleteAllResponse{}, nil
}

// StartGRPCServer starts a gRPC server on port 50051.
// If enableLogging is true, it initializes logging to the specified file with a rotation interval of 1 hour.
func StartGRPCServer(enableLogging bool, enableSecurity bool) error {
	log.Printf("Starting server on port 7878...")

	// initialize the keyvalue store and logging
	server := NewGRPCServer()
	if enableLogging {
		if err := server.kv.InitLogging("/app/log/transaction.log", 1*time.Hour); err != nil {
			return fmt.Errorf("failed to initialize logging: %w", err)
		}
	}

	// create a new gRPC server with or without tls
	s, serverFactoryErr := grpcServerFactory(enableSecurity)
	if serverFactoryErr != nil {
		return fmt.Errorf("failed to create server: %w", serverFactoryErr)
	}

	// register the KeyValueService server
	proto.RegisterKeyValueServiceServer(s, server)

	// setup listener
	lis, listenErr := net.Listen("tcp", "0.0.0.0:7878")
	if listenErr != nil {
		return fmt.Errorf("failed to listen: %w", listenErr)
	}

	// serve the gRPC server
	serveErr := s.Serve(lis)
	if serveErr != nil {
		return fmt.Errorf("failed to serve: %w", serveErr)
	}
	return nil
}

// grpcServerFactory creates a new gRPC server with or without security enabled.
func grpcServerFactory(enableSecurity bool) (*grpc.Server, error) {
	if enableSecurity {
		// load the server's certificate and private key
		cert, certPairErr := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
		if certPairErr != nil {
			return nil, fmt.Errorf("failed to load X509 key pair: %w", certPairErr)
		}

		// setup and load the CA's certificate
		ca := x509.NewCertPool()
		caFilePath := "certs/ca.crt"
		caBytes, caBytesErr := os.ReadFile(caFilePath)
		if caBytesErr != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", caBytesErr)
		}
		if ok := ca.AppendCertsFromPEM(caBytes); !ok {
			return nil, fmt.Errorf("failed to append CA certificate")
		}

		// create a new TLS configuration with the server's certificate and the CA's certificate
		tlsConfig := &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    ca,
		}

		// create a new gRPC server with the TLS configuration
		return grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig))), nil
	}

	return grpc.NewServer(), nil
}
