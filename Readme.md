# Introduction
This golang package offers parsing and editing for VCards (a common format for saving contact information). 
It aims to give a unified interface for the different VCard versions so you don't have to think about them.

The library aims to be
* easy to use
* performant
* correct (verified by an extensive test suite)


It aims to support as many versions of VCards as possible. Here is the current progress:
- [x] VCard 2.1 (beta)
- [ ] VCard 3
- [ ] VCard 4

More detailed progress can be seen in [in our docs](doc/progress.md)

# Code examples
Here is a simple example of how to use the code:

```Go
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
```