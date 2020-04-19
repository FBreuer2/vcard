package vcard

import (
	"bufio"
	"errors"
	"io"
	"mime/quotedprintable"
	"strings"
)

func vCard21Parse(scanner *bufio.Scanner) (*VCard, error) {
	reachedEnd := false

	currentVCard := &VCard{
		Version: VERSION_21,
	}

	for reachedEnd != true {
		if scanner.Scan() == false {
			if err := scanner.Err(); err != nil {
				return nil, err
			} else {
				return currentVCard, nil
			}
		}

		line := scanner.Text()

		if line == VCARD_END {
			return currentVCard, nil
		}

		err := handleLine(line, currentVCard)

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func handleLine(line string, currentCard *VCard) error {
	if strings.HasPrefix(line, "TEL") {
		return handleTEL(line, currentCard)
	}

	if strings.HasPrefix(line, "FN") {
		return handleFN(line, currentCard)
	}

	if strings.HasPrefix(line, "N") {
		return handleN(line, currentCard)
	}

	if strings.HasPrefix(line, "KIND") {
		return handleKIND(line, currentCard)
	}

	return nil
}

func handleTEL(line string, currentCard *VCard) error {
	// contentline  = TEL *(";" param) ":" value CRLF
	splitLine := strings.Split(line, ":")
	if len(splitLine) != 2 {
		return errors.New("TEL line malformed.")
	}

	tel := TEL{
		Number: splitLine[1],
	}

	// TEL;PREF;WORK;MSG;FAX
	attributes := strings.Split(splitLine[0], ";")[1:]

	attributeString := ""
	for _, attr := range attributes {
		attributeString += attr + ","
	}

	tel.Attributes = strings.TrimSuffix(attributeString, ",")

	currentCard.Numbers = append(currentCard.Numbers, tel)
	return nil
}

func handleFN(line string, currentCard *VCard) error {
	// contentline  = FN *(";" param) ":" value CRLF
	splitLine := strings.Split(line, ":")

	if len(splitLine) != 2 {
		return errors.New("FN line malformed.")
	}

	// get Params
	params := getParams(strings.Split(splitLine[0], ";")[1:])

	if params["ENCODING"] == "" {
		currentCard.FN = strings.ReplaceAll(splitLine[1], "\\,", ",")
	} else if params["ENCODING"] == "QUOTED-PRINTABLE" {
		decodedString, err := getDecodedQuotedPrintable(splitLine[1])
		if err != nil {
			return err
		}

		currentCard.FN = strings.ReplaceAll(decodedString, "\\,", ",")
	}

	return nil
}

func handleN(line string, currentCard *VCard) error {
	// contentline  = N *(";" param) ":" value CRLF
	splitLine := strings.Split(line, ":")

	// get Params
	params := getParams(strings.Split(splitLine[0], ";")[1:])

	values := strings.Split(splitLine[1], ";")
	if len(values) != 5 {
		return errors.New("N line malformed!")
	}

	if params["ENCODING"] == "" {
		currentCard.N = N{
			FamilyName:        values[0],
			GivenName:         values[1],
			AdditionalNames:   values[2],
			HonorificPrefixes: values[3],
			HonorificSuffixes: values[4],
		}
		return nil
	} else if params["ENCODING"] == "QUOTED-PRINTABLE" {
		n := N{}

		decodedString := make([]string, 5)

		for index, value := range values {
			decodedValue, err := getDecodedQuotedPrintable(value)
			decodedString[index] = decodedValue
			if err != nil {
				return err
			}
		}

		n.FamilyName = decodedString[0]
		n.GivenName = decodedString[1]
		n.AdditionalNames = decodedString[2]
		n.HonorificPrefixes = decodedString[3]
		n.HonorificSuffixes = decodedString[4]

		currentCard.N = n

		return nil
	} else {
		return errors.New("Unkown encoding in N field")
	}
}

func handleKIND(line string, currentCard *VCard) error {
	return nil
}

func getParams(params []string) map[string]string {
	parsedParams := make(map[string]string)

	for _, param := range params {
		paramContent := strings.Split(param, "=")
		parsedParams[paramContent[0]] = paramContent[1]
	}

	return parsedParams
}

func getDecodedQuotedPrintable(encoded string) (string, error) {
	var decodedString = ""
	decodingReader := quotedprintable.NewReader(strings.NewReader(encoded))
	buf := make([]byte, 200)

	for {
		readBytes, err := decodingReader.Read(buf)

		decodedString += string(buf[:readBytes])
		if err == io.EOF {
			return decodedString, nil
		} else if err != nil {
			return "", err
		}
	}
}
