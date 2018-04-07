package bencode

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// bencoding encoding function
// Beispiel:
// map[ein test:foo zahl:5 liste:[string in liste 2 3]]
// wird zu
// d5:listel15:string in listei2ei3ee8:ein test3:foo4:zahli5ee
func Encode(intf interface{}) []byte {
	t := reflect.TypeOf(intf).String()
	switch t {
	case "int":
		return []byte("i" + strconv.Itoa(intf.(int)) + "e")
	case "string":
		return []byte(strconv.Itoa(len(intf.(string))) + ":" + intf.(string))
	case "[]interface {}":
		var text string
		for _, v := range intf.([]interface{}) {
			text += string(Encode(v))
		}
		return []byte("l" + text + "e")
	case "map[string]interface {}":
		var text string
		for k, v := range intf.(map[string]interface{}) {
			text += string(Encode(k)) + string(Encode(v))
		}
		return []byte("d" + text + "e")
	}
	return []byte("")
}

// bencoding decoding function
// Beispiel:
// d5:listel15:string in listei2ei3ee8:ein test3:foo4:zahli5ee
// wird zu
// map[ein test:foo zahl:5 liste:[string in liste 2 3]]
func Decode(s []byte) interface{} {
	text := string(s)
	k, rest := decodeFirst(text)

	// wenn ein Rest vorhanden ist, heißt das,
	// das
	if rest != "" {
		panic("Input nicht gut formatiert")
	}
	return k
}

// eigentliche decoding Funktion
// schaut sich das erste Zeichen an und bestimmt anhand dessen,
// was für ein Datenentyp es kodieren soll.
// gibt das decodete Ergebnis aus und den restlichen String,
// der noch dekodiert werden muss
func decodeFirst(s string) (interface{}, string) {
	if strings.HasPrefix(s, "i") {
		esplit := strings.SplitN(s[1:], "e", 2)
		zahl, _ := strconv.Atoi(esplit[0])
		return zahl, esplit[1]
	}

	if unicode.IsDigit(rune(s[0])) {
		esplit := strings.SplitN(s, ":", 2)
		length, _ := strconv.Atoi(esplit[0])
		return esplit[1][:length], esplit[1][length:]
	}

	if strings.HasPrefix(s, "l") || strings.HasPrefix(s, "d") {
		var liste []interface{}
		lmap := make(map[string]interface{})
		var rest = s[1:]
		var ergebnis interface{}
		for !strings.HasPrefix(rest, "e") {
			ergebnis, rest = decodeFirst(rest)
			liste = append(liste, ergebnis)
		}
		rest = rest[1:]

		if strings.HasPrefix(s, "l") {
			return liste, rest
		} else {
			for i := 0; i < len(liste); i += 2 {
				lmap[liste[i].(string)] = liste[i+1]
			}
			return lmap, rest
		}
	}

	return "foo", "bar"
}
