package product_type

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductType struct {
	ProductTypeID int    `json:"product_type_id"`
	ProductType   string `json:"product_type"`
}

// Get all product types
func GetProductTypes(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT Product_Type_ID, Product_Type FROM Product_Type")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var productTypes []ProductType
	for rows.Next() {
		var productType ProductType
		if err := rows.Scan(&productType.ProductTypeID, &productType.ProductType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		productTypes = append(productTypes, productType)
	}
	c.JSON(http.StatusOK, productTypes)
}

// Create a new product type
func CreateProductType(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	productTypeIDStr := c.Query("product_type_id")
	productType := c.Query("product_type")
	// Convert productTypeID to int (since it will be a string from the query)
	productTypeID, err := strconv.Atoi(productTypeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_type_id"})
		return
	}
	// Prepare the new product type struct
	newProductType := ProductType{
		ProductTypeID: productTypeID,
		ProductType:   productType,
	}
	// Execute the SQL statement with parameters
	result, err := db.Exec("INSERT INTO Product_Type (Product_Type_ID, Product_Type) VALUES (?, ?);", newProductType.ProductTypeID, newProductType.ProductType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Retrieve and set the last inserted ID if needed
	id, _ := result.LastInsertId()
	newProductType.ProductTypeID = int(id)
	// Send the response
	c.JSON(http.StatusCreated, newProductType)
}

// Update a product type
func UpdateProductType(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	productTypeIDStr := c.Query("product_type_id")
	productType := c.Query("product_type")
	// Convert productTypeID to an integer
	productTypeID, err := strconv.Atoi(productTypeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_type_id"})
		return
	}
	// Prepare the updated product type struct
	updatedProductType := ProductType{
		ProductTypeID: productTypeID,
		ProductType:   productType,
	}
	// Execute the SQL statement with parameters
	_, err = db.Exec("UPDATE Product_Type SET Product_Type = ? WHERE Product_Type_ID = ?", updatedProductType.ProductType, updatedProductType.ProductTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Send the response
	c.JSON(http.StatusOK, updatedProductType)
}

// Delete a product type
func DeleteProductType(c *gin.Context, db *sql.DB) {
	id := c.Param("product_type_id")
	_, err := db.Exec("DELETE FROM Product_Type WHERE Product_Type_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product Type deleted"})
}
