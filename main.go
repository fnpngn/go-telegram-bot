package main

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"strings"
	"time"
)

const (
	CHOOSING      = "CHOOSING"
	TYPING_REPLY  = "TYPING_REPLY"
	TYPING_CHOICE = "TYPING_CHOICE"
)

var (
	replyKeyboard = [][]telebot.ReplyButton{
		{{Text: "Age"}, {Text: "Favourite colour"}},
		{{Text: "Number of siblings"}, {Text: "Something else..."}},
		{{Text: "Done"}},
	}
)

type UserData map[string]string

func factsToStr(data UserData) string {
	var facts []string
	for k, v := range data {
		facts = append(facts, fmt.Sprintf("%s - %s", k, v))
	}
	return strings.Join(facts, "\n")
}

func main() {
	pref := telebot.Settings{
		Token:  "TOKEN",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	userStates := make(map[int64]string)
	userData := make(map[int64]UserData)

	bot.Handle("/start", func(c telebot.Context) error {
		chatID := c.Chat().ID
		userStates[chatID] = CHOOSING

		if _, exists := userData[chatID]; !exists {
			userData[chatID] = make(UserData)
		}

		replyText := "Hi! My name is Doctor Botter."
		if len(userData[chatID]) > 0 {
			replyText += " You already told me your " + strings.Join(mapKeys(userData[chatID]), ", ") + "."
			replyText += " Why don't you tell me something more about yourself? Or change anything I already know."
		} else {
			replyText += " I will hold a more complex conversation with you. Why don't you tell me something about yourself?"
		}

		return c.Send(replyText, telebot.ReplyMarkup{ReplyKeyboard: replyKeyboard})
	})

	bot.Handle("Age", func(c telebot.Context) error {
		return handleChoice(c, userStates, userData, "Age")
	})
	bot.Handle("Favourite colour", func(c telebot.Context) error {
		return handleChoice(c, userStates, userData, "Favourite colour")
	})
	bot.Handle("Number of siblings", func(c telebot.Context) error {
		return handleChoice(c, userStates, userData, "Number of siblings")
	})
	bot.Handle("Something else...", func(c telebot.Context) error {
		userStates[c.Chat().ID] = TYPING_CHOICE
		return c.Send("Alright, please send me the category first, for example \"Most impressive skill\".")
	})
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		chatID := c.Chat().ID
		state := userStates[chatID]

		if state == TYPING_CHOICE {
			userStates[chatID] = TYPING_REPLY
			userData[chatID]["choice"] = c.Text()
			return c.Send(fmt.Sprintf("Your category \"%s\"? Yes, I would love to hear about that!", c.Text()))
		} else if state == TYPING_REPLY {
			choice := userData[chatID]["choice"]
			userData[chatID][choice] = strings.ToLower(c.Text())
			delete(userData[chatID], "choice")

			userStates[chatID] = CHOOSING
			return c.Send("Neat! Just so you know, this is what you already told me:\n"+factsToStr(userData[chatID]),
				telebot.ReplyMarkup{ReplyKeyboard: replyKeyboard})
		}
		return nil
	})
	bot.Handle("Done", func(c telebot.Context) error {
		chatID := c.Chat().ID
		state := userStates[chatID]
		if state == TYPING_REPLY {
			delete(userData[chatID], "choice")
		}

		return c.Send("I learned these facts about you:\n"+factsToStr(userData[chatID])+"\nUntil next time!",
			&telebot.ReplyMarkup{RemoveKeyboard: true})
	})

	log.Println("Bot is running...")
	bot.Start()
}

func handleChoice(c telebot.Context, userStates map[int64]string, userData map[int64]UserData, choice string) error {
	chatID := c.Chat().ID
	userStates[chatID] = TYPING_REPLY
	userData[chatID]["choice"] = choice
	if val, exists := userData[chatID][choice]; exists {
		return c.Send(fmt.Sprintf("Your %s? I already know the following about that: %s", choice, val))
	}
	return c.Send(fmt.Sprintf("Your %s? Yes, I would love to hear about that!", choice))
}

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
