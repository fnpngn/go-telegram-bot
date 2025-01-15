package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/mock"
)

func TestStartCommand(t *testing.T) {
	// Create a mock bot
	mockBot := mock.New()
	defer mockBot.Finish()

	userData := make(map[int64]UserData)
	state := make(map[int64]string)

	// Simulate the /start command handler
	mockBot.On("/start", func(c telebot.Context) error {
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

	// Test /start command
	mockBot.Send(&telebot.Message{Text: "/start"})
	reply := mockBot.LastMessage()
	assert.Contains(t, reply.Text, "Hi! My name is Doctor Botter.")
	assert.Contains(t, reply.Text, "Why don't you tell me something about yourself?")
}

func TestChoiceHandler(t *testing.T) {
	// Create a mock bot
	mockBot := mock.New()
	defer mockBot.Finish()

	userData := make(map[int64]UserData)
	state := make(map[int64]string)

	// Simulate the choice handler
	mockBot.On("Age", func(c telebot.Context) error {
		uid := c.Sender().ID
		state[uid] = TYPING_REPLY
		userData[uid]["choice"] = "Age"

		data := userData[uid]
		var reply string
		if value, ok := data["Age"]; ok {
			reply = "Your Age? I already know the following about that: " + value
		} else {
			reply = "Your Age? Yes, I would love to hear about that!"
		}
		return c.Reply(reply)
	})

	// Test Age button
	mockBot.Send(&telebot.Message{Text: "Age"})
	reply := mockBot.LastMessage()
	assert.Contains(t, reply.Text, "Your Age? Yes, I would love to hear about that!")
}

func TestDoneHandler(t *testing.T) {
	// Create a mock bot
	mockBot := mock.New()
	defer mockBot.Finish()

	userData := make(map[int64]UserData)
	state := make(map[int64]string)

	// Simulate the done handler
	mockBot.On("Done", func(c telebot.Context) error {
		uid := c.Sender().ID
		data := userData[uid]
		reply := "I learned these facts about you: " + factsToStr(data) + "\nUntil next time!"
		delete(state, uid)
		return c.Reply(reply, &telebot.ReplyMarkup{RemoveKeyboard: true})
	})

	// Test Done button
	mockBot.Send(&telebot.Message{Text: "Done"})
	reply := mockBot.LastMessage()
	assert.Contains(t, reply.Text, "I learned these facts about you:")
	assert.Contains(t, reply.Text, "Until next time!")
}
