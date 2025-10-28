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
	grid = append(grid, lines...)
	return grid
}

func Game(c *gin.Context) {
	// Load the game state
	cookie, err := c.Cookie("game")
	bytes := []byte(cookie)
	gs, jsonLoadErr := gamestate.FromJson(bytes)

	if _, reset := c.GetQuery("reset"); (err != nil) || reset || jsonLoadErr != nil {
		if jsonLoadErr != nil {
			panic(fmt.Errorf("error loading gamestate json: %w", jsonLoadErr))
		}

		// Initiialize a fresh game
		blankGame := gamestate.GameState{
			X: 0,
			Y: 0,
			Room: 0,
			State: 0,
		}
		gs.CurrGrid = loadMap(gs.GetCurrRoom().Path)

		blankGameJson, err := blankGame.ToJson()
		if err != nil {
			// TODO - dumb and lazy
			panic(err)
		}
		c.SetCookie("game", string(blankGameJson), cookieAge, "/", domain, false, true)
	}

	gs.CurrGrid = loadMap(gs.GetCurrRoom().Path)
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

	// Save the game state back
	j, err := gs.ToJson()
	if err != nil {
		// TODO dumb, lazy
		panic(err)
	}
	c.SetCookie("game", string(j), cookieAge, "/", domain, false, true)

	c.HTML(http.StatusOK, "game.html", gin.H{
		"title": "Game Page",
		"x": gs.X,
		"y": gs.Y,
		"room": gs.GetCurrRoom(),
		"gs": &gs,
		"xrange": gs.GetMapRange(gs.X),
		"yrange": gs.GetMapRange(gs.Y),
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
	err := r.SetTrustedProxies(nil)  // Only allow our own clientIP
	if err != nil {
		// We don't really have a way to recover from this.
		panic(err)
	}

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
	err = r.Run()
	if err != nil {
		// No way to recover from this, really.
		panic(err)
	}
}
