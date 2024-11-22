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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// struct definition
type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBName     string `json:"db_name"`
	DBPort     string `json:"db_port"`
}

func loadConfig() (Config, error) {
	var config Config
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, fmt.Errorf("reading config file: %v", err)
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("parsing config file: %v", err)
	}
	return config, err
}

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting application...")

	// Set up Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Load configuration from JSON file
	config, err := loadConfig()
	if err != nil {
		log.Printf("Config load error: %v", err)
		log.Fatalf("Cannot load configuration: %v", err)
	}
	log.Printf("Config loaded successfully")

	// For App Engine, modify your DSN to use Unix socket
	var dsn string
	if os.Getenv("GAE_ENV") == "standard" {
		// Running on App Engine
		dsn = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s",
			config.DBUser,
			config.DBPassword,
			config.DBHost, // This should be your instance connection name in config
			config.DBName)
	} else {
		// Local development
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBName)
	}

	// Configure database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Printf("Database ping error: %v", err)
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Printf("Database connected successfully!")

	// Add basic health check endpoint
	router.GET("/health", func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Database connection failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is healthy",
		})
	})

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
		concernIDStr := c.Param("concern_id")
		concernID, err := strconv.Atoi(concernIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid concern_id"})
			return
		}
		skinTypeIDStr := c.Param("skin_type_id")
		skinTypeID, err := strconv.Atoi(skinTypeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skin_type_id"})
			return
		}
		products.GetSelectProducts(c, db, concernID, skinTypeID)
	})

	router.GET("/products/selectspec/:concern_id/:skin_type_id/:product_type_id", func(c *gin.Context) {
		concernIDStr := c.Param("concern_id")
		concernID, err := strconv.Atoi(concernIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid concern_id"})
			return
		}

		skinTypeIDStr := c.Param("skin_type_id")
		skinTypeID, err := strconv.Atoi(skinTypeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skin_type_id"})
			return
		}

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// Enhanced logging for database connection
	log.Printf("Environment: %s", os.Getenv("GAE_ENV"))
	log.Printf("Database host: %s", config.DBHost)

	// Modify your products endpoint with logging
	router.GET("/products/select/:concern_id/:skin_type_id", func(c *gin.Context) {
		log.Printf("Received request for products with concern_id: %s, skin_type_id: %s",
			c.Param("concern_id"), c.Param("skin_type_id"))

		concernIDStr := c.Param("concern_id")
		concernID, err := strconv.Atoi(concernIDStr)
		if err != nil {
			log.Printf("Error converting concern_id: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid concern_id"})
			return
		}

		skinTypeIDStr := c.Param("skin_type_id")
		skinTypeID, err := strconv.Atoi(skinTypeIDStr)
		if err != nil {
			log.Printf("Error converting skin_type_id: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skin_type_id"})
			return
		}

		log.Printf("Calling GetSelectProducts with concernID: %d, skinTypeID: %d",
			concernID, skinTypeID)
		products.GetSelectProducts(c, db, concernID, skinTypeID)
	})

	// Add middleware to log all requests
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
}
