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

			log.Println("Received push event")
			log.Println("Pulling latest code from GitHub")
			gitCmd := exec.Command("git", "pull")
			gitCmd.Dir = "/root/blog"
			err := gitCmd.Run()
			if err != nil {
				log.Println("Failed to pull", err)
				return
			}
			log.Println("Pulled latest code")

			log.Println("Initiating build process")
			buildCmd := exec.Command("yarn", "build")
			buildCmd.Dir = "/root/blog"
			err = buildCmd.Run()
			if err != nil {
				log.Println("Failed to build", err)
				return
			}
			log.Println("Built latest code")
			time.Sleep(10 * time.Second)

			log.Println("Stopping the server")
			pm2StopCmd := exec.Command("pm2", "delete", "blog")
			pm2StopCmd.Dir = "/root/blog"
			err = pm2StopCmd.Run()
			if err != nil {
				log.Println("Failed to stop", err)
				return
			}
			log.Println("Stopped the server")
			time.Sleep(5 * time.Second)

			log.Println("Fetching articles for sitemap")
			sitemapCmd := exec.Command("npm", "run", "fetch-articles")
			sitemapCmd.Dir = "/root/blog"
			err = sitemapCmd.Run()
			if err != nil {
				log.Println("Failed to fetch articles for sitemap", err)
				return
			}
			log.Println("Got articles for sitemap")
			time.Sleep(5 * time.Second)

			log.Println("Generating sitemap")
			sitemapGenCmd := exec.Command("npm", "run", "generate-static-files")
			sitemapGenCmd.Dir = "/root/blog"
			err = sitemapGenCmd.Run()
			if err != nil {
				log.Println("Failed to generate sitemap", err)
				return
			}
			log.Println("Generated sitemap")
			time.Sleep(5 * time.Second)

			log.Println("Starting the server")
			pm2Cmd := exec.Command("pm2", "start", "ecosystem.config.js")
			pm2Cmd.Dir = "/root/blog"
			err = pm2Cmd.Run()
			if err != nil {
				log.Println("Failed to start", err)
				return
			}
			log.Println("Started the server")

			log.Println("Deployed successfully")
		}()
	})

	router.Run(":5000")
}
