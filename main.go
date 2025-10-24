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


type Room struct {
	Id int
	Name string
	Path string
}

type GameState struct {
	x int
	y int
	room int
	currGrid []string
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

func (gs *GameState) GetStatusBlurb() string {
	return "Status: All good!"
}

func (gs *GameState) GetX() int {
	return gs.x
}

func (gs *GameState) GetY() int {
	return gs.y
}

func (gs *GameState) AddX(a int) int {
	return gs.x + a
}

func (gs *GameState) AddY(a int) int {
	return gs.y + a
}

func (gs *GameState) CanMove(dx int, dy int) bool {
	// For now, we just want to check if it's water or mountain
	// Or, the edge of the map.
	// Eventually could handle a case where we fly or swim.
	loc := gs.GetGridLoc(dx, dy)
	switch loc {
	case	".":
		return true
	case "^":
		return false
	case "~":
		return false
	case "Q": // placeholder for edge of map?
		return false
	default:
		return true
	}
}

func (gs *GameState) MoveUp() {
	if gs.CanMove(gs.x, gs.y-1) {
		gs.y--;
	}
}

func (gs *GameState) MoveDown() {
	if gs.CanMove(gs.x, gs.y+1) {
		gs.y++;
	}
}

func (gs *GameState) MoveLeft() {
	if gs.CanMove(gs.x-1, gs.y) {
		gs.x--;
	}
}

func (gs *GameState) MoveRight() {
	if gs.CanMove(gs.x+1, gs.y) {
		gs.x++;
	}
}

func (gs *GameState) GetGrid() []string {
	return gs.currGrid
}

func (gs *GameState) SetGrid(grid []string) {
	gs.currGrid = grid
}

func (gs *GameState) GetGridLoc(x int, y int) string {
	adjustedX := max(0, len(gs.currGrid[0]) / 2 + x)
	adjustedY := max(0, len(gs.currGrid) / 2 + y)
	// Dumb exit to send weird letter we can map to some style
	if ( adjustedX < 0 || adjustedX >= len(gs.currGrid[0]) ) || (adjustedY < 0 || adjustedY >= len(gs.currGrid) ) {
		return "Q"
	}
	loc := gs.GetGrid()[adjustedY][adjustedX]
	return string(loc)
}

func (gs *GameState) GetGridLocColor(x int, y int) string {
	var colorMap = map[string]string{
		".": "#FFE4B5",
		"^": "#444444",
		"~": "#4682B4",
		"#": "#9ACD32",
		"Q": "#000000",
	}
	return colorMap[gs.GetGridLoc(x, y)]
}

func (gs *GameState) GetGridLocClass(x int, y int) string {
	var classMap = map[string]string{
		".": "desert",
		"^": "mountain",
		"~": "water",
		"#": "grass",
		"Q": "abyss",
	}
	return classMap[gs.GetGridLoc(x, y)]
}

func (gs *GameState) GetRoomHash() map[int]Room {
	var roomMap = map[int]Room {
		0: Room{Id: 0, Name: "Continent of Euniciar", Path: "map-mh04i224.txt"},
		1: Room{Id: 1, Name: "Land of Patricolia", Path: "map-mh04dw5i.txt"},
	}
	return roomMap
}

func (gs *GameState) GetRoom(i int) Room {
	return gs.GetRoomHash()[i]
}

func (gs *GameState) GetCurrRoom() Room {
	return gs.GetRoomHash()[gs.room]
}

func (gs *GameState) GetCurrRoomName() string {
	return gs.GetRoomHash()[gs.room].Name
}

func (gs *GameState) GetMapRange(coord int) []int {
	var coordRange []int
	for i := coord - 2; i < coord + 3; i++ {
		coordRange = append(coordRange, i)
	}
	return coordRange
}

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

	gs := deserializeGameState(cookie)
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
	c.SetCookie("game", serializeGameState(gs), cookieAge, "/", domain, false, true)

	c.HTML(http.StatusOK, "game.html", gin.H{
		"title": "Game Page",
		"x": gs.x,
		"y": gs.y,
		"room": gs.room,
		"gs": &gs,
		"xrange": gs.GetMapRange(gs.x),
		"yrange": gs.GetMapRange(gs.y),
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
