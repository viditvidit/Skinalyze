package concern

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Concern struct {
	ConcernID int    `json:"concern_id"`
	Concern   string `json:"concern"`
}

// Get all concerns
func GetConcerns(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT Concern_ID, Concern FROM Concern")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var concerns []Concern
	for rows.Next() {
		var concern Concern
		if err := rows.Scan(&concern.ConcernID, &concern.Concern); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		concerns = append(concerns, concern)
	}
	c.JSON(http.StatusOK, concerns)
}

// Create a new concern
func CreateConcern(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	concernIDStr := c.Query("concern_id")
	concern := c.Query("concern")
	// Convert concernID to int (since it will be a string from the query)
	concernID, err := strconv.Atoi(concernIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid concern_id"})
		return
	}
	// Prepare the new concern struct
	newConcern := Concern{
		ConcernID: concernID,
		Concern:   concern,
	}
	// Execute the SQL statement with parameters
	result, err := db.Exec("INSERT INTO Concern (Concern_ID, Concern) VALUES (?, ?);", newConcern.ConcernID, newConcern.Concern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Retrieve and set the last inserted ID if needed
	id, _ := result.LastInsertId()
	newConcern.ConcernID = int(id)
	// Send the response
	c.JSON(http.StatusCreated, newConcern)
}

// Update a concern
func UpdateConcern(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	concernIDStr := c.Query("concern_id")
	concern := c.Query("concern")
	// Convert concernID to an integer
	concernID, err := strconv.Atoi(concernIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid concern_id"})
		return
	}
	// Prepare the updated concern struct
	updatedConcern := Concern{
		ConcernID: concernID,
		Concern:   concern,
	}
	// Execute the SQL statement with parameters
	_, err = db.Exec("UPDATE Concern SET Concern = ? WHERE Concern_ID = ?", updatedConcern.Concern, updatedConcern.ConcernID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Send the response
	c.JSON(http.StatusOK, updatedConcern)
}

// Delete a concern
func DeleteConcern(c *gin.Context, db *sql.DB) {
	id := c.Param("concern_id")
	_, err := db.Exec("DELETE FROM Concern WHERE Concern_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Concern deleted"})
}
