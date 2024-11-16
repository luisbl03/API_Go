package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luideiz/API_Go/api"
	"github.com/luideiz/API_Go/constants"
	//"fmt"
)

func main() {
	api := gin.Default()
	api.GET("/version", GetVersion) 
	api.POST("/signup", Register)
	api.Run(":8080")

}

func GetVersion(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json")
	c.JSON(200, gin.H{"version":"1.0.0"})
}

func Register(c *gin.Context) {
	//miramos si esta el cuerpo del mensaje
	if c.Request.Body == nil {
		c.JSON(400, gin.H{"error": "No request body found"})
		return
	}
	//miramos si el cuerpo del mensaje es un json con un username y un password
	var login map[string]string //declaramos un map con dos claves string
	err := c.BindJSON(&login)
	if err != nil {
		c.JSON(400, gin.H{"error": "No username or password found"})
		return
	}
	//miramos si estan los campos username y password
	if login["username"] == "" || login["password"] == "" {
		c.JSON(400, gin.H{"error": "No username or password found"})
		return
	}
	status := api.Register(login["username"], login["password"])
	if status == constants.EXISTS {
		c.JSON(400, gin.H{"error": "User already exists"})
		return
	}
	if status == constants.ERROR {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"status": "User created"})
	
}