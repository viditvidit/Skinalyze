package brand

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Brand struct {
	BrandID int    `json:"brand_id"`
	Brand   string `json:"brand"`
}

// Get all brands
func GetBrands(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT Brand_ID, Brand FROM BRAND")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var brands []Brand
	for rows.Next() {
		var brand Brand
		if err := rows.Scan(&brand.BrandID, &brand.Brand); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		brands = append(brands, brand)
	}
	c.JSON(http.StatusOK, brands)
}

// Create a new brand
func CreateBrand(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	brandIDStr := c.Query("brand_id")
	brand := c.Query("brand")
	// Convert brandID to int (since it will be a string from the query)
	brandID, err := strconv.Atoi(brandIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand_id"})
		return
	}
	// Prepare the new brand struct
	newBrand := Brand{
		BrandID: brandID,
		Brand:   brand,
	}
	// Execute the SQL statement with parameters
	result, err := db.Exec("INSERT INTO BRAND (Brand_ID, Brand) VALUES (?, ?);", newBrand.BrandID, newBrand.Brand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Retrieve and set the last inserted ID if needed
	id, _ := result.LastInsertId()
	newBrand.BrandID = int(id)
	// Send the response
	c.JSON(http.StatusCreated, newBrand)
}

// Update an existing brand
func UpdateBrand(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	brandIDStr := c.Query("brand_id")
	brand := c.Query("brand")

	// Convert brandID to an integer
	brandID, err := strconv.Atoi(brandIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand_id"})
		return
	}
	// Prepare the updated brand struct
	updatedBrand := Brand{
		BrandID: brandID,
		Brand:   brand,
	}
	// Execute the SQL statement with parameters
	_, err = db.Exec("UPDATE BRAND SET Brand = ? WHERE Brand_ID = ?", updatedBrand.Brand, updatedBrand.BrandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Send the response
	c.JSON(http.StatusOK, updatedBrand)
}

// Delete a brand
func DeleteBrand(c *gin.Context, db *sql.DB) {
	id := c.Param("brand_id")
	_, err := db.Exec("DELETE FROM BRAND WHERE Brand_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Brand deleted"})
}
