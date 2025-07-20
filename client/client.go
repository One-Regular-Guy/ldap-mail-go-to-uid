package main

import (
	"context"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	pb "github.com/One-Regular-Guy/ldap-mail-go-to-uid/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	// Carrega o certificado do servidor
	cert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		log.Fatalf("Erro ao ler cert.pem: %v", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	creds := credentials.NewClientTLSFromCert(certPool, "")

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewSecureServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "meu-secret")

	resp, err := client.SecureEcho(ctx, &pb.EncryptedRequest{Payload: "email@exemplo.com"})
	if err != nil {
		log.Fatalf("Erro ao chamar SecureEcho: %v", err)
	}
	log.Printf("Resposta: %s", resp.Payload)
}
