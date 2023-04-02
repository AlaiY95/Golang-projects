package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alaiy95/golang-projects/api2-mongodb/models"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Define a UserController struct to hold a reference to a MongoDB client
type UserController struct {
	client *mongo.Client
}

// Define the NewUserController function to create a new UserController instance
func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Extract the ID parameter from the request parameters
	id := p.ByName("id")

	// Check if the ID parameter is a valid hexadecimal string
	if !primitive.IsValidHexID(id) {
		log.Printf("Invalid ObjectID provided: %v", id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Convert the ID string to a MongoDB ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error parsing ObjectID from request parameter: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("Converted ID parameter to ObjectID: %v\n", oid)

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve a user with that ID from the "users" collection in the "mongo-golang" database
	u := models.User{}
	filter := bson.M{"_id": oid}
	err = uc.client.Database("mongo-golang").Collection("users").FindOne(ctx, filter).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Error retrieving user: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("Error retrieving user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the retrieved user object into a JSON-encoded byte slice
	uj, err := json.Marshal(u)
	if err != nil {
		log.Printf("Error encoding user as JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON-encoded user object to the response body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

// The CreateUser method creates a new user in the database
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Create a new User object and decode the JSON-encoded user data from the request body
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %s", err.Error())
		return
	}

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate a new MongoDB ObjectID for the user and set the user's ID field
	log.Println("Before generating ObjectID")
	u.Id = primitive.NewObjectID()
	log.Printf("Generated ObjectID: %s\n", u.Id.Hex())

	// Insert the user into the "users" collection in the "mongo-golang" database
	_, err1 := uc.client.Database("mongo-golang").Collection("users").InsertOne(ctx, u)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the inserted user object into a JSON-encoded byte slice
	uj, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling user: %s", err.Error())
		return
	}

	// Set the response headers and write the JSON-encoded user object to the response body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

// The DeleteUser method deletes a user from the database by ID
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := p.ByName("id")
	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid user ID")
		return
	}

	collection := uc.client.Database("mongo-golang").Collection("users")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error deleting user")
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "User not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User deleted successfully")
}
