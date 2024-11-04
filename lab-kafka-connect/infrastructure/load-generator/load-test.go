package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection settings
	host := "127.0.0.1"
	port := 5432
	user := "rafi"
	password := "123456"
	dbname := "defaultdb"
	sslmode := "disable"

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Connect to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}
	fmt.Println("Successfully connected to the database!")

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate random names and families and insert them into the table
	for i := 0; i < 1000000; i++ {
		name := randomString(5, 10)
		family := randomString(5, 12)

		query := `INSERT INTO employee (name, family) VALUES ($1, $2)`
		_, err := db.Exec(query, name, family)
		if err != nil {
			log.Fatal("Failed to insert data:", err)
		}

		if i%1000 == 0 {
			fmt.Printf("Inserted %d records\n", i)
		}
	}

	fmt.Println("Inserted all 10,000 records.")
}

// randomString generates a random string of a given length between min and max
func randomString(min, max int) string {
	length := rand.Intn(max-min+1) + min
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]rune, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
