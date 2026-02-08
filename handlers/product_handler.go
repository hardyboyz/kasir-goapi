package handlers

import (
	"myshop-api/models"
	"myshop-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HandleProducts - GET /api/product
func (h *ProductHandler) HandleProducts(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		h.GetAll(c)
	case "POST":
		h.Create(c)
	default:
		c.JSON(405, gin.H{"error": "Method not allowed"})
	}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	name := c.Query("name")
	categoryIDStr := c.Query("category_id")
	categoryID := 0
	var err error

	if categoryIDStr != "" {
		categoryID, err = strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid category_id"})
			return
		}
	}

	var products []models.Product
	if name != "" || categoryID > 0 {
		products, err = h.service.SearchProducts(name, categoryID)
	} else {
		products, err = h.service.GetAll()
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, products)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var product models.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.service.Create(&product)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, product)
}

// HandleProductByID - GET/PUT/DELETE /api/product/{id}
func (h *ProductHandler) HandleProductByID(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		h.GetByID(c)
	case "PUT":
		h.Update(c)
	case "DELETE":
		h.Delete(c)
	default:
		c.JSON(405, gin.H{"error": "Method not allowed"})
	}
}

// GetByID - GET /api/product/{id}
func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	err = c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	product.ID = id
	err = h.service.Update(&product)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, product)
}

// Delete - DELETE /api/product/{id}
func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}
