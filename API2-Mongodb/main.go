// The code is a basic implementation of a RESTful API in Go, using the httprouter package for routing and the mgo package for connecting to a MongoDB database. Here's a breakdown of the code:

// The code imports the necessary packages: httprouter, mgo, and net/http. The controllers package is also imported, which will contain the handlers for the various HTTP routes.

// The main function sets up a new httprouter instance and creates a new UserController, passing in the result of the getSession function. The UserController is used to handle requests to the various routes.

// The main function then sets up the routes using the httprouter instance's GET, POST, and DELETE methods. Each route corresponds to a method on the UserController.

// Finally, the main function starts an HTTP server listening on port 9000 of the local machine.

// The getSession function returns a new mgo session connected to a MongoDB database running on the default port (27017) of the local machine. If there is an error connecting to the database, the function panics.

// Overall, this code demonstrates a simple implementation of a RESTful API in Go, with MongoDB as the underlying database. The code is relatively concise and easy to follow, although it could be improved with additional error handling and more robust database connections.

package main

import (
	"net/http" // The http library provides a set of functions for working with HTTP requests and responses.

	"github.com/julienschmidt/httprouter" // The httprouter library provides a fast HTTP router implementation for Go
	"gopkg.in/mgo.v2"                     // The mgo library is a MongoDB driver for Go

	// The controllers package contains the user controller logic that will handle incoming HTTP requests.
	"github.com/akhil/mongo-golang/controllers"
)

// The main function starts by creating a new httprouter instance and a new user controller instance.
func main() {

	r := httprouter.New() // new httprouter instance r
	// new user controller instance uc using the getSession() function
	// which returns a new mgo.Session object that connects to the MongoDB instance running on localhost:27017.
	uc := controllers.NewUserController(getSession())
	// setting up several HTTP routes
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	// Start the HTTP server by calling http.ListenAndServe().
	// Function takes two arguments:
	// The first is the address to listen on, and the second is the httprouter instance to use for handling incoming HTTP requests.
	http.ListenAndServe("localhost:9000", r)
}

// getSession() function creates a new mgo.Session object that connects to the MongoDB instance running on localhost:27017
func getSession() *mgo.Session {

	s, err := mgo.Dial("mongodb://localhost:27017")
	// If an error occurs while connecting to the database, the function panics and terminates the program.
	if err != nil {
		panic(err)
	}
	// The function returns the mgo.Session object, which is then used to create a new user controller instance.
	return s
}
