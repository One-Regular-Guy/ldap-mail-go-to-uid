package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"

	pb "github.com/One-Regular-Guy/ldap-mail-go-to-uid/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const secretKey = "meu-secret"

func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) == 0 || md["authorization"][0] != secretKey {
		return nil, status.Error(codes.Unauthenticated, "invalid secret")
	}
	return handler(ctx, req)
}

type server struct {
	pb.UnimplementedSecureServiceServer
}

func (s *server) SecureEcho(ctx context.Context, req *pb.EncryptedRequest) (*pb.EncryptedResponse, error) {
	log.Printf("Recebido email: %s", req.Payload)
	uid := "uid123"
	return &pb.EncryptedResponse{Payload: uid}, nil
}

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Erro ao carregar certificados: %v", err)
	}
	creds := credentials.NewServerTLSFromCert(&cert)

	listener, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(authInterceptor),
	)
	pb.RegisterSecureServiceServer(s, &server{})
	log.Println("Running in :50051 (TLS)")
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("Could not serve: %v", err)
	}
}
