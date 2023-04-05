package main

import (
	"context"   // Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
	"log"       // Package log implements a simple logging package. It defines a type, Logger, with methods for formatting output.
	"math/rand" // Package rand implements pseudo-random number generators.
	"net"       // Package net provides a portable interface for network I/O, including TCP/IP, UDP, domain name resolution, and Unix domain sockets.

	// pb is a package auto-generated by the protocol buffer compiler from the usermgmt.proto file.
	pb "github.com/alaiy95/golang-projects/grpc-unary-service/usermgmt"
	"google.golang.org/grpc" // Package grpc implements an RPC system called gRPC.
)

const (
	port = ":50051" // The port on which the server listens for incoming requests.
)

// UserManagementServer implements the pb.UserManagementServer interface.
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer // We embed the UnimplementedUserManagementServer to satisfy the interface's methods.
}

// CreateNewUser creates a new user with a randomly generated ID.
func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())                               // Log the name of the received user.
	userID := int32(rand.Intn(100))                                        // Generate a random user ID.
	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: userID}, nil // Return a User message with the received name and age, and the generated ID.
}

func main() {
	lis, err := net.Listen("tcp", port) // Create a TCP listener on the specified port.
	if err != nil {
		log.Fatalf("failed to listen: %v", err) // If an error occurs, log it and exit.
	}
	s := grpc.NewServer()                                       // Create a new gRPC server.
	pb.RegisterUserManagementServer(s, &UserManagementServer{}) // Register our implementation of the UserManagementServer interface with the server.
	log.Printf("server listening at %v", lis.Addr())            // Log that the server is listening on the address.
	if err := s.Serve(lis); err != nil {                        // Start serving the gRPC server on the listener. If an error occurs, log it and exit.
		log.Fatalf("failed to serve: %v", err)
	}
}
