package main

import (
	"log"
	"time"
	"encoding/json"
	"net/http"
	"sort"
	//"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)
var recentCommits []GHApiCommit

type GHAuthor struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
}

type GHCommit struct {
	Author  GHCommitAuthor `json:"author"`
	Message string         `json:"message"`
}

type GHCommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type GHApiCommit struct {
	Author  GHAuthor `json:"author"`
	Commit  GHCommit `json:"commit"`
	HtmlUrl string   `json:"html_url"`
}

type GHCommitsByDate []GHApiCommit

func (d GHCommitsByDate) Len() int      { return len(d) }
func (d GHCommitsByDate) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d GHCommitsByDate) Less(i, j int) bool {
	time1, _ := time.Parse(time.RFC3339, d[i].Commit.Author.Date)
	time2, _ := time.Parse(time.RFC3339, d[j].Commit.Author.Date)
	return time1.After(time2)
}


func backgroundTask() {
	for {
		time.Sleep(15 * time.Minute)
		//getRecentsCommits()
	}
}

func main() {
	// Load dotenv
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found!")
	}
	//discordToken = os.Getenv("DISCORD_TOKEN")
	//discordStatusChannel = os.Getenv("DISCORD_STATUS_CHANNEL")
	//discordUpdateChannel = os.Getenv("DISCORD_UPDATES_CHANNEL")

	// get recent commits
	getRecentsCommits()

	// start Discord bot
	//go startDiscordBot()

	// start background task
	//go backgroundTask()

	// initialize fiber app
	app := fiber.New()
	app.Use(cors.New())

	// root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("online")
	})

	// status endpoint
	//app.Get("/status", func(c *fiber.Ctx) error {
		//return c.JSON(currentStatus)
	//})

	// latest updates endpoint
	//app.Get("/updates", func(c *fiber.Ctx) error {
		//return c.JSON(currentUpdate)
	//})

	// recent commits endpoint
	app.Get("/commits", func(c *fiber.Ctx) error {
		return c.JSON(recentCommits)
	})

	// start fiber app
	app.Listen(":3000")
}

func getRecentsCommits() {
	githubCommitApis := []string{
		"https://api.github.com/repos/DevevolperPlus/ScratchTurbo-ObjectLibraries/commits?per_page=50",
	}

	var newRecentCommits []GHApiCommit
	for i := 0; i < len(githubCommitApis); i++ {
		resp, err := http.Get(githubCommitApis[i])
		if err != nil {
			//log.Errorf("Failed fetching %s: %s", githubCommitApis[i], err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			//log.Errorf("Failed fetching %s: Non-OK status code; %s", githubCommitApis[i], strconv.Itoa(resp.StatusCode))
			continue
		}

		var apiResp []GHApiCommit
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		if err != nil {
			//log.Errorf("Failed decoding response from %s: %s", githubCommitApis[i], err)
			continue
		}

		newRecentCommits = append(newRecentCommits, apiResp...)
	}

	sort.Sort(GHCommitsByDate(newRecentCommits))
	if len(newRecentCommits) >= 200 {
		recentCommits = newRecentCommits[:200]
	} else {
		recentCommits = newRecentCommits
	}
}
