package main

import (
	// "fmt"
	"embed"
	// "path/filepath"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/kemcbride/gin-quest/internal/gamestate"
)

//go:embed static
var server embed.FS

//go:embed favi
var favicons embed.FS

const cookieAge int = 3600 * 24 * 7 // 1 week?
var domain string = "kemcbride.noho.st"

func Game(c *gin.Context) {
	// Load the game state
	cookie, err := c.Cookie("game")
	bytes := []byte(cookie)
	gsave, jsonLoadErr := gamestate.GameSaveFromJson(bytes)

	if _, reset := c.GetQuery("reset"); (err != nil) || reset || jsonLoadErr != nil {
		gsave = &gamestate.GameSave{
			X:       0,
			Y:       0,
			RoomKey: "mh04i224",
			State:   0,
		}

		// Initiialize a fresh game
		gs := gamestate.GameState{
			Save: *gsave,
		}
		gs.Room = gs.GetCurrRoom()
		gs.Room.LoadMap(server)
		gs.Room.LoadMeta(server)

		blankGameSaveJson, err := gs.Save.ToJson()
		if err != nil {
			// TODO - dumb and lazy
			panic(err)
		}
		c.SetCookie("game", string(blankGameSaveJson), cookieAge, "/", domain, false, true)
	}

	gs := gamestate.GameState{
		Save: *gsave,
	}
	gs.Room = gs.GetCurrRoom()
	gs.Room.LoadMap(server)
	gs.Room.LoadMeta(server)


	if _, talk := c.GetQuery("talk"); talk {
		_ = gs.Talk(server)
	} else if _, move := c.GetQuery("move"); move { // If doing talk, don't do other motion actions
		gs.Save.State = gamestate.StateExplore
	}

	if gs.Save.State == gamestate.StateExplore {
		// Handle Explore-based query options
		// if the query involved a portal, try that before other motion:
		if _, portal := c.GetQuery("portal"); portal {
			_ = gs.Portal(server)
		}
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
	}

	// Save the game state back
	j, err := gs.Save.ToJson()
	if err != nil {
		// TODO dumb, lazy
		panic(err)
	}
	c.SetCookie("game", string(j), cookieAge, "/", domain, false, true)

	c.HTML(http.StatusOK, "game", gin.H{
		"title":  "Game Page",
		"room":   gs.GetCurrRoom(),
		"gs":     &gs,
		"xrange": gs.GetMapRange(gs.Save.X, 3),
		"yrange": gs.GetMapRange(gs.Save.Y, 3),
	})
}

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// From https://github.com/gin-gonic/gin/issues/3675
		if strings.HasPrefix(c.Request.URL.Path, "/gin-quest/static/") || strings.HasPrefix(c.Request.URL.Path, "/favi") {
			c.Header("Cache-Control", "private, max-age=86400")
		}
		c.Next()
	}
}

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()
	r.HTMLRender = createMyRender()

	err := r.SetTrustedProxies(nil) // Only allow our own clientIP
	if err != nil {
		// We don't really have a way to recover from this.
		panic(err)
	}

	if os.Getenv("DOMAIN") != domain {
		domain = ""
	}

	fs, err := static.EmbedFolder(server, "static")
	if err != nil {
		panic(err)
	}

	// FFS for favicon file system, obviously!
	ffs, err := static.EmbedFolder(favicons, "favi")
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
		c.HTML(http.StatusOK, "beepboop", gin.H{
			"title": "Hello World Index",
		})
	})

	r.GET("/gin-quest/game", Game)
	r.GET("/gin-quest/", Game)

	r.Use(staticCacheMiddleware())
	// Serve static stuff so we can template it into html, etc
	r.Use(static.Serve("/gin-quest/static/", fs))
	r.Use(static.Serve("/", ffs))

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

// https://gin-gonic.com/en/docs/rendering/multiple-template/
func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "./templates/layouts/index.html")
	r.AddFromFiles(
		"game",
		"./templates/layouts/game.html",
		"templates/includes/map.html",
		"templates/includes/menu.html",
		"templates/includes/status.html",
		"templates/includes/conversation.html",
		"templates/includes/misc.html", // Resolves to top-level Menus (items, skills)
		"templates/includes/battle.html",
	)
	return r
}
