package main

import (
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

const (
	CHOOSING      = "choosing"
	TYPING_REPLY  = "typing_reply"
	TYPING_CHOICE = "typing_choice"
)

var (
	replyKeyboard = [][]telebot.Btn{
		{telebot.Btn{Text: "Age"}, telebot.Btn{Text: "Favourite colour"}},
		{telebot.Btn{Text: "Number of siblings"}, telebot.Btn{Text: "Something else..."}},
		{telebot.Btn{Text: "Done"}},
	}
)

type UserData map[string]string

// Helper function for formatting the gathered user info
func factsToStr(userData UserData) string {
	var facts []string
	for key, value := range userData {
		facts = append(facts, key+" - "+value)
	}
	return strings.Join(facts, "\n")
}

func main() {
	// Initialize the bot
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TOKEN"), // Use your bot token here
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	userData := make(map[int64]UserData)
	state := make(map[int64]string)

	// /start command handler
	bot.Handle("/start", func(c telebot.Context) error {
		uid := c.Sender().ID
		data, exists := userData[uid]
		if !exists {
			data = make(UserData)
			userData[uid] = data
		}

		var reply string
		if len(data) > 0 {
			reply = "Hi! My name is Doctor Botter. You already told me your " +
				strings.Join(keys(data), ", ") + ". Why don't you tell me more about yourself?"
		} else {
			reply = "Hi! My name is Doctor Botter. I will hold a more complex conversation with you. Why don't you tell me something about yourself?"
		}

		state[uid] = CHOOSING
		return c.Reply(reply, telebot.ReplyMarkup{ReplyKeyboard: replyKeyboard})
	})

	// Handle predefined choices
	bot.Handle(&telebot.Btn{Text: "Age"}, choiceHandler(bot, "Age", userData, state))
	bot.Handle(&telebot.Btn{Text: "Favourite colour"}, choiceHandler(bot, "Favourite colour", userData, state))
	bot.Handle(&telebot.Btn{Text: "Number of siblings"}, choiceHandler(bot, "Number of siblings", userData, state))

	// Custom category
	bot.Handle(&telebot.Btn{Text: "Something else..."}, func(c telebot.Context) error {
		uid := c.Sender().ID
		state[uid] = TYPING_CHOICE
		return c.Reply("Alright, please send me the category first, for example 'Most impressive skill'")
	})

	// Done handler
	bot.Handle(&telebot.Btn{Text: "Done"}, func(c telebot.Context) error {
		uid := c.Sender().ID
		data := userData[uid]
		reply := "I learned these facts about you: " + factsToStr(data) + "\nUntil next time!"
		delete(state, uid)
		return c.Reply(reply, &telebot.ReplyMarkup{RemoveKeyboard: true})
	})

	// Handle text messages for choice and reply
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		uid := c.Sender().ID
		text := c.Text()

		switch state[uid] {
		case TYPING_CHOICE:
			userData[uid]["choice"] = text
			state[uid] = TYPING_REPLY
			return c.Reply("Your " + text + "? Yes, I would love to hear about that!")

		case TYPING_REPLY:
			choice := userData[uid]["choice"]
			userData[uid][choice] = text
			delete(userData[uid], "choice")
			state[uid] = CHOOSING

			reply := "Neat! Just so you know, this is what you already told me:\n" +
				factsToStr(userData[uid]) +
				"\nYou can tell me more, or change your opinion on something."
			return c.Reply(reply, telebot.ReplyMarkup{ReplyKeyboard: replyKeyboard})

		default:
			return c.Reply("Please use the menu options.")
		}
	})

	log.Println("Bot is running...")
	bot.Start()
}

func choiceHandler(bot *telebot.Bot, choice string, userData map[int64]UserData, state map[int64]string) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		uid := c.Sender().ID
		state[uid] = TYPING_REPLY
		userData[uid]["choice"] = choice

		data := userData[uid]
		var reply string
		if value, ok := data[choice]; ok {
			reply = "Your " + choice + "? I already know the following about that: " + value
		} else {
			reply = "Your " + choice + "? Yes, I would love to hear about that!"
		}
		return c.Reply(reply)
	}
}

func keys(m map[string]string) []string {
	var k []string
	for key := range m {
		k = append(k, key)
	}
	return k
}
