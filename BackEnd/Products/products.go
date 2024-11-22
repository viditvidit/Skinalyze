package products

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Product struct {
	ProductID        int    `json:"product_id"`
	ProductName      string `json:"product_name"`
	AllIngredients   string `json:"all_ingredients"`
	ConcernID        int    `json:"concern_id"`
	Concern          string `json:"concern"`
	SkinTypeID       int    `json:"skin_type_id"`
	SkinType         string `json:"skin_type"`
	BrandID          int    `json:"brand_id"`
	Brand            string `json:"brand"`
	ProductTypeID    int    `json:"product_type_id"`
	KeyIngredientsID int    `json:"key_ingredients_id"`
	KeyIngredients   string `json:"key_ingredients"`
}

// Get all products
func GetProducts(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
	SELECT Product_ID, Product_Name, All_Ingredients, Concern_ID, Skin_Type_ID,
	Brand_ID, Product_Type_ID, Key_Ingredients_ID
	FROM Products`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ProductID, &product.ProductName, &product.AllIngredients, &product.ConcernID,
			&product.SkinTypeID, &product.BrandID, &product.ProductTypeID, &product.KeyIngredientsID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// Get Select Products
func GetSelectProducts(c *gin.Context, db *sql.DB, concernID, skinTypeID int) {
	rows, err := db.Query(`
    SELECT 
        p.Product_Name,
        p.All_Ingredients,
        b.Brand,
        c.Concern,
        k.Key_Ingredients,
        s.Skin_Type
    FROM PRODUCTS p
    INNER JOIN brand b ON p.Brand_ID = b.Brand_ID
    INNER JOIN concern c ON p.Concern_ID = c.Concern_ID
    INNER JOIN key_ingredients k ON p.Key_Ingredients_ID = k.Key_Ingredients_ID
    INNER JOIN skin_type s ON p.Skin_Type_ID = s.Skin_Type_ID
    WHERE p.Concern_ID = ? 
    AND p.Skin_Type_ID = ?
    ORDER BY p.PRODUCT_TYPE_ID
    `, concernID, skinTypeID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ProductName,
			&product.AllIngredients,
			&product.Brand,
			&product.Concern,
			&product.KeyIngredients,
			&product.SkinType,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}
	c.JSON(http.StatusOK, products)
}

// Get Select Products of specific type
func GetSelectProductsByType(c *gin.Context, db *sql.DB, concernID, skinTypeID, productTypeID int) {
	rows, err := db.Query(`
    SELECT 
        p.Product_Name,
        p.All_Ingredients,
        b.Brand,
        c.Concern,
        k.Key_Ingredients,
        s.Skin_Type
    FROM PRODUCTS p
    INNER JOIN brand b ON p.Brand_ID = b.Brand_ID
    INNER JOIN concern c ON p.Concern_ID = c.Concern_ID
    INNER JOIN key_ingredients k ON p.Key_Ingredients_ID = k.Key_Ingredients_ID
    INNER JOIN skin_type s ON p.Skin_Type_ID = s.Skin_Type_ID
    WHERE p.Concern_ID = ? 
    AND p.Skin_Type_ID = ?
    AND p.Product_Type_ID = ?
    ORDER BY p.PRODUCT_TYPE_ID
    `, concernID, skinTypeID, productTypeID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ProductName,
			&product.AllIngredients,
			&product.Brand,
			&product.Concern,
			&product.KeyIngredients,
			&product.SkinType,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}
	c.JSON(http.StatusOK, products)
}

// Create a new product
func CreateProduct(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	productIDStr := c.Query("product_id")
	productName := c.Query("product_name")
	allIngredients := c.Query("all_ingredients")
	concernIDStr := c.Query("concern_id")
	skinTypeIDStr := c.Query("skin_type_id")
	brandIDStr := c.Query("brand_id")
	productTypeIDStr := c.Query("product_type_id")
	keyIngredientsIDStr := c.Query("key_ingredients_id")
	// Convert string IDs to integers
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
		return
	}
	concernID, err := strconv.Atoi(concernIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid concern_id"})
		return
	}
	skinTypeID, err := strconv.Atoi(skinTypeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skin_type_id"})
		return
	}
	brandID, err := strconv.Atoi(brandIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand_id"})
		return
	}
	productTypeID, err := strconv.Atoi(productTypeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_type_id"})
		return
	}
	keyIngredientsID, err := strconv.Atoi(keyIngredientsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key_ingredients_id"})
		return
	}
	// Prepare the new product struct
	newProduct := Product{
		ProductID:        productID,
		ProductName:      productName,
		AllIngredients:   allIngredients,
		ConcernID:        concernID,
		SkinTypeID:       skinTypeID,
		BrandID:          brandID,
		ProductTypeID:    productTypeID,
		KeyIngredientsID: keyIngredientsID,
	}
	// Execute the SQL statement with parameters
	result, err := db.Exec(`
    INSERT INTO PRODUCTS (Product_ID, Product_Name, All_Ingredients, Concern_ID, Skin_Type_ID, Brand_ID, Product_Type_ID, Key_Ingredients_ID)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		newProduct.ProductID,
		newProduct.ProductName,
		newProduct.AllIngredients,
		newProduct.ConcernID,
		newProduct.SkinTypeID,
		newProduct.BrandID,
		newProduct.ProductTypeID,
		newProduct.KeyIngredientsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Retrieve and set the last inserted ID if needed
	id, _ := result.LastInsertId()
	newProduct.ProductID = int(id)
	// Send the response
	c.JSON(http.StatusCreated, newProduct)
}

// Update a product
func UpdateProduct(c *gin.Context, db *sql.DB) {
	// Read the query parameters
	productIDStr := c.Query("product_id")
	productName := c.Query("product_name")
	allIngredients := c.Query("all_ingredients")
	concernIDStr := c.Query("concern_id")
	skinTypeIDStr := c.Query("skin_type_id")
	brandIDStr := c.Query("brand_id")
	productTypeIDStr := c.Query("product_type_id")
	keyIngredientsIDStr := c.Query("key_ingredients_id")
	// Convert string IDs to integers
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
		return
	}
	concernID, err := strconv.Atoi(concernIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid concern_id"})
		return
	}
	skinTypeID, err := strconv.Atoi(skinTypeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skin_type_id"})
		return
	}
	brandID, err := strconv.Atoi(brandIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand_id"})
		return
	}
	productTypeID, err := strconv.Atoi(productTypeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_type_id"})
		return
	}
	keyIngredientsID, err := strconv.Atoi(keyIngredientsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key_ingredients_id"})
		return
	}
	// Prepare the updated product struct
	updatedProduct := Product{
		ProductID:        productID,
		ProductName:      productName,
		AllIngredients:   allIngredients,
		ConcernID:        concernID,
		SkinTypeID:       skinTypeID,
		BrandID:          brandID,
		ProductTypeID:    productTypeID,
		KeyIngredientsID: keyIngredientsID,
	}
	// Execute the SQL statement with parameters
	_, err = db.Exec(`
    UPDATE PRODUCTS SET 
        Product_Name = ?, 
        All_Ingredients = ?, 
        Concern_ID = ?, 
        Skin_Type_ID = ?,
        Brand_ID = ?, 
        Product_Type_ID = ?, 
        Key_Ingredients_ID = ? 
    WHERE Product_ID = ?`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Send the response
	c.JSON(http.StatusOK, updatedProduct)
}

// Delete a product
func DeleteProduct(c *gin.Context, db *sql.DB) {
	id := c.Param("products_id")
	_, err := db.Exec("DELETE FROM PRODUCTS WHERE Product_ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
