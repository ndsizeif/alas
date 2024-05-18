package main

import (
	"fmt"
	"os"
	"strings"
)

func Hex2Hash(hex string) string {
	return strings.Replace(hex, "0x", "#", -1)
}

func ExitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func RemoveDuplicate(strSlice []string) []string {
	var list []string
	keys := make(map[string]bool)

	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func PrintList(list []string) {
	for _, v := range list {
		fmt.Printf("%v\n", v)
	}
}

// make all strings lowercase, except "CellBackground", which can be used as a replacement for hex color values
func ToLowerCase(config []byte) []byte {
	s := strings.ToLower(string(config))
	s = strings.Replace(s, "cellbackground", "CellBackground", -1)
	s = strings.Replace(s, "cellforeground", "CellForeground", -1)
	return []byte(s)
}
