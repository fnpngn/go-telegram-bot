package main

import (
	"reflect"
	"testing"
)

func TestFactsToStr(t *testing.T) {
	data := UserData{
		"Age":                "25",
		"Favourite colour":   "blue",
		"Number of siblings": "3",
	}

	expected := "Age - 25\nFavourite colour - blue\nNumber of siblings - 3"
	result := factsToStr(data)

	if result != expected {
		t.Errorf("factsToStr() = %q; want %q", result, expected)
	}
}

func TestHandleChoice(t *testing.T) {
	mockStates := make(map[int64]string)
	mockData := make(map[int64]UserData)
	mockChatID := int64(1234)

	mockData[mockChatID] = make(UserData)

	choice := "Age"
	handleChoice(mockChatID, mockStates, mockData, choice)

	if mockStates[mockChatID] != TYPING_REPLY {
		t.Errorf("State = %q; want %q", mockStates[mockChatID], TYPING_REPLY)
	}
	if mockData[mockChatID]["choice"] != choice {
		t.Errorf("Choice = %q; want %q", mockData[mockChatID]["choice"], choice)
	}
}

func TestMapKeys(t *testing.T) {
	data := map[string]string{
		"Age":                "25",
		"Favourite colour":   "blue",
		"Number of siblings": "3",
	}

	result := mapKeys(data)
	expected := []string{"Age", "Favourite colour", "Number of siblings"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("mapKeys() = %v; want %v", result, expected)
	}
}
