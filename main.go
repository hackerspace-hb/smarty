package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"smarty/discord"
	"smarty/githubcli"
	"syscall"
	"time"
)

func loadConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error while reading config file: %w \n", err))
	}
}

func setupLog() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func updateDiscordChannelText() {
	var text = ""
	for _, repo := range gi.GetReposByStars() {
		log.Info(repo.Name)
		text += repo.Name + "\nRate: " + string(repo.Stars) + "🌟\n\n" + repo.Url + "\n"
	}
	log.Info("Update text to: " + text)
	dg.SendToChoosenChannel(text)
}

var dg discord.Discord
var gi githubcli.GitHubInterface

func main() {
	loadConfig()
	setupLog()
	log.Info("Init")

	//setup discord and github cli
	dg.Setup()
	gi.SetOrgaName("hackerspace-hb")

	//oneshot text update
	updateDiscordChannelText()

	//setup scheduler/timer
	setupTimer()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func setupTimer() {
	ticker := time.NewTicker(60 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				updateDiscordChannelText()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
