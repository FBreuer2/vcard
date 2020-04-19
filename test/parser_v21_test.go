package parser_21_test

import (
	"strings"
	"testing"

	"github.com/FBreuer2/vcard"
)

var wellformedVCARD = []struct {
	entityString string
	entity       vcard.VCardEntity
}{
	{"BEGIN:VCARD\r\nVERSION:2.1\r\nFN:Florian\r\n\t Breuer\r\nEND:VCARD\r\n", vcard.VCardEntity{
		VCards: []*vcard.VCard{
			{FN: "Florian Breuer", Version: "2.1"},
		},
	}},
	{"BEGIN:VCARD\r\nVERSION:2.1\r\nFN:Florian\r\n\t Bre\r\n\tue\r\n\tr\r\nEND:VCARD\r\n", vcard.VCardEntity{
		VCards: []*vcard.VCard{
			{FN: "Florian Breuer", Version: "2.1"},
		},
	}},
	{"BEGIN:VCARD\r\nVERSION:2.1\r\nN:Breuer;Florian;;;\r\nFN:Florian Breuer\r\nEND:VCARD\r\n", vcard.VCardEntity{
		VCards: []*vcard.VCard{
			{FN: "Florian Breuer", Version: "2.1", N: vcard.N{FamilyName: "Breuer", GivenName: "Florian"}},
		},
	}},
	{"BEGIN:VCARD\r\nVERSION:2.1\r\nN;CHARSET=UTF-8;ENCODING=QUOTED-PRINTABLE:=47=72=65=67;=47=72=65=67;=47=72=65=67;=47=72=65=67;=47=72=65=67\r\nFN;CHARSET=UTF-8;ENCODING=QUOTED-PRINTABLE:=47=72=65=67=47=72=65=67\r\nEND:VCARD\r\n", vcard.VCardEntity{
		VCards: []*vcard.VCard{
			{FN: "GregGreg", Version: "2.1", N: vcard.N{FamilyName: "Greg", GivenName: "Greg", AdditionalNames: "Greg", HonorificPrefixes: "Greg", HonorificSuffixes: "Greg"}},
		},
	}},
}

func TestWellformedVCard(t *testing.T) {
	for _, instance := range wellformedVCARD {
		reader := strings.NewReader(instance.entityString)
		entity, err := vcard.Parse(reader)

		if err != nil {
			t.Errorf("Parsing a wellformed VCard created an error: %s \n VCard: %s", err.Error(), instance.entityString)
		}

		for index, vcard := range entity.VCards {
			if vcard.Equal(instance.entity.VCards[index]) == false {
				t.Errorf("Parsing a wellformed VCard didn't create expected VCard structure: VCard: %s", instance.entityString)
			}
		}
	}
}
