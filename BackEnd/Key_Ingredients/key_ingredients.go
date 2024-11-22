package key_ingredients

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type KeyIngredients struct {
	KeyIngredientsID int    `json:"key_ingredients_id"`
	KeyIngredient    string `json:"ingredient"`
}

// Get all key ingredients
func GetKeyIngredients(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT Key_Ingredients_ID, Key_Ingredients FROM Key_Ingredients")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var keyIngredients []KeyIngredients
	for rows.Next() {
		var ingredient KeyIngredients
		if err := rows.Scan(&ingredient.KeyIngredientsID, &ingredient.KeyIngredient); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		keyIngredients = append(keyIngredients, ingredient)
	}

	c.JSON(http.StatusOK, keyIngredients)
}

// Create a new key ingredient
func CreateKeyIngredient(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	keyIngredientsIDStr := c.Query("key_ingredients_id")
	keyingredient := c.Query("key_ingredients")
	// Convert keyIngredientsID to int (since it will be a string from the query)
	keyIngredientsID, err := strconv.Atoi(keyIngredientsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key_ingredients_id"})
		return
	}
	// Prepare the new key ingredient struct
	newKeyIngredient := KeyIngredients{
		KeyIngredientsID: keyIngredientsID,
		KeyIngredient:    keyingredient,
	}
	// Execute the SQL statement with parameters
	result, err := db.Exec("INSERT INTO Key_Ingredients (Key_Ingredients_ID, Key_Ingredients) VALUES (?, ?);", newKeyIngredient.KeyIngredientsID, newKeyIngredient.KeyIngredient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Retrieve and set the last inserted ID if needed
	id, _ := result.LastInsertId()
	newKeyIngredient.KeyIngredientsID = int(id)
	// Send the response
	c.JSON(http.StatusCreated, newKeyIngredient)
}

// Update a key ingredient
func UpdateKeyIngredient(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	keyingredientidStr := c.Query("key_ingredients_id")
	keyingredient := c.Query("key_ingredients")

	// Convert id to an integer
	id, err := strconv.Atoi(keyingredientidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	// Prepare the updated key ingredient struct
	updatedKeyIngredient := KeyIngredients{
		KeyIngredientsID: id,
		KeyIngredient:    keyingredient,
	}
	// Execute the SQL statement with parameters
	_, err = db.Exec("UPDATE Key_Ingredients SET Key_Ingredients = ? WHERE Key_Ingredients_ID = ?", updatedKeyIngredient.KeyIngredient, updatedKeyIngredient.KeyIngredientsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Send the response
	c.JSON(http.StatusOK, updatedKeyIngredient)
}

// Delete a key ingredient
func DeleteKeyIngredient(c *gin.Context, db *sql.DB) {
	id := c.Param("key_ingredients_id")
	_, err := db.Exec("DELETE FROM Key_Ingredients WHERE Key_Ingredients_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Key Ingredient deleted"})
}
