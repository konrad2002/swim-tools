package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"swim-tools/dto"
	"swim-tools/service"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run() {
	port := os.Getenv("SERVER_PORT")

	if port == "" {
		fmt.Println("no application port given! Please set SERVER_PORT. Using 8080")
		port = "8080"
	}

	router.GET("/actuator", actuator)

	router.POST("/convert", convert)

	fmt.Println("Starting server on port " + port)

	err := router.Run(":" + port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}

func actuator(c *gin.Context) {
	state := "OPERATIONAL"

	c.String(http.StatusOK, state)
}

func convert(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is required",
		})
		return
	}

	metadataStr := c.PostForm("metadata")

	var request dto.DsvLenexConvertRequest
	if err := json.Unmarshal([]byte(metadataStr), &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid metadata",
		})
		return
	}

	err = c.SaveUploadedFile(file, "./uploads/"+file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err1 := service.Convert(request.Filename, request.Origin, request.Target)

	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err1.Error(),
		})
	} else {
		c.String(http.StatusOK, "File uploaded successfully: %s", file.Filename)
	}

}
