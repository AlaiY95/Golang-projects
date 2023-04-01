package main

// This block of code imports the necessary packages for:
// the web server, the chi router, the renderer package, and the MongoDB driver.
import (
	// context: The context package provides a way to carry deadlines, cancellation signals,
	// and other request-scoped values across API boundaries and between processes.
	"context"

	// encoding/json: The encoding/json package provides functions for encoding and decoding JSON data in Go.
	"encoding/json"

	// log: The log package provides a simple logging package that allows you to write messages to the console or to a file.
	"log"

	// net/http: The net/http package provides a set of functions and types for working with HTTP requests and responses in Go.
	"net/http"

	// os: The os package provides a way to interact with the operating system, including accessing environment variables,
	// command-line arguments, and file I/O.
	"os"

	// os/signal: The os/signal package provides a way to receive signals from the operating system,
	// such as SIGINT or SIGTERM, and handle them in Go.
	"os/signal"

	// strings: The strings package provides a set of functions for working with strings in Go, such as splitting, joining, and searching
	"strings"

	//
	"time"

	// github.com/go-chi/chi: The go-chi/chi package is a lightweight and flexible router for building HTTP services in Go.
	"github.com/go-chi/chi"

	// github.com/go-chi/chi/middleware: The go-chi/chi/middleware package provides middleware functions for the chi router, such as logging and rate limiting.
	"github.com/go-chi/chi/middleware"

	// github.com/thedevsaddam/renderer: The thedevsaddam/renderer package is a Go package for rendering JSON and HTML templates.
	"github.com/thedevsaddam/renderer"

	// gopkg.in/mgo.v2: The mgo.v2 package is the official MongoDB driver for the Go programming language. It provides a way to interact with MongoDB databases using Go code.
	mgo "gopkg.in/mgo.v2"

	// gopkg.in/mgo.v2/bson: The mgo.v2/bson package provides functions for working with BSON (Binary JSON), which is the binary representation of JSON used by MongoDB.
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostName       string = "localhost:27017"
	dbName         string = "demo_todo"
	collectionName string = "todo"
	port           string = ":9000"
)

// The code also defines two structs: todoModel and todo.

type (
	// todoModel represents the data model used to store todo items in the MongoDB database,
	todoModel struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"createAt"`
	}

	// while todo represents the JSON format used to represent todo items in HTTP requests and responses.
	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}
)

// Initializes a session with a MongoDB database using the mgo package,
// connects to a MongoDB instance specified by hostName, sets the session to use monotonic consistency,
// and sets the db variable to the database specified by dbName.
func init() {
	rnd = renderer.New()
	sess, err := mgo.Dial(hostName)
	checkErr(err)
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}

// homeHandler: handles the home page request and renders a template called "static/home.tpl".
func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	checkErr(err)
}

// createTodo: handles POST requests to create a new todo item.
// It reads the request body and decodes it into a todo struct.
// It then performs some simple validation and saves the todo to the MongoDB database.

// createTodo is an HTTP handler function responsible for creating a new todo item.
// it takes 2 arguments:
// http.ResponseWriter is used to write the response that will be sent back to the client
// http.Request contains information about the incoming request.
func createTodo(w http.ResponseWriter, r *http.Request) {

	// This line declares a variable t of type todo. todo is a custom struct
	// that is defined elsewhere in the code. It is used to hold information
	// about the new todo item that is being created.
	var t todo

	// 1st Stage
	// This code block decodes the JSON data from the request body into the todo struct.
	// If there is an error decoding the JSON data, it returns an HTTP response with a
	// 422 status code (Unprocessable Entity) and an error message.

	// json.NewDecoder(r.Body) creates a new json.Decoder that reads from r.Body, which is the request body,
	// and .Decode(&t) decodes the JSON data from the request body and stores it in the todo struct.
	// &t is a pointer to the t variable, which allows the Decode method to modify the t variable.
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		// // If there is an error decoding the request body, return a 400 Bad Request response.
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	// 2nd Stage
	// This code block checks if the Title field of the todo struct is empty.
	// If it is empty, it returns an HTTP response with a 400 status code (Bad Request)
	// and an error message indicating that the Title field is required.
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is requried",
		})
		return
	}

	// 3rd Stage
	// This code block creates a new todoModel struct with the ID, Title, Completed, and CreatedAt fields.
	// The ID field is generated using bson.NewObjectId(), which creates a new ObjectId that can be used
	// as a unique identifier for the new todo item. The Title field is set to the Title field of the todo struct
	// that was decoded earlier. The Completed field is set to false by default, since a new todo item
	// is not completed when it is first created. The CreatedAt field is set to the current time using time.Now().
	tm := todoModel{
		ID:        bson.NewObjectId(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	// 4th Stage
	// This code block inserts the new todoModel struct into a database collection.
	// The db.C(collectionName) expression returns a collection object from the database that corresponds to
	// the collection name specified by collectionName. The Insert method in MongoDB is called on the collection object with
	// a pointer to the todoModel struct as an argument. If there is an error inserting the new todo item into the database,
	// it returns an HTTP response with a 422 status code (Unprocessable Entity) and an error message.
	if err := db.C(collectionName).Insert(&tm); err != nil { // err != means if error is present
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save todo",
			"error":   err,
		})
		return
	}

	// 5th Stage
	// Return a 201 Created response with the created todo item in the response body.
	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"todo_id": tm.ID.Hex(),
	})
}

// updateTodo: handles PUT requests to update an existing todo item.
// It reads the request body and decodes it into a todo struct.
// It then performs some simple validation and updates the corresponding todo in the MongoDB database.

// This line defines a function called updateTodo that takes two arguments: an http.ResponseWriter and an http.Request.
// These are standard parameters for an HTTP handler function in Go.
func updateTodo(w http.ResponseWriter, r *http.Request) {

	// This line extracts the value of the id parameter from the URL using chi.URLParam().
	// It then trims any leading or trailing whitespace from the parameter value using strings.TrimSpace()
	// and assigns the result to the id variable.
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	// This code block checks whether the id variable contains a valid MongoDB ObjectID.
	// It uses the bson.IsObjectIdHex() function from the MongoDB driver to perform the check.
	// If the ID is invalid, it returns a JSON response with a status code of http.StatusBadRequest and an error message.
	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	// This code block initializes an empty todo struct and decodes the request body into it using json.NewDecoder() and Decode().
	// If there is an error during decoding, it returns a JSON response with a status code of http.StatusProcessing and the error message.
	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	// This code block checks that the Title field of the todo struct is not empty.
	// If it is, it returns a JSON response with a status code of http.StatusBadRequest and an error message.
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is requried",
		})
		return
	}

	// This code block updates the todo item in the MongoDB collection. It uses the
	// db.C(collectionName).Update() function to perform the update, passing in a MongoDB query that specifies
	// which document to update and a set of updates to apply. If there is an error during the update operation,
	// it returns a JSON response with a status code of http.StatusProcessing and the error message.
	if err := db.C(collectionName).
		Update(
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"title": t.Title, "completed": t.Completed},
		); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to update todo",
			"error":   err,
		})
		return
	}

	// This code block returns a JSON response with a status code of http.StatusOK and a success message if the update operation succeeds.
	// It uses the rnd.JSON() function to serialize the response data to JSON and write it to the response writer.
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo updated successfully",
	})
}

// fetchTodos: handles GET requests to fetch a list of all the todo items.
// It queries the MongoDB database and returns the result in JSON format.

func fetchTodos(w http.ResponseWriter, r *http.Request) {

	// This line declares an empty slice of todoModel structs to hold the results of the MongoDB query.
	// this will be json that we want to send to the frontend.
	todos := []todoModel{}

	// This code block fetches all documents from the MongoDB collection specified in the collectionName variable.
	// It uses the db.C(collectionName).Find() function to perform the query and
	// the All() method to retrieve all matching documents and decode them into the todos slice.
	// If there is an error during the query, it returns a JSON response with a status code of http.StatusProcessing and the error message.
	if err := db.C(collectionName).
		Find(bson.M{}).
		All(&todos); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
		})
		return
	}

	// This block of code is creating a new slice of todo structs named todoList, and then populating it with new todo struct instances
	// based on the todoModel structs that are passed in through the todos slice. See For loop below.
	todoList := []todo{}

	// Each todoModel struct in todos is looped over with a range loop, and for each one, a new todo struct is
	// created with the ID, Title, Completed, and CreatedAt fields set to the corresponding values of the todoModel struct.
	//The new todo struct instances are then appended to the todoList slice using the append() function.
	// After this loop completes, todoList contains a list of todo structs that have been populated with the relevant data from todos.

	// We're using the for
	for _, t := range todos {
		todoList = append(todoList, todo{ // range over todos which is the BSON data received from database and append 1 by 1 to the todolist
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	// This code block returns a JSON response with a status code of http.StatusOK and a data field to the FRONTEND that contains
	// the todoList slice serialized as JSON. It uses the rnd.JSON() function to serialize the response data to JSON and write it to the response writer.
	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todoList,
	})
}

// deleteTodo: handles DELETE requests to delete an existing todo item.
// It reads the ID of the todo from the request URL, performs some validation,
// and deletes the corresponding todo from the MongoDB database.

// This line defines a function called deleteTodo that takes two arguments: an http.ResponseWriter and an http.Request.
// These are standard parameters for an HTTP handler function in Go.
func deleteTodo(w http.ResponseWriter, r *http.Request) {

	// This line extracts the id parameter from the URL using the chi.URLParam() function and stores it as a string.
	// It also trims any leading or trailing whitespace from the id string using the strings.TrimSpace() function.
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	// This code block checks whether the id parameter is a valid MongoDB object ID by calling the bson.IsObjectIdHex() function.
	// If the id is invalid, it returns a JSON response with a status code of http.StatusBadRequest and an error message.

	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	// This code block removes a document from the MongoDB collection specified in the collectionName variable using the RemoveId() method.
	// It passes the bson.ObjectIdHex(id) value to the RemoveId() method to identify the document to remove.
	// If there is an error during the delete operation, it returns a JSON response with a status code of http.StatusProcessing and an error message.
	if err := db.C(collectionName).RemoveId(bson.ObjectIdHex(id)); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
		return
	}

	// This code block returns a JSON response with a status code of http.StatusOK and a success message.
	// It uses the rnd.JSON() function to serialize the response data to JSON and write it to the response writer.
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo deleted successfully",
	})
}

// Finally, the main function sets up the HTTP server and starts listening on the specified port.
// It uses the stopChan channel to listen for a signal to stop the server gracefully when the user presses Ctrl+C.

func main() {
	// create a channel to receive OS interrupt signal
	stopChan := make(chan os.Signal)

	// signal.Notify(stopChan, os.Interrupt) notifies the channel stopChan when an OS interrupt signal is received,
	// such as when the user presses Ctrl+C to terminate the program.
	signal.Notify(stopChan, os.Interrupt)

	// create a new router using Chi
	r := chi.NewRouter()

	// use middleware to log each incoming HTTP request and its corresponding response to the console.
	r.Use(middleware.Logger)

	// define a route for the home page when a GET request is made to the root path.
	r.Get("/", homeHandler)

	// mount the todoHandlers router at the "/todo" path
	r.Mount("/todo", todoHandlers())

	// create an http.Server instance with the specified configurations
	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	// The server will run indefinitely until it receives an interrupt signal.
	go func() {
		log.Println("Listening on port ", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Waits for an OS interrupt signal to be received, which will unblock the main thread.
	<-stopChan

	// gracefully shut down the server
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Server gracefully stopped!")
}

// This is Go code that defines an HTTP router for handling
// CRUD (Create, Read, Update, Delete) operations for a "todo" resource.
func todoHandlers() http.Handler {
	// The code uses the chi router package to define routes for the different CRUD operations.
	//  The todoHandlers function returns a http.Handler object that can be used to handle incoming HTTP requests.
	rg := chi.NewRouter()

	// The r.Group method is used to group the routes under a common prefix,
	// in this case, the root path for the "todo" resource. This allows the router to apply middleware
	// or other routing rules to all the routes in the group.
	rg.Group(func(r chi.Router) {

		// The routes are defined using r.Get, r.Post, r.Put, and r.Delete methods.
		// These methods map the HTTP methods to the corresponding CRUD operations.
		// For example, r.Get("/", fetchTodos) maps a GET request to the / URL path to the
		// fetchTodos function, which retrieves a list of todos.
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)

		// The routes use path parameters to identify individual todos.
		// For example, r.Put("/{id}", updateTodo) maps a PUT request to a URL path that includes a todo ID.
		// The ID is passed as a parameter to the updateTodo function, which updates the corresponding todo.
	})
	return rg
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err) //respond with error page or message
	}
}
