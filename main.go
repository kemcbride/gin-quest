package main

import (
	"fmt"
	"embed"
  "net/http"
	"strconv"
	"strings"

  "github.com/gin-contrib/static"
  "github.com/gin-gonic/gin"
)

//go:embed static
var server embed.FS

const cookieAge int = 3600 * 24 * 7 // 1 week?
const domain string = "kemcbride.noho.st"


type GameState struct {
	x int
	y int
	room int
}

func serializeGameState(g GameState) string {
	return fmt.Sprintf("%d,%d,%d", g.x, g.y, g.room)
}

func deserializeGameState(s string) GameState {
	segments := strings.Split(s, ",")
	var values [3]int
	for i, seg := range segments {
		val, err := strconv.Atoi(seg)
		if err != nil {
			// Lazy lazy bad
			panic(err)
		}
		values[i] = val
	}
	return GameState{
		x: values[0],
		y: values[1],
		room: values[2],
	}
}

func main() {
  // Create a Gin router with default middleware (logger and recovery)
  r := gin.Default()
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
	r.GET("/gin-quest", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Hello World Index",
		})
	})

	r.GET("/gin-quest/game", func(c *gin.Context) {

		gs := GameState{}
		cookie, err := c.Cookie("game")
		if err != nil {
			cookie = "0,0,0"
			c.SetCookie("game", cookie, cookieAge, "/", domain, false, true)
		} 
		gs = deserializeGameState(cookie)
		// Let's check the query path and respond to up, down, left, right.
		if _, up := c.GetQuery("up"); up {
			gs.y--
		} else if _, down := c.GetQuery("down"); down {
			gs.y++
		}
		if _, left := c.GetQuery("left"); left {
			gs.x--
		} else if _, right := c.GetQuery("right"); right {
			gs.x++
		}
		c.SetCookie("game", serializeGameState(gs), cookieAge, "/", domain, false, true)

		c.HTML(http.StatusOK, "game.html", gin.H{
			"title": "Game Page",
			"x": gs.x,
			"y": gs.y,
			"room": gs.room,
		})
	})
	// Serve static stuff so we can template it into html, etc
	r.Use(static.Serve("/static/", fs))

	// r.NoRoute(func(c *gin.Context) {
	// 	fmt.Printf("%s doesn't exist, redirect on /\n", c.Request.URL.Path)
	// 	c.Redirect(http.StatusMovedPermanently, "/gin-quest")
	// })

  // Start server on port 8080 (default)
  // Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
  r.Run()
}
