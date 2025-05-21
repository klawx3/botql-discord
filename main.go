package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type ServerCommand int

const DiscordEnvVar string = "discord_token"

const (
	RebootServer ServerCommand = iota
)

var commandName = map[string]ServerCommand{
	"reboot": RebootServer,
}

var allowedDiscordUsers = [...]string{
	"560168832296288266", // kronoz
	"201899463613218816", // k64
}

func commandFromString(command string) (ServerCommand, bool) {
	d, ok := commandName[command]
	return d, ok
}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(".env file not loaded")
	}

	authToken := os.Getenv(DiscordEnvVar)

	if len(authToken) == 0 {
		log.Fatal("no enviroment variable " + DiscordEnvVar)
	}

	discord, err := discordgo.New("Bot " + authToken)

	if err != nil {
		log.Fatal("Error with discord bot load", err)
		panic(err)
	}

	discord.AddHandler(messageHandler)

	discord.Open()
	defer discord.Close()

	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	hasPrefix := !strings.HasPrefix(m.Content, "!ql")

	if hasPrefix {
		return
	}

	if !isDiscordUserAllowed(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "No tienes permisos para ejecutar comandos")
		return
	} else {
		s.ChannelMessageSend(m.ChannelID, "Permisos OK")
	}

	tokens := strings.Split(m.Content, " ")
	if len(tokens) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Mal formato de mensaje de comando , ejemplo: !ql <commando>")
		return
	}

	commandToken := tokens[1]
	userCommand, ok := commandFromString(commandToken)

	if !ok {
		message := fmt.Sprintf("el comando %s no existe", commandToken)
		s.ChannelMessageSend(m.ChannelID, message)
		return
	}

	var commandExecutionStatus string
	err := executeCommand(userCommand)
	if err == nil {
		commandExecutionStatus = fmt.Sprintf("comando %s ejecutado correctamente", commandToken)
	} else {
		commandExecutionStatus = "Error en la ejecucion: `" + err.Error() + "`"
	}
	fmt.Println(commandExecutionStatus)
	s.ChannelMessageSend(m.ChannelID, commandExecutionStatus)
}

func isDiscordUserAllowed(userId string) bool {
	isAllowed := false

	for i, _ := range allowedDiscordUsers {
		if allowedDiscordUsers[i] == userId {
			isAllowed = true
			break
		}
	}
	return isAllowed
}

func executeCommand(command ServerCommand) error {
	switch command {
	case RebootServer:
		{
			cmd := exec.Command("", `javascript:alert("Messagentoo!");close();`)
			_, err := cmd.Output()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}
