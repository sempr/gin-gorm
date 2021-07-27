package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// App hold an gin.Engine and a gorm.DB
type App struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

// Initialize the app with db config
func (a *App) Initialize(dbType, dbString string) {
	var err error
	a.DB, err = gorm.Open(sqlite.Open(dbString))
	if err != nil {
		log.Fatal(err)
	}
	a.Engine = gin.Default()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Engine.GET("/products", a.getProducts)
	a.Engine.POST("/product", a.createProduct)
	a.Engine.GET("/product/:id", a.getProduct)
	a.Engine.PUT("/product/:id", a.updateProduct)
	a.Engine.DELETE("/product/:id", a.deleteProduct)

}

// Run will start the app
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Engine))
}

func (a *App) getProduct(c *gin.Context) {
	var p product
	if err := c.BindUri(&p); err != nil {
		c.JSON(http.StatusNotFound, "Product Not Found")
		return
	}
	if err := p.getProduct(a.DB); err != nil {
		switch err.Error() {
		case "record not found":
			c.JSON(http.StatusNotFound, gin.H{"message": "Product Not Found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	respondWithJSON(c, http.StatusOK, p)
	return
}

func (a *App) getProducts(c *gin.Context) {
	var obj struct {
		Count int `json:"count"`
		Start int `json:"start"`
	}
	c.BindQuery(&obj)
	count, start := obj.Count, obj.Start
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(c, http.StatusOK, products)
}

func respondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}

func (a *App) createProduct(c *gin.Context) {
	var p product
	if err := c.BindJSON(&p); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := p.createProduct(a.DB); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(c, http.StatusCreated, p)
}

func (a *App) updateProduct(c *gin.Context) {
	var obj struct {
		ID uint `uri:"id"`
	}
	c.BindUri(&obj)

	var p product
	if err := c.BindJSON(&p); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	p.ID = obj.ID

	if err := p.updateProduct(a.DB); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(c, http.StatusOK, p)
}

func (a *App) deleteProduct(c *gin.Context) {
	var obj struct {
		ID uint `uri:"id"`
	}
	c.BindUri(&obj)

	var p product
	p.ID = obj.ID
	if err := p.deleteProduct(a.DB); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(c, http.StatusOK, gin.H{"result": "success"})
}
