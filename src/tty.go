package main

import (
	"errors"
	"fmt"
	"strings"
)

// ColorScheme struct colors in string slice
func ColorsToString(c ColorScheme) []string {
	black := c.Colors.Normal.Black
	red := c.Colors.Normal.Red
	green := c.Colors.Normal.Green
	yellow := c.Colors.Normal.Yellow
	blue := c.Colors.Normal.Blue
	magenta := c.Colors.Normal.Magenta
	cyan := c.Colors.Normal.Cyan
	white := c.Colors.Normal.White

	Black := c.Colors.Bright.Black
	Red := c.Colors.Bright.Red
	Green := c.Colors.Bright.Green
	Yellow := c.Colors.Bright.Yellow
	Blue := c.Colors.Bright.Blue
	Magenta := c.Colors.Bright.Magenta
	Cyan := c.Colors.Bright.Cyan
	White := c.Colors.Bright.White

	colors := []string{
		black, red, green, yellow, blue, magenta, cyan, white,
		Black, Red, Green, Yellow, Blue, Magenta, Cyan, White,
	}

	return colors
}

// transform string slice into a bash script that sets tty colors
func StringsToScript(colors []string) ([]byte, error) {
	var s []string
	var label string
	var prefix = []string{
		"P0", "P1", "P2", "P3", "P4", "P5", "P6", "P7",
		"P8", "P9", "PA", "PB", "PC", "PD", "PE", "PF",
	}

	s = append(s, fmt.Sprintf("#!/bin/bash\n")) // shebang

	for key, color := range colors {
		var entry string
		if color == "" && key < 8 { // cancel if normal colors are not present
			return nil, errors.New("color %v is invalid")
		}
		if color == "" && key > 7 { // if bright colors are not present, use normal
			label = strings.Replace(color, "#", prefix[key-8], 1)
			label = strings.Replace(color, "0x", prefix[key-8], 1)
		} else {
			label = strings.Replace(color, "#", prefix[key], 1)
			label = strings.Replace(color, "0x", prefix[key], 1)
		}
		entry = fmt.Sprintf("echo -en \"\\e]%v\"", label)

		s = append(s, entry)
	}

	s = append(s, fmt.Sprintf("clear\n")) // remove artifacts
	script := strings.Join(s, "\n")
	return []byte(script), nil
}

// write tty script to file
func SaveTTYScript(paths []string, script []byte) error {
	if script == nil {
		return errors.New("empty tty file")
	}

	for _, path := range paths {
		err := WriteToFile(path, ttyFile, script)
		if err != nil {
			continue
		}
	}

	return nil
}
