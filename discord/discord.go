package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Discord struct {
	session *discordgo.Session
}

func (dg *Discord) Setup() {
	var err error
	dg.session, err = discordgo.New("Bot " + viper.GetString("discord.token"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.session.AddHandler(dg.messageCreate)

	dg.session.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
}

func (dg *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!register" {

		msg, _ := s.ChannelMessageSend(m.ChannelID, "-->placeholder<--")
		log.Info("Register: ", msg.ID, msg.ChannelID)
		viper.Set("discord.messageid", msg.ID)
		viper.Set("discord.channelid", msg.ChannelID)
		viper.WriteConfig()
	}

	if m.Content == "!test" {
		dg.SendToChoosenChannel("test back")
	}
}

func (dg *Discord) SendToChoosenChannel(text string) {
	var mid = viper.GetString("discord.messageid")
	var cid = viper.GetString("discord.channelid")
	dg.session.ChannelMessageEdit(cid, mid, text)
}
func (dg *Discord) Close() {
	dg.Close()
}
