package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"

	"github.com/bwmarrin/discordgo"
)

// func init() {
// 	flag.StringVar(&token, "t", "", "Bot Token")
// 	flag.Parse()
// }

var token string
var adminUsername string
var adminPassword string

func main() {

	// Load Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token = os.Getenv("DISCORD_TOKEN")
	adminUsername = os.Getenv("API_USERNAME")
	adminPassword = os.Getenv("API_PASSWORD")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// DEBUG - MOVE ME
	client := &http.Client{}

	var returnString string

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "current" {
		// DEBUG - Change this to the output of the API
		var track *spotify.FullTrack

		resp, err := http.Get("http://localhost:8888/tracks/current")
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal([]byte(string(body)), &track)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, artist := range track.Album.Artists {
			returnString = fmt.Sprintf("%s%s", returnString, artist.Name+", ")
		}
		// String Fix
		returnString = returnString[:len(returnString)-2]
		returnString = track.Name + " - " + returnString
	}

	if strings.HasPrefix(m.Content, "search") {
		var apiSearchOutput APIResponseSearchOutput

		searchTerm := strings.Replace(m.Content, "search ", "", -1)

		resp, err := http.Get("http://localhost:8888/search/" + searchTerm)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal([]byte(string(body)), &apiSearchOutput)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, track := range apiSearchOutput.Tracks {
			returnString = fmt.Sprintf("%s%s", returnString, track.Name+" - "+track.Artist+" - "+track.ID+"\n")
		}
	}

	if strings.HasPrefix(m.Content, "add") {

		songId := strings.Replace(m.Content, "add ", "", -1)

		addTrackInput := APIAddTrack{
			URI: spotify.URI("spotify:track:" + songId),
		}
		postBody, _ := json.Marshal(addTrackInput)

		req, err := http.NewRequest("POST", "http://localhost:8888/tracks/add", bytes.NewBuffer(postBody))

		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(adminUsername, adminPassword)
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		returnString = string(body)
	}

	if strings.HasPrefix(m.Content, "skip") {
		req, err := http.NewRequest("POST", "http://localhost:8888/player/skip", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(adminUsername, adminPassword)
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		returnString = string(body)
	}

	if strings.HasPrefix(m.Content, "play") {
		req, err := http.NewRequest("POST", "http://localhost:8888/player/play", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(adminUsername, adminPassword)
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		returnString = string(body)
	}

	if strings.HasPrefix(m.Content, "pause") {
		req, err := http.NewRequest("POST", "http://localhost:8888/player/pause", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(adminUsername, adminPassword)
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		returnString = string(body)
	}

	s.ChannelMessageSend(m.ChannelID, returnString)
}
