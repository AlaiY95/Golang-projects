// Note grpc.WithInsecure() is deprecated and should instead use grpc.WithTransportCredentials()
// with credentials.NewTLS(nil) to create a secure connection without any certificate verification.
// Reference: https://pkg.go.dev/google.golang.org/grpc#section-readme
package main

import (
	"context"
	"log"
	"time"

	pb "github.com/alaiy95/golang-projects/grpc-unary-service/usermgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Dial the gRPC server
	// 2nd and 3rd arguments are:
	// grpc.WithInsecure() option to create an insecure connection (i.e., without SSL/TLS)
	// and grpc.WithBlock() to block until the connection is established.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	c := pb.NewUserManagementClient(conn)

	// Create a context with a timeout of 1 second
	// This context will be used to make RPC calls to the server.
	// The defer statement at the end of the function ensures that the context is canceled once the function exits.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a map of new user details
	var new_users = make(map[string]int32)
	new_users["Alice"] = 43
	new_users["Bob"] = 30

	// Loop through the map and call the CreateNewUser RPC for each user
	for name, age := range new_users {
		// This function takes the context as its first argument,
		// followed by a pointer to a NewUser message containing the user's name and age.
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}

		// Print the details of the newly created user
		log.Printf(`User Details:
        NAME: %s
        AGE: %d
        ID: %d`, r.GetName(), r.GetAge(), r.GetId())
	}
}
