package main

import (
	"encoding/json"
	"context"
	"flag"
	"fmt"
	"strings"

	//"io/ioutil"
	//"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ayush6624/go-chatgpt"
	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

var GPTClient *chatgpt.Client

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	key := os.Getenv("OPENAI_KEY")

	GPTClient, err = chatgpt.NewClient(key)
	if err != nil {
		fmt.Println("Error creating chatgpt client")
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
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	s.Identify.Intents |= discordgo.IntentMessageContent

	if strings.HasPrefix(m.Content, "!hi") {
        GPTMessage("hi")
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "I received your message")
	if err != nil {
		fmt.Println(err)
	}
}

func GPTMessage(prompt string) {
	ctx := context.Background()

	res, err := GPTClient.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT35Turbo0613,
		Messages: []chatgpt.ChatMessage{
			{
				Role: chatgpt.ChatGPTModelRoleSystem,
				Content: prompt,
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	a, _ := json.MarshalIndent(res, "", "  ");
	fmt.Println(string(a))
}
