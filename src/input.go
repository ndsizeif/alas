package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

var (
	listSchemes  bool
	printData    bool
	randomScheme bool
	selectFont   bool
	listFonts    bool
	bell         bool
	help         bool
	clearScheme  bool
	ttyScheme    bool
	mode         bool
)

func UserInput() []string {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s\n", programName)

		order := []string{
			"list", "random", "bell", "print", "font", "fonts", "tty", "mode", "clear", "help"}

		for _, name := range order {
			f := flag.Lookup(name)
			fmt.Printf("  -%v, -%v\n", Shorthand[f.Name], f.Name)
			fmt.Printf("\t%s\n", UsageString[f.Name])
		}
		fmt.Printf("\n")
		fmt.Printf("Example: \"%v <colorscheme>\" to apply a color scheme\n", programName)
		fmt.Printf("\n")
	}

	flag.BoolVar(&help, "help", false, UsageString["help"])
	flag.BoolVar(&help, "h", false, UsageString["help"])

	flag.BoolVar(&clearScheme, "clear", false, UsageString["clear"])
	flag.BoolVar(&clearScheme, "x", false, UsageString["clear"])

	flag.BoolVar(&ttyScheme, "tty", false, UsageString["tty"])
	flag.BoolVar(&ttyScheme, "t", false, UsageString["tty"])

	flag.BoolVar(&mode, "mode", false, UsageString["mode"])
	flag.BoolVar(&mode, "m", false, UsageString["mode"])

	flag.BoolVar(&listSchemes, "list", false, UsageString["list"])
	flag.BoolVar(&listSchemes, "l", false, UsageString["list"])

	flag.BoolVar(&printData, "print", false, UsageString["print"])
	flag.BoolVar(&printData, "p", false, UsageString["print"])

	flag.BoolVar(&randomScheme, "random", false, UsageString["random"])
	flag.BoolVar(&randomScheme, "r", false, UsageString["random"])

	flag.BoolVar(&listFonts, "fonts", false, UsageString["fonts"])
	flag.BoolVar(&listFonts, "F", false, UsageString["fonts"])

	flag.BoolVar(&selectFont, "font", false, UsageString["font"])
	flag.BoolVar(&selectFont, "f", false, UsageString["font"])

	flag.BoolVar(&bell, "bell", false, UsageString["bell"])
	flag.BoolVar(&bell, "b", false, UsageString["bell"])

	flag.Parse()
	cmd := flag.Args()

	if help {
		flag.Usage()
		os.Exit(0)
	}
	return cmd
}

func ValidateInput(input []string) error {
	if len(input) == 0 {
		return errors.New("error: no user input")
	}
	return nil
}

func HandleClear(paths []string) {
	err := ColorschemeClear(paths)
	ExitOnError(err)
	os.Exit(0)
}
func HandleListFonts() {
	args := flag.Args()
	list, err := FontList(args)
	ExitOnError(err)
	PrintList(list)
	os.Exit(0)
}
func HandleListColors(bytes []byte) {
	colors, err := ColorschemeNames(bytes)
	ExitOnError(err)
	PrintList(colors)
	os.Exit(0)
}
func HandlePrintData(bytes []byte, input []string) {
	err := ValidateInput(input) // print current color data if no input
	if err != nil {
		fmt.Printf("%v", string(currentColorschemeBytes))
		os.Exit(0)
	}
	name := strings.Join(input[0:], " ")
	data, err := ColorschemeValues(bytes, name)
	ExitOnError(err)
	fmt.Printf("%v", string(data))
	os.Exit(0)
}
func HandleRandom(bytes []byte, paths []string) {
	colors, err := ColorschemeNames(bytes)
	ExitOnError(err)
	name := colors[rand.Intn(len(colors))]
	Colorscheme(paths, bytes, name)
	os.Exit(0)
}
func HandleSelectFont(paths, input []string) {
	fonts, err := FontsInstalled()
	if err != nil {
		fmt.Println(err)
		return
	}
	ExitOnError(err)
	err = ValidateInput(input)
	ExitOnError(err)
	// FIXME not guaranteed to be input[0] when multiple flags are called
	font := strings.Join(input[0:], " ") // handles multi-word strings with/without quotes
	err = FontMatch(fonts, font)

	ExitOnError(err)
	err = Font(paths, currentFontBytes, font)
	ExitOnError(err)
	os.Exit(0)
}
func HandleBellColor(bytes []byte, paths, input []string) {
	err := ValidateInput(input)
	ExitOnError(err)
	if len(input) == 1 { // handle single color word by prefixing normal
		input = append(input, input[0])
		input[0] = "normal"
	}

	bellColor, err := BellMatchColor(currentColorschemeBytes, input)
	ExitOnError(err)
	BellReplace(paths, currentColorschemeBytes, bellColor)
	os.Exit(0)
}
func HandleTTY(paths []string) {
	data, err := GetCurrentStruct(currentColorschemeBytes)
	ExitOnError(err)
	colors := ColorsToString(data)
	ttycolors, err := StringsToScript(colors)
	ExitOnError(err)
	SaveTTYScript(paths, ttycolors)
}
func HandleMode() {
	dark, err := ColorschemeIsDark(currentColorschemeBytes)
	ExitOnError(err)
	if dark {
		fmt.Println("dark")
		return
	}
	fmt.Println("light")
}
// default case of single string input
func HandleDefaultInput(bytes []byte, paths, input []string) error {
	err := ValidateInput(input)
	if err != nil { // continue program if no input color
		return err
	}
	name := input[0] // colorscheme string should not contain a space character
	Colorscheme(paths, bytes, name)
	os.Exit(0)
	return nil
}
// pass bytes from colorscheme & files, user input
func HandleFlags(bytes []byte, paths, input []string) {
	if mode {
		HandleMode()
		os.Exit(0)
	}
	if ttyScheme {
		HandleTTY(paths)
		os.Exit(0)
	}
	if listFonts { // XXX
		err := ValidateFontProgram()
		ExitOnError(err)
		HandleListFonts()
		os.Exit(0)
	}
	if listSchemes {
		HandleListColors(bytes)
	}
	if printData {
		HandlePrintData(bytes, input)
	}
	if randomScheme {
		HandleRandom(bytes, paths)
	}
	if bell {
		HandleBellColor(bytes, paths, input)
	}
	if selectFont {
		HandleSelectFont(paths, input)
	}
	if clearScheme {
		HandleClear(paths)
	}
	err := HandleDefaultInput(bytes, paths, input)
	if err != nil {
		flag.Usage() // print help if no input
		return
	}
}

var UsageString = map[string]string{
	"list":   "return a list of available color schemes",
	"select": "apply <colorscheme>",
	"random": "apply a random color scheme",
	"bell":   "set bell color to a base-8 color or color scheme property",
	"print":  "return string data from <colorscheme> (no input returns current)",
	"fonts":  "return a list of available fonts",
	"font":   "apply <font>",
	"help":   "print help for " + programName,
	"clear":  "clear current color settings and use default",
	"tty":    "convert color scheme into sourceable shell script for tty colors",
	"mode":   "return if color scheme is a light or dark mode scheme",
}
var Shorthand = map[string]string{
	"list":   "l",
	"select": "s",
	"random": "r",
	"bell":   "b",
	"print":  "p",
	"font":   "f",
	"fonts":  "F",
	"help":   "h",
	"clear":  "x",
	"tty":    "t",
	"mode":   "m",
}
