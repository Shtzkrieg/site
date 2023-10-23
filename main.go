package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func uploadHandler(c *gin.Context) {
	// Single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open the file from the uploaded file header
	uploadedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer uploadedFile.Close()

	// Create a new file on the server to store the uploaded file
	newFile, err := os.Create("uploads/" + file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer newFile.Close()

	// Copy the uploaded file's content to the new file
	_, err = io.Copy(newFile, uploadedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully!"})
}

func main() {
	r := gin.Default()

	// Static html server
	r.Static("/static", "./static")

	// Access the file dumper
	r.GET("/dumper", func(c *gin.Context) {
		c.File("./static/dumper.html")
	})

	// File upload endpoint
	r.POST("/upload", uploadHandler)

	// Start the Gin server
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
