package msg

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	Reasons       []Reason
	ReasonChoices []*discordgo.ApplicationCommandOptionChoice
)

type Reason struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func LoadReasons() {
	jsonFile, err := os.Open("reasons.json")
	if err != nil {
		log.Panicf("Error opening reasons.json: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &Reasons)
	if err != nil {
		log.Panicf("Error unmarshalling reasons.json: %v", err)
	}

	for i, reason := range Reasons {
		ReasonChoices = append(ReasonChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  reason.Name,
			Value: strconv.Itoa(i),
		})
	}

	log.Print("Reasons loaded")
}
