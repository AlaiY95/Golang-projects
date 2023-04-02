// This program provides a simple API for a book library that allows users to get a list of books,
// get a specific book by ID, check out a book, return a book, and create a new book.

package main

// First, the necessary packages are imported: "net/http" for HTTP protocol, "errors" for handling errors, and "github.com/gin-gonic/gin" for the Gin web framework.
import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

// Next, a struct book is defined to represent a book with its properties like ID, Title, Author, and Quantity.
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// This is a slice of book structs that serves as a simple database of books for this book management system.
var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

// Functions are defined to perform CRUD operations on the books.

// The getBooks function handles requests to the /books endpoint and returns a list of all books in the library as JSON.
// getBooks retrieves all the books from the slice of books and returns them in JSON format.
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

// bookById functionHandles requests to '/books/id/ and retrieves a single book by its ID.
// It calls the function to retrieve the book with the specified ID,
// and returns the book or an error message in JSON format.
func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	// If the book is not found, it returns a 404 Not Found status code.
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// The checkoutBook function is a handler for the HTTP PATCH method on the "/checkout" endpoint.
func checkoutBook(c *gin.Context) {
	// Extract the "id" query parameter from the HTTP request.
	id, ok := c.GetQuery("id")

	// If the "id" query parameter is not present, return a HTTP response with a 400 Bad Request status code ana JSON Object with the error.
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}
	// Get the book from the "books" slice with the given "id" using the `getBookById` function.
	book, err := getBookById(id)

	// If the book is not found in the "books" slice, return a HTTP response with a 404 Not Found status code.
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	// If the book is found but its quantity is already zero, return a HTTP response with a 400 Bad Request status code.
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	// Decrement the quantity of the book by one and return a HTTP response with a 200 OK status code and the book object as a JSON response.
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

// The returnBook function handles requests to the /return endpoint and allows a user to return a book to the library.
// returnBook function is used to handle the PATCH request to return a book. It takes a gin.Context object as its only parameter
func returnBook(c *gin.Context) {

	// check if the id query parameter is present in the request URL by calling the GetQuery method of the gin.Context object.
	id, ok := c.GetQuery("id")

	// If the id parameter is not present, the function returns a 400 Bad Request status code with a JSON message indicating that the id parameter is missing.
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	// If the id parameter is present, the function calls the getBookById function to retrieve the book with the specified ID.
	book, err := getBookById(id)

	// If the book is not found, the function returns a 404 Not Found status code with a JSON message indicating that the book was not found.
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	// If the book is found, the function increments the Quantity field of the book object by 1
	book.Quantity += 1

	// and returns a 200 OK status code with the updated book object in the response body as a JSON object.
	c.IndentedJSON(http.StatusOK, book)
}

// getBookById is a helper function that takes an ID string and returns a pointer to a book struct and an error.
// It loops through the global books slice, checks if the book with the given ID exists, and returns a pointer to that book if it does.
// If the book is not found, it returns a nil pointer and an error message.
func getBookById(id string) (*book, error) {
	// Loop through the global books slice
	for i, b := range books {
		// Check if the ID of the current book in the loop matches the given ID
		if b.ID == id {
			// Return a pointer to the current book in the loop and no error
			return &books[i], nil
		}
	}
	// If the function reaches this point, the book was not found, so return a nil pointer and an error message
	return nil, errors.New("book not found")
}

// The createBook function handles requests to the /books endpoint with a POST request method
// and allows a user to create a new book by sending a JSON object in the request body. The function parses the
func createBook(c *gin.Context) {
	var newBook book

	// Bind the request body JSON data to the newBook variable
	// If an error occurs while binding, return early from the function
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	// Append the new book to the books slice
	books = append(books, newBook)

	// Respond with a JSON representation of the new book and a 201 Created status code
	c.IndentedJSON(http.StatusCreated, newBook)
}

// The router object sets up the routing for the API by defining the endpoints for each of the above functions and starting the server on port 8080.

// A PATCH request is an HTTP method that is used to partially update a resource on the server.
// It is similar to the HTTP PUT method, which is used to completely replace a resource on the server.
// However, the PATCH method allows for more fine-grained updates of a resource by specifying only the fields that need to be updated,
// instead of sending the entire resource with all its fields.
// The PATCH request can be used in situations where updating a whole resource is unnecessary or undesirable.
func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
