package v1

import (
	"database/sql"
	"fmt"
	"log"
)

type PostgresStruct struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
}

func (po *PostgresStruct) Init() {
	// Initialize PostgreSQL connection
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", po.PostgresHost, po.PostgresPort, po.PostgresUser, po.PostgresPassword, po.PostgresDb))
	if err != nil {
		log.Fatalf("Failed to open a connection to PostgreSQL: %v\n", err)
	}
	defer db.Close()
	log.Println("Successfully connected to PostgreSQL")
}

func InsertMessage(db *sql.DB, message string) error {
	// Insert message into PostgreSQL
	_, err := db.Exec("INSERT INTO messages (message) VALUES ($1)", message)
	if err != nil {
		return fmt.Errorf("Error inserting message into PostgreSQL: %v", err)
	}
	return nil
}

func GetMessages(db *sql.DB) ([]string, error) {
	// Retrieve messages from PostgreSQL
	rows, err := db.Query("SELECT message FROM messages")
	if err != nil {
		return nil, fmt.Errorf("Error retrieving messages from PostgreSQL: %v", err)
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var message string
		err := rows.Scan(&message)
		if err != nil {
			return nil, fmt.Errorf("Error scanning message: %v", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}
