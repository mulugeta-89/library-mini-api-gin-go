package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
}

var books = []book {
	{"1", "war and peace", "Mulugeta Hailegnaw", 2},
	{"2", "The long journey", "Efriem shimels", 4},
	{"3", "The brief history of time", "Steven Hawking", 3},
}
func getBooks(c *gin.Context){
	c.IndentedJSON(http.StatusOK, books)
}
func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return 
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)

}
func getById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func getBookById(c *gin.Context){
	id := c.Param("id")
	book, err := getById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}
func checkoutBook(c *gin.Context){
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing query parameter"})
		return
	}
	book, err := getById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	}
	
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}

func checkinBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing query parameter"})
		return
	}

	book, err := getById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}
func main(){
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", getBookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/checkin", checkinBook)
	router.Run("localhost:8080")
}