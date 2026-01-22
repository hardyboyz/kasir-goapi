package routes

import (
	"net/http"
	"strconv"

	"go-api/models"

	"github.com/gin-gonic/gin"
)

var categories = []models.Category{
	{ID: 1, Name: "Electronics", Description: "Electronic devices"},
	{ID: 2, Name: "Books", Description: "Various books"},
}

func CategoryRoutes(r *gin.Engine) {
	r.GET("/categories", getCategories)
	r.POST("/categories", createCategory)
	r.GET("/categories/:id", getCategory)
	r.PUT("/categories/:id", updateCategory)
	r.DELETE("/categories/:id", deleteCategory)
}

func getCategories(c *gin.Context) {
	c.JSON(http.StatusOK, categories)
}

func createCategory(c *gin.Context) {
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)
	c.JSON(http.StatusCreated, newCategory)
}

func getCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	for _, category := range categories {
		if category.ID == id {
			c.JSON(http.StatusOK, category)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
}

func updateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, category := range categories {
		if category.ID == id {
			updatedCategory.ID = id
			categories[i] = updatedCategory
			c.JSON(http.StatusOK, updatedCategory)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
}

func deleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	for i, category := range categories {
		if category.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
}
