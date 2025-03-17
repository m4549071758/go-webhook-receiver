package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/webhook", func(c *gin.Context) {
		event := c.GetHeader("X-GitHub-Event")
		if event == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-GitHub-Event header"})
			return
		}

		if event != "push" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only push events are supported", "event": event})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "received push event"})

		gitCmd := exec.Command("git", "pull")
		gitCmd.Dir = "/root/blog"
		err := gitCmd.Run()
		if err != nil {
			log.Println("Failed to pull", err)
			return
		}

		buildCmd := exec.Command("yarn", "build")
		buildCmd.Dir = "/root/blog"
		err = buildCmd.Run()
		if err != nil {
			log.Println("Failed to build", err)
			return
		}

		err = exec.Command("pm2", "restart", "blog").Start()
		if err != nil {
			log.Println("Failed to restart", err)
			return
		}
	})

	router.Run(":5000")
}
