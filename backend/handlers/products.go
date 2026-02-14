package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/SolieSoftware/product-feeds/db"
	"github.com/SolieSoftware/product-feeds/models"
	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	// Read query params for pagination
	cursor := c.Query("cursor")
	limitStr := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	// Build query based on whether a cursor was provided
	var query string
	var args []interface{}

	if cursor != "" {
		query = `SELECT id, name, price, currency, image_paths, source_url, company_name, scraped_at
				 FROM products
				 WHERE scraped_at < (SELECT scraped_at FROM products WHERE id = $1)
				 ORDER BY scraped_at DESC
				 LIMIT $2`
		args = []interface{}{cursor, limit}
	} else {
		query = `SELECT id, name, price, currency, image_paths, source_url, company_name, scraped_at
				 FROM products
				 ORDER BY scraped_at DESC
				 LIMIT $1`
		args = []interface{}{limit}
	}

	// Execute the query
	rows, err := db.Pool.Query(c.Request.Context(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
		return
	}
	defer rows.Close()

	// Scan rows into Product structs
	var products []models.Product
	for rows.Next() {
		var p models.Product
		var imagePaths []byte // image_paths is jsonb, comes as raw bytes

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Currency, &imagePaths, &p.SourceURL, &p.CompanyName, &p.ScrapedAt)
		if err != nil {
			log.Println("Scan error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product"})
			return
		}

		// Parse the jsonb image_paths into a string slice
		json.Unmarshal(imagePaths, &p.ImagePaths)

		products = append(products, p)
	}

	// Build the next cursor from the last product
	var nextCursor string
	if len(products) == limit {
		nextCursor = products[len(products)-1].ID
	}

	c.JSON(http.StatusOK, gin.H{
		"products":    products,
		"next_cursor": nextCursor,
	})
}
