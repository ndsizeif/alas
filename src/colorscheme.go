package main

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml"
	"sort"
	"strings"
	"strconv"
	"os"
)

func GetMapInterface(data []byte) (map[string]interface{}, error) {
	var tree map[string]interface{}
	err := toml.Unmarshal(data, &tree)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// return list of colorscheme names/keys from a slice of bytes
func ColorschemeNames(bytes []byte) ([]string, error) {
	data, err := GetMapInterface(bytes)
	if err != nil {
		return nil, err
	}
	var names []string
	for key := range data {
		names = append(names, key)
	}
	sort.Strings(names)
	return names, nil
}

// set bell color and return structured toml bytes
func BellMarshal(structure ColorScheme) ([]byte, error) {
	var color string
	var b BellStruct
	color = structure.Colors.Bright.Yellow
	b.Bell.Color = color
	data, err := toml.Marshal(b)
	if err != nil {
		return data, err
	}
	return data, nil
}

// use type assertions to extract the color category and then the color value
func BellMatchColor(bytes []byte, input []string) (string, error) {
	var bellColor string
	if len(input) < 2 {
		return bellColor, errors.New("bell input requires type and color")
	}

	data, err := GetMapInterface(bytes)
	if err != nil {
		return bellColor, err
	}

	colorTable := map[string]interface{}{"colors": nil}
	colorTable["colors"] = data["colors"]

	attr := string(ToLowerCase([]byte(input[0])))
	value := string(ToLowerCase([]byte(input[1])))

	// will crash even with type assertion if "colors" is upper "Colors
	colorType, ok := colorTable["colors"].(map[string]interface{})[attr]
	if !ok {
		return bellColor, errors.New(fmt.Sprintf("no color type %v", attr))
	}

	colorValue := colorType.(map[string]interface{})[value]
	bellColor, ok = colorValue.(string)
	if !ok {
		return bellColor, errors.New(fmt.Sprintf("color: %v not found", value))
	}

	return bellColor, nil
}

func BellReplace(paths []string, colorscheme []byte, color string) error {
	// place the chosen color value in a bell struct
	var b BellStruct
	b.Bell.Color = color
	bellBytes, err := toml.Marshal(&b)
	if err != nil {
		return err
	}
	// get the current colorscheme minus bell struct
	var current ColorScheme
	err = toml.Unmarshal(colorscheme, &current)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// convert the colorscheme to bytes
	bytes, err := toml.Marshal(current)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// add bell struct bytes to colorscheme
	bytes = append(bytes, bellBytes...)
	bytes = ToLowerCase(bytes)
	// write to file(s)
	for _, path := range paths {
		err := WriteToFile(path, colorFile, bytes)
		if err != nil {
			continue
		}
	}
	return nil
}

// return the colorscheme's color table values
func ColorschemeValues(bytes []byte, input string) ([]byte, error) {
	if input == "" {
		return nil, errors.New("no valid colorscheme")
	}

	data, err := GetMapInterface(bytes)
	if err != nil {
		return nil, err
	}

	dataMap := data[input]
	table := map[string]interface{}{"colors": nil}
	table["colors"] = dataMap

	dataStructure, err := BuildColorStruct(table)
	if err != nil {
		return nil, err
	}
	dataBytes, err := ColorschemeMarshal(dataStructure)
	if err != nil {
		return nil, err
	}
	bell, err := BellMarshal(dataStructure)
	if err != nil {
		return nil, err
	}
	dataBytes = append(dataBytes, bell...)
	dataBytes = ToLowerCase(dataBytes)
	return dataBytes, nil
}

// Colorscheme struct -> byte slice
func ColorschemeMarshal(structure ColorScheme) ([]byte, error) {
	data, err := toml.Marshal(structure)
	if err != nil {
		return data, err
	}
	return data, nil
}

// apply the named colorscheme if included in the provided slice of bytes
// writes to each qualifying alacritty path (contains alacritty.toml)
func Colorscheme(paths []string, bytes []byte, name string) error {
	// fmt.Println("colorscheme")
	colorscheme, err := ColorschemeValues(bytes, name)
	if err != nil {
		// fmt.Print(err)
		return err
	}
	// fmt.Println("write colorscheme")
	for _, path := range paths {
		err := WriteToFile(path, colorFile, colorscheme)
		if err != nil {
			continue
		}
	}
	return nil
}

// write 0 bytes to color file
func ColorschemeClear(paths []string) error {
	var empty []byte
	for _, path := range paths {
		err := WriteToFile(path, colorFile, empty)
		if err != nil {
			continue
		}
	}
	return nil
}

func ColorschemeIsDark(colorscheme []byte) (bool, error) {
	var current ColorScheme
	var dark bool
	var base = 16
	var bitsize = 32 // int not uint
	err := toml.Unmarshal(colorscheme, &current)
	if err != nil {
		return dark, err
	}
	background := TrimCode(current.Colors.Primary.Background)
	if len(background) != 6 {
		return dark, errors.New("not a valid color string for background")
	}
	r, err := strconv.ParseInt(background[0:2], base, bitsize)
	if err != nil {
		return dark, err
	}
	g, err := strconv.ParseInt(background[2:4], base, bitsize)
	if err != nil {
		return dark, err
	}
	b, err := strconv.ParseInt(background[4:6], base, bitsize)
	if err != nil {
		return dark, err
	}
	l := Luminosity(r, g, b)
	dark = l < 0.5
	return dark, err
}

func TrimCode(color string) string{
	if len(color) == 7 {
		if strings.HasPrefix(color, "#") {
			return color[1:]
		}
	}
	if len(color) == 8 {
		if strings.HasPrefix(color, "0x") {
			return color[2:]
		}
	}
	return color
}

func Luminosity(r, g, b int64) float64 {
	return (float64(0.2126)*float64(r) +
		(float64(0.7152) * float64(g)) +
		(float64(0.0722) * float64(b)))
}

// use to set env variable for this process or child processes
func SetColorschemeMode(value string) error {
	if value == "" {
		return errors.New("no mode value set")
	}
	err := os.Setenv(modeEnvVariable, value)
	if err != nil {
		return err
	}
	return nil
}

func GetCurrentStruct(colorscheme []byte) (ColorScheme, error) {
	var current ColorScheme
	err := toml.Unmarshal(colorscheme, &current)
	if err != nil {
		return current, err
	}
	return current, nil
}
