package main

import (
	"fmt"
	"embed"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kemcbride/gin-quest/internal/gamestate"
)

//go:embed static
var server embed.FS

const cookieAge int = 3600 * 24 * 7 // 1 week?
const domain string = "kemcbride.noho.st"


func loadMap(roomPath string) []string {
	file, err := server.ReadFile(fmt.Sprintf("static/%s", roomPath))
	if err != nil {
		panic(err)
	}

	var grid []string
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		grid = append(grid, line)
	}
	return grid
}

func Game(c *gin.Context) {
	cookie, err := c.Cookie("game")
	if _, reset := c.GetQuery("reset"); (err != nil) || reset {
		cookie = "0,0,0"
		c.SetCookie("game", cookie, cookieAge, "/", domain, false, true)
	}

	gs := gamestate.DeserializeGameState(cookie)
	grid := loadMap(gs.GetCurrRoom().Path)
	gs.SetGrid(grid)
	// Let's check the query path and respond to up, down, left, right.
	if _, up := c.GetQuery("up"); up {
		gs.MoveUp()
	} else if _, down := c.GetQuery("down"); down {
		gs.MoveDown()
	}
	if _, left := c.GetQuery("left"); left {
		gs.MoveLeft()
	} else if _, right := c.GetQuery("right"); right {
		gs.MoveRight()
	}
	c.SetCookie("game", gamestate.SerializeGameState(gs), cookieAge, "/", domain, false, true)

	c.HTML(http.StatusOK, "game.html", gin.H{
		"title": "Game Page",
		"x": gs.GetX(),
		"y": gs.GetY(),
		"room": gs.GetCurrRoom(),
		"gs": &gs,
		"xrange": gs.GetMapRange(gs.GetX()),
		"yrange": gs.GetMapRange(gs.GetY()),
	})
}

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// From https://github.com/gin-gonic/gin/issues/3675
		if strings.HasPrefix(c.Request.URL.Path, "/gin-quest/static/") {
			c.Header("Cache-Control", "private, max-age=86400")
		}
		c.Next()
	}
}

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()
	r.SetTrustedProxies(nil)  // Only allow our own clientIP

	r.LoadHTMLGlob("templates/*")
	fs, err := static.EmbedFolder(server, "static")
	if err != nil {
		panic(err)
	}

	// Define a simple GET endpoint
	r.GET("/gin-quest/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/gin-quest/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Hello World Index",
		})
	})

	r.GET("/gin-quest/game", Game)
	r.GET("/gin-quest/", Game)

	r.Use(staticCacheMiddleware())
	// Serve static stuff so we can template it into html, etc
	r.Use(static.Serve("/gin-quest/static/", fs))

	// r.NoRoute(func(c *gin.Context) {
	// 	fmt.Printf("%s doesn't exist, redirect on /\n", c.Request.URL.Path)
	// 	c.Redirect(http.StatusMovedPermanently, "/gin-quest")
	// })

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
