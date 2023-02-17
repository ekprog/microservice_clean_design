package tests

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"microservice_clean_design/domain"
	pb "microservice_clean_design/pkg/pb/api"
	"os"
	"testing"
	"time"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
}

func makeClient() *grpc.ClientConn {

	addr := "localhost:" + os.Getenv("GRPC_PORT")

	tslEnable := os.Getenv("TSL_ENABLE") == "true"
	if tslEnable {

		crt := "../cert/ca.cert"
		key := "../cert/ca.key"
		caN := "../cert/ca.cert"

		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			log.Fatalf("could not load client key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caN)
		if err != nil {
			log.Fatalf("could not read ca certificate: %s", err)
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("failed to append ca certs")
		}

		creds := credentials.NewTLS(&tls.Config{
			ServerName:   addr, // NOTE: this is required!
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})

		// Create a connection with the TLS credentials
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Fatalf("could not dial %s: %s", addr, err)
		}
		return conn
	} else {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return conn
	}
}

func Test_CreateTaskGRPC(t *testing.T) {

	conn := makeClient()
	defer conn.Close()
	c := pb.NewTasksServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.CreateTask(ctx, &pb.CreateTaskRequest{
		Name: "New Task",
	})
	require.NoError(t, err, "should be success creating ucase")
	require.Equal(t, res.Status.Code, domain.Success, "should be success status code")
	log.Infof("%v", res)
}

func Test_GetAllTasksGRPC(t *testing.T) {

	conn := makeClient()
	defer conn.Close()
	c := pb.NewTasksServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.AllTasks(ctx, &pb.AllTasksRequest{})
	require.NoError(t, err, "should be success all_tasks ucase")
	require.Equal(t, res.Status.Code, domain.Success, "should be success status code")
	log.Infof("%v", res)
}
