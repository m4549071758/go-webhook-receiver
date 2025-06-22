package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"

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
		go func() {

			log.Println("========================================")
			log.Println("Received push event")
			log.Println("========================================")

			log.Println("----------------------------------------")
			log.Println("STARTING: Pulling latest code from GitHub")
			log.Println("----------------------------------------")
			gitCmd := exec.Command("git", "pull")
			gitCmd.Dir = "/root/blog"
			gitCmd.Stdout = log.Writer()
			gitCmd.Stderr = log.Writer()
			err := gitCmd.Run()
			if err != nil {
				log.Println("Failed to pull", err)
				return
			}
			log.Println("----------------------------------------")
			log.Println("FINISHED: Pulled latest code")
			log.Println("----------------------------------------")

			log.Println("----------------------------------------")
			log.Println("STARTING: Build process")
			log.Println("----------------------------------------")
			buildCmd := exec.Command("yarn", "build")
			buildCmd.Dir = "/root/blog"
			buildCmd.Stdout = log.Writer()
			buildCmd.Stderr = log.Writer()
			err = buildCmd.Run()
			if err != nil {
				log.Println("Failed to build", err)
				return
			}
			log.Println("----------------------------------------")
			log.Println("FINISHED: Built latest code")
			log.Println("----------------------------------------")
			time.Sleep(10 * time.Second)

			log.Println("----------------------------------------")
			log.Println("STARTING: Stopping the server")
			log.Println("----------------------------------------")
			pm2StopCmd := exec.Command("pm2", "delete", "blog")
			pm2StopCmd.Dir = "/root/blog"
			pm2StopCmd.Stdout = log.Writer()
			pm2StopCmd.Stderr = log.Writer()
			err = pm2StopCmd.Run()
			if err != nil {
				log.Println("Failed to stop", err)
				return
			}
			log.Println("----------------------------------------")
			log.Println("FINISHED: Stopped the server")
			log.Println("----------------------------------------")
			time.Sleep(5 * time.Second)

			log.Println("----------------------------------------")
			log.Println("STARTING: Starting the server")
			log.Println("----------------------------------------")
			pm2Cmd := exec.Command("pm2", "start", "ecosystem.config.js")
			pm2Cmd.Dir = "/root/blog"
			pm2Cmd.Stdout = log.Writer()
			pm2Cmd.Stderr = log.Writer()
			err = pm2Cmd.Run()
			if err != nil {
				log.Println("Failed to start", err)
				return
			}
			log.Println("----------------------------------------")
			log.Println("FINISHED: Started the server")
			log.Println("----------------------------------------")

			log.Println("========================================")
			log.Println("Deployed successfully")
			log.Println("========================================")
		}()
	})

	router.Run(":5000")
}
