package main

import (
	"BackEnd/Brand"
	"BackEnd/Concern"
	"BackEnd/Key_Ingredients"
	"BackEnd/Product_Type"
	"BackEnd/Products"
	"BackEnd/Skin_Type"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// struct definition
type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBName     string `json:"db_name"`
}

func loadConfig() (Config, error) {
	var config Config
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	return config, err
}

func main() {
	// Load configuration from JSON file
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Cannot load configuration: %v", err)
	}

	// Set up the Gin router
	router := gin.Default()

	// Create the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Brand CRUD routes
	router.GET("/brand", func(c *gin.Context) {
		brand.GetBrands(c, db)
	})
	router.POST("/brand/create", func(c *gin.Context) {
		brand.CreateBrand(c, db)
	})
	router.PUT("/brand/update", func(c *gin.Context) {
		brand.UpdateBrand(c, db)
	})
	router.DELETE("/brand/delete/:brand_id", func(c *gin.Context) {
		brand.DeleteBrand(c, db)
	})

	// Concern CRUD routes
	router.GET("/concerns", func(c *gin.Context) {
		concern.GetConcerns(c, db)
	})
	router.POST("/concerns/create", func(c *gin.Context) {
		concern.CreateConcern(c, db)
	})
	router.PUT("/concerns/update", func(c *gin.Context) {
		concern.UpdateConcern(c, db)
	})
	router.DELETE("/concerns/delete/:concern_id", func(c *gin.Context) {
		concern.DeleteConcern(c, db)
	})

	// Skin Type CRUD routes
	router.GET("/skin_type", func(c *gin.Context) {
		skin_type.GetSkinTypes(c, db)
	})
	router.POST("/skin_type/create", func(c *gin.Context) {
		skin_type.CreateSkinType(c, db)
	})
	router.PUT("/skin_type/update", func(c *gin.Context) {
		skin_type.UpdateSkinType(c, db)
	})
	router.DELETE("/skin_type/delete/:skin_type_id", func(c *gin.Context) {
		skin_type.DeleteSkinType(c, db)
	})

	// Product Type CRUD routes
	router.GET("/product_type", func(c *gin.Context) {
		product_type.GetProductTypes(c, db)
	})
	router.POST("/product_type/create", func(c *gin.Context) {
		product_type.CreateProductType(c, db)
	})
	router.PUT("/product_type/update", func(c *gin.Context) {
		product_type.UpdateProductType(c, db)
	})
	router.DELETE("/product_type/delete/:product_type_id", func(c *gin.Context) {
		product_type.DeleteProductType(c, db)
	})

	// Key Ingredients CRUD routes
	router.GET("/key_ingredients", func(c *gin.Context) {
		key_ingredients.GetKeyIngredients(c, db)
	})
	router.POST("/key_ingredients/create", func(c *gin.Context) {
		key_ingredients.CreateKeyIngredient(c, db)
	})
	router.PUT("/key_ingredients/update", func(c *gin.Context) {
		key_ingredients.UpdateKeyIngredient(c, db)
	})
	router.DELETE("/key_ingredients/delete/:key_ingredients_id", func(c *gin.Context) {
		key_ingredients.DeleteKeyIngredient(c, db)
	})

	// Products CRUD routes
	router.GET("/products", func(c *gin.Context) {
		products.GetProducts(c, db)
	})

	router.GET("/products/select/:concern_id/:skin_type_id", func(c *gin.Context) {
		// Get concern ID
		concernIDStr := c.Param("concern_id")
		concernID, err := strconv.Atoi(concernIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid concern_id"})
			return
		}
		// Get skin type ID
		skinTypeIDStr := c.Param("skin_type_id")
		skinTypeID, err := strconv.Atoi(skinTypeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skin_type_id"})
			return
		}
		products.GetSelectProducts(c, db, concernID, skinTypeID)
	})

	router.GET("/products/selectspec/:concern_id/:skin_type_id/:product_type_id", func(c *gin.Context) {
		// Get concern ID
		concernIDStr := c.Param("concern_id")
		concernID, err := strconv.Atoi(concernIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid concern_id"})
			return
		}

		// Get skin type ID
		skinTypeIDStr := c.Param("skin_type_id")
		skinTypeID, err := strconv.Atoi(skinTypeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skin_type_id"})
			return
		}

		// Get product type ID
		productTypeIDStr := c.Param("product_type_id")
		productTypeID, err := strconv.Atoi(productTypeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_type_id"})
			return
		}

		products.GetSelectProductsByType(c, db, concernID, skinTypeID, productTypeID)
	})

	router.POST("/products/create", func(c *gin.Context) {
		products.CreateProduct(c, db)
	})
	router.PUT("/products/update", func(c *gin.Context) {
		products.UpdateProduct(c, db)
	})
	router.DELETE("/products/delete/:products_id", func(c *gin.Context) {
		products.DeleteProduct(c, db)
	})
	router.Run(":8080")
}
