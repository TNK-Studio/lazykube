package utils

import (
	"github.com/spkg/bom"
	"regexp"
	"sort"
	"strings"
)

func CleanString(s string) string {
	output := string(bom.Clean([]byte(s)))
	return NormalizeLinefeeds(output)
}

// NormalizeLinefeeds - Removes all Windows and Mac style line feeds
func NormalizeLinefeeds(str string) string {
	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

func OptionsMapToString(optionsMap map[string]string) string {
	optionsArray := make([]string, 0)
	for key, description := range optionsMap {
		optionsArray = append(optionsArray, key+": "+description)
	}
	sort.Strings(optionsArray)
	return strings.Join(optionsArray, ", ")
}

func DeleteExtraSpace(s string) string {
	s1 := strings.Replace(s, "	", " ", -1)
	regx := "\\s{2,}"
	reg, _ := regexp.Compile(regx)
	s2 := make([]byte, len(s1))
	copy(s2, s1)
	spcIndex := reg.FindStringIndex(string(s2))
	for len(spcIndex) > 0 {
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...)
		spcIndex = reg.FindStringIndex(string(s2))
	}
	return string(s2)
}
