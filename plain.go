package charon

import (
	"context"
	"fmt"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	debug              = false
)

// unaryInterceptor is an example unary interceptor.
func UnaryInterceptorClient(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var header metadata.MD // variable to store header and trailer
	err := invoker(ctx, method, req, reply, cc, grpc.Header(&header))
	// log the header's info
	logger("[Received Resp]:	The metadata for response is %s\n", header)
	return err
}

// unaryInterceptor is an example unary interceptor.
func UnaryInterceptorEnduser(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return errMissingMetadata
	}

	methodName := md["method"][0]
	// print all the k-v pairs in the metadata md
	if debug {
		logger("[Before Req]:	Node Client is calling %s\n", methodName)
		var metadataLog string
		for k, v := range md {
			metadataLog += fmt.Sprintf("%s: %s, ", k, v)
		}
		if metadataLog != "" {
			logger("[Sending Req Enduser]: The metadata for request is %s\n", metadataLog)
		}
	}

	// Set a timer for the client
	startTime := time.Now()
	ctx = metadata.AppendToOutgoingContext(ctx, "time", startTime.Format(time.RFC3339Nano))

	var header metadata.MD // variable to store header and trailer
	err := invoker(ctx, method, req, reply, cc, grpc.Header(&header))
	// log the header's info
	logger("[Received Resp Enduser]:	The metadata for response is %s\n", header)
	return err
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// This is the server side interceptor, it should check tokens, update price, do overload handling and attach price to response

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if debug {
		var metadataLog string
		for k, v := range md {
			metadataLog += fmt.Sprintf("%s: %s, ", k, v)
		}
		if metadataLog != "" {
			logger("[Received Req]: The metadata for request is %s\n", metadataLog)
		}
	}

	m, err := handler(ctx, req)

	if err != nil {
		logger("RPC failed with error %v", err)
	}
	return m, err
}
