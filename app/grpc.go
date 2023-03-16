package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net"
	"strconv"
)

var (
	grpcServer *grpc.Server
	grpcMux    *runtime.ServeMux
)

// GRPC
func InitGRPCServer() (*grpc.Server, *runtime.ServeMux, error) {

	var options []grpc.ServerOption

	// TSL
	tslEnable := viper.GetString("app.grpc.tsl") == "true"
	if tslEnable {
		crt := "./cert/service.pem"
		key := "./cert/service.key"
		caN := "./cert/ca.cert"

		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return nil, nil, errors.New("cannot initialize GRPC Server")
		}

		// CreateIfNotExists a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caN)
		if err != nil {
			return nil, nil, errors.New("cannot initialize GRPC Server")
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, nil, errors.New("failed to append client certs")
		}

		// CreateIfNotExists the TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})
		options = append(options, grpc.Creds(creds))
	}

	// Middleware
	options = append(options, []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(errorLogging),
	}...)

	// CreateIfNotExists server
	grpcServer = grpc.NewServer(options...)
	if grpcServer == nil {
		return nil, nil, errors.New("cannot initialize GRPC Server")
	}

	grpcMux = runtime.NewServeMux()

	return grpcServer, grpcMux, nil
}

func RunGRPCServer() {

	gRPCPort := viper.GetString("app.grpc.port")

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", "localhost:"+gRPCPort)
	if err != nil {
		log.Fatal("%v", err)
	}
	log.Info("GRPC server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("%v", err)
	}
}

// Register

type DeliveryService interface {
	Init() error
}

func InitDelivery(d ...DeliveryService) error {
	for _, service := range d {
		err := service.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func InitGRPCService[T any](s func(grpc.ServiceRegistrar, T), src T) {
	s(grpcServer, src)
}

// Logging interceptor

func errorLogging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// Calls the handler
	h, err := handler(ctx, req)

	// Log if error
	if err != nil {
		log.Error("%v", err)
		return h, status.Error(codes.Internal, err.Error())
	}

	return h, nil
}

// Tools

func ExtractRequestUserId(ctx context.Context) (int32, error) {
	m, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userIds := m.Get("user_id")
		if len(userIds) > 0 {
			userId, err := strconv.ParseInt(userIds[0], 10, 32)
			if err != nil {
				return -1, errors.Wrap(err, "cannot parse user_id")
			}
			return int32(userId), nil
		}
	}
	return -1, errors.New("user_id was not found into context")
}
