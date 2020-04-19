package main

import (
	"fmt"
	"strings"

	"github.com/FBreuer2/vcard"
)

func main() {
	cardString := "BEGIN:VCARD\r\nVERSION:2.1\r\nN:Breuer;Florian;;;\r\nFN:Florian Breuer\r\nEND:VCARD\r\n"

	reader := strings.NewReader(cardString)
	entity, err := vcard.Parse(reader)

	if err != nil {
		// Parsing did not work
		fmt.Println("Error while parsing VCards: ", err.Error())
		return
	}

	// Do something with the parsed information
	fmt.Println("Parsed ", len(entity.VCards), " VCards.")

	return
}
