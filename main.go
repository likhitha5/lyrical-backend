package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Song struct
type Song struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Artist         string `json:"artist"`
	MusicDirector  string `json:"music_director"`  // New Field
	Lyrics         string `json:"lyrics"`
}

var songs []Song

// Load songs from data.json
func loadSongs() error {
	file, err := os.Open("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &songs)
	if err != nil {
		return err
	}

	fmt.Println("✅ Loaded songs:", len(songs))
	return nil
}

func main() {
	// Load songs from file
	err := loadSongs()
	if err != nil {
		fmt.Println("❌ Error loading songs:", err)
		return
	}

	// Initialize Gin router
	r := gin.Default()
	r.Use(cors.Default()) // Enable CORS

	// Get all songs
	r.GET("/songs", func(c *gin.Context) {
		c.JSON(http.StatusOK, songs)
	})

	// Get song by ID
	r.GET("/song/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
			return
		}

		for _, song := range songs {
			if song.ID == id {
				c.JSON(http.StatusOK, song)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
	})

	// Search songs by title or artist
	r.GET("/search", func(c *gin.Context) {
		query := strings.ToLower(c.Query("q")) // Get search query
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty search query"})
			return
		}

		var results []Song
		for _, song := range songs {
			if strings.Contains(strings.ToLower(song.Title), query) || strings.Contains(strings.ToLower(song.Artist), query) {
				results = append(results, song)
			}
		}

		c.JSON(http.StatusOK, results)
	})

	// Start server
	fmt.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}
