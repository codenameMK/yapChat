package v1

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Import the pq driver
)

type PostgresStruct struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
}

func (po *PostgresStruct) Init() (*sql.DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", po.PostgresUser, po.PostgresPassword, po.PostgresHost, po.PostgresPort, po.PostgresDb)
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to open a connection to PostgreSQL: %v", err)
	}

	// Optionally, you can ping the database to verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping the PostgreSQL database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return db, nil
}

// ExecuteQuery executes the passed query with arguments and returns the results or an error
func ExecuteQuery(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	// Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		// For non-SELECT queries like INSERT/UPDATE, return exec result instead
		result, errExec := db.Exec(query, args...)
		if errExec != nil {
			return nil, fmt.Errorf("Error executing query: %v", errExec)
		}
		// For non-SELECT queries, return a success message (empty result set)
		rowsAffected, _ := result.RowsAffected()
		return []map[string]interface{}{{"rowsAffected": rowsAffected}}, nil
	}
	defer rows.Close()

	// Parse the result for SELECT queries
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Error fetching columns: %v", err)
	}

	// Prepare the result set to hold rows
	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of interface{}'s to hold each column's value
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the result into the pointers
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}

		// Create a map to hold column:value
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]

			// Handle NULL values
			switch val.(type) {
			case []byte:
				v = string(val.([]byte))
			default:
				v = val
			}
			rowMap[col] = v
		}
		results = append(results, rowMap)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("Error iterating rows: %v", rows.Err())
	}

	return results, nil
}

func InsertMessage(db *sql.DB, conversationId int32, timeStamp time.Time, message string, userId int32, receiverId int32) error {
	// Insert message into PostgreSQL
	query := `
        INSERT INTO Messages (conversation_id, sender_id, receiver_id, timestamp, content, status)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	
	_, err := db.Exec(query, conversationId, userId, receiverId, timeStamp, message, "Unread")
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
