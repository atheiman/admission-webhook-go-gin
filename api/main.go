package main

import (
	"github.com/gin-gonic/gin"
	admission "k8s.io/api/admission/v1beta1"
	"log"
	"net/http"
)

type Person struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"-"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/admission-review", func(c *gin.Context) {
		var ar admission.AdmissionRequest

		if err := c.ShouldBindJSON(&ar); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println(ar.UID)
		log.Println(ar.Kind)

		resp := admission.AdmissionResponse{
			UID:     ar.UID,
			Allowed: true,
		}

		c.JSON(200, resp)
	})

	router.POST("/json-binding", func(c *gin.Context) {
		var person Person

		if err := c.ShouldBindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println(person.Name)
		log.Println(person.Age)

		if person.Name == "script-kiddie" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
			})
			return
		}

		c.JSON(200, gin.H{
			"welcome": person.Name,
		})
	})

	return router
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
