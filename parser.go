package vcard

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"strings"
)

const (
	VERSION_21 = "2.1"
	VERSION_3  = "3.0"
	VERSION_4  = "4.0"
)

const (
	VCARD_BEGIN   = "BEGIN:VCARD"
	VCARD_END     = "END:VCARD"
	VCARD_VERSION = "VERSION:"
)

const (
	NEW_LINE      = "\r\n"
	FOLDING_SPACE = " "
	FOLDING_TAB   = "\t"
)

type VCardEntity struct {
	VCards []*VCard
}

type VCard struct {
	Version string
	Group   string
	N       N
	FN      string
}

func (card *VCard) Equal(otherCard *VCard) bool {
	return card.Group == otherCard.Group &&
		card.N.Equal(&otherCard.N) &&
		card.FN == otherCard.FN
}

type N struct {
	FamilyName        string
	GivenName         string
	AdditionalNames   string
	HonorificPrefixes string
	HonorificSuffixes string
}

func (n *N) Equal(otherN *N) bool {
	return n.FamilyName == otherN.FamilyName &&
		n.GivenName == otherN.GivenName &&
		n.AdditionalNames == otherN.AdditionalNames &&
		n.HonorificPrefixes == otherN.HonorificPrefixes &&
		n.HonorificSuffixes == otherN.HonorificSuffixes
}

func splitVCardContent(data []byte, atEOF bool) (advance int, token []byte, err error) {
	startIndex := 0

	for {

		newLineIndex := bytes.Index(data[startIndex:], []byte(NEW_LINE))

		if newLineIndex == -1 {
			// Read more data
			return 0, nil, nil
		}

		if startIndex+newLineIndex+3 > len(data) {
			// Cannot determine if folding happened -> need more data
			return 0, nil, nil
		}

		if bytes.Equal(data[startIndex+newLineIndex+2:startIndex+newLineIndex+3], []byte(FOLDING_SPACE)) == true ||
			bytes.Equal(data[startIndex+newLineIndex+2:startIndex+newLineIndex+3], []byte(FOLDING_TAB)) == true {
			// Folded, line actually goes on
			startIndex += newLineIndex + 3
			continue
		}

		clearTabs := bytes.ReplaceAll(data[:startIndex+newLineIndex], []byte("\r\n\t"), []byte(""))
		clearSpaces := bytes.ReplaceAll(clearTabs, []byte("\r\n "), []byte(""))

		// found a newline and it's not folded, so give out the line
		return startIndex + newLineIndex + 2, clearSpaces, nil
	}

}

func Parse(inputReader io.Reader) (*VCardEntity, error) {
	// make input buffered
	vCardScanner := bufio.NewScanner(inputReader)
	vCardScanner.Split(splitVCardContent)

	entity := &VCardEntity{}

	// This implements parsing of BEGIN, VERSION and END of VCard
	// Parsing content lines is done by the version appropriate parser
	for {
		if vCardScanner.Scan() == false {
			if err := vCardScanner.Err(); err != nil {
				return nil, err
			} else {
				return entity, nil
			}
		}

		line := vCardScanner.Text()

		if strings.ToUpper(line) != VCARD_BEGIN {
			return nil, errors.New("File is not formatted correctly!")
		}

		if vCardScanner.Scan() == false {
			if err := vCardScanner.Err(); err != nil {
				return nil, err
			} else {
				return entity, nil
			}
		}

		line = vCardScanner.Text()

		if strings.ToUpper(line) == "VERSION:"+VERSION_21 == true {
			newCard, err := vCard21Parse(vCardScanner)

			if err != nil {
				log.Println("Error while parsing input: ", err.Error())
				return entity, err
			}

			entity.VCards = append(entity.VCards, newCard)
		}

	}

}

type VCardEntityParser struct {
}
