package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Read database connection details from environment variables set by Railway
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	
	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify the connection to the database
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	// Set up the Gin router
	r := gin.Default()

	// Define a handler to get all students from the database
	r.GET("/students", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name FROM student")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
			return
		}
		defer rows.Close()

		var students []Student
		for rows.Next() {
			var student Student
			if err := rows.Scan(&student.ID, &student.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
				return
			}
			students = append(students, student)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row iteration error"})
			return
		}

		c.JSON(http.StatusOK, students)
	})

	// Start the server
	r.Run() // listen and serve on 0.0.0.0:8080
}
