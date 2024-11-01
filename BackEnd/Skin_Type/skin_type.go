package skin_type

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type SkinType struct {
	SkinTypeID int    `json:"skin_type_id"`
	SkinType   string `json:"skin_type"`
}

// Get all skin types
func GetSkinTypes(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT Skin_Type_ID, Skin_Type FROM SKIN_TYPE")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var skinTypes []SkinType
	for rows.Next() {
		var skinType SkinType
		if err := rows.Scan(&skinType.SkinTypeID, &skinType.SkinType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		skinTypes = append(skinTypes, skinType)
	}
	c.JSON(http.StatusOK, skinTypes)
}

// Create a new skin type
func CreateSkinType(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	skinTypeIDStr := c.Query("skin_type_id")
	skinType := c.Query("skin_type")
	// Convert skinTypeID to int (since it will be a string from the query)
	skinTypeID, err := strconv.Atoi(skinTypeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skin_type_id"})
		return
	}
	// Prepare the new skin type struct
	newSkinType := SkinType{
		SkinTypeID: skinTypeID,
		SkinType:   skinType,
	}
	// Execute the SQL statement with parameters
	result, err := db.Exec("INSERT INTO SKIN_TYPE (Skin_Type_ID, Skin_Type) VALUES (?, ?);", newSkinType.SkinTypeID, newSkinType.SkinType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Retrieve and set the last inserted ID if needed
	id, _ := result.LastInsertId()
	newSkinType.SkinTypeID = int(id)
	// Send the response
	c.JSON(http.StatusCreated, newSkinType)
}

// Update a skin type
func UpdateSkinType(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	skinTypeIDStr := c.Query("skin_type_id")
	skinType := c.Query("skin_type")
	// Convert skinTypeID to an integer
	skinTypeID, err := strconv.Atoi(skinTypeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skin_type_id"})
		return
	}
	// Prepare the updated skin type struct
	updatedSkinType := SkinType{
		SkinTypeID: skinTypeID,
		SkinType:   skinType,
	}
	// Execute the SQL statement with parameters
	_, err = db.Exec("UPDATE SKIN_TYPE SET Skin_Type = ? WHERE Skin_Type_ID = ?", updatedSkinType.SkinType, updatedSkinType.SkinTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Send the response
	c.JSON(http.StatusOK, updatedSkinType)
}

// Delete a skin type
func DeleteSkinType(c *gin.Context, db *sql.DB) {
	id := c.Param("skin_type_id")
	_, err := db.Exec("DELETE FROM SKIN_TYPE WHERE Skin_Type_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Skin Type deleted"})
}
