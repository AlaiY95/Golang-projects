// This is a Go program that sets up a server with an API to perform CRUD operations on a database of books.
// The program uses the Fiber web framework and GORM as the database ORM. Let's go through the code to understand it better.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alaiy95/go-fiber-postgres/models"
	"github.com/alaiy95/go-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// The Book struct has three fields: Author, Title, and Publisher, all of which are of type string.
// Each field also has a corresponding JSON tag that specifies the name of the field when the struct is marshaled to or unmarshaled from JSON.
type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

// The Repository struct has a single field DB, which is a pointer to a gorm.DB object.
// The gorm.DB object is assumed to be a database connection object that provides a set of methods for querying and manipulating data in the database.
// Kind of like Django but with Go's own flavour :)
type Repository struct {
	DB *gorm.DB
}

// This function creates a new book record by parsing the JSON data in the HTTP request body.
// CreateBook is a struct method.  All our struct methods will have access to the repository, r.
func (r *Repository) CreateBook(context *fiber.Ctx) error {

	if context.Method() != fiber.MethodPost {
		// Return a "Method Not Allowed" error if the method is not POST
		return fiber.NewError(http.StatusMethodNotAllowed, "Computer says mmmmmmNO")
	}
	book := Book{}

	// It uses the BodyParser method of the fiber.Ctx struct to parse the request body and populate the book struct with the JSON data.
	// Fibre already has the ability to convert JSON data to a struct method
	err := context.BodyParser(&book)
	// If there is an error parsing the request body, it returns an HTTP status code of 422 (Unprocessable Entity) and an error message.
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	// Use the Create method of the gorm.DB struct to create a new database record with the book data
	err = r.DB.Create(&book).Error
	// If there is an error creating the record, it returns an HTTP status code of 400 (Bad Request) and an error message.
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	return nil
}

// DeleteBook removes a book from the database based on its ID
func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	// Extract the book ID from the URL parameters
	id := context.Params("id")

	// Create a new empty book model to store the results of the query
	bookModel := models.Books{}

	// Check if the ID is empty, and if so, send a JSON response with a status code of 500 to the client indicating an error
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	// Delete the book from the database using the ID
	err := r.DB.Delete(bookModel, id)

	// If the delete operation fails, return an error and send a JSON response with a status code of 400 to the client
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}

	// If the delete operation succeeds, send a JSON response with a status code of 200 to the client indicating success
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})

	// Return nil to indicate that there were no errors
	return nil
}

// It uses the Params method of the fiber.Ctx struct to get the value of the id parameter from the URL.
// GetBooks retrieves a list of books from the database and sends it as a JSON response to the client
func (r *Repository) GetBooks(context *fiber.Ctx) error {
	// Create a new empty slice of book models to store the results of the query
	bookModels := &[]models.Books{}

	// Retrieve the list of books from the database
	err := r.DB.Find(bookModels).Error
	if err != nil {
		// If the query fails, return an error and send a JSON response with a status code of 400 to the client
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	// If the query succeeds, send a JSON response with a status code of 200 to the client containing the list of book models in the "data" field
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})

	// Return nil to indicate that there were no errors
	return nil
}

// GetBookByID retrieves a single book from the database by its ID and sends it as a JSON response to the client
func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	// Extract the book ID from the URL parameters
	id := context.Params("id")

	// Create a new empty book model to store the results of the query
	bookModel := &models.Books{}

	// Check if the ID is empty, and if so, send a JSON response with a status code of 500 to the client indicating an error
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	// Print the ID to the console for debugging purposes
	fmt.Println("the ID is", id)

	// Retrieve the book from the database by its ID
	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		// If the query fails, return an error and send a JSON response with a status code of 400 to the client
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}

	// If the query succeeds, send a JSON response with a status code of 200 to the client containing the book model in the "data" field
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    bookModel,
	})

	// Return nil to indicate that there were no errors
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)

	app.Get("*", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Invalid request. Please check the URL and try again.")
	})
}

func main() {
	// use the godotenv package to load environment variables from a .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config) //create a new database connection

	if err != nil {
		log.Fatal("could not load the database")
	}
	// Applies any necessary database migrations. It calls the MigrateBooks function from the models package,
	// passing in the database connection as an argument
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	// creates a new instance of the Repository struct, passing in the database connection as an argument.
	r := Repository{
		DB: db,
	}
	// creates a new instance of the fiber.App struct and sets up the HTTP routes using the SetupRoutes method of the Repository struct.
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
