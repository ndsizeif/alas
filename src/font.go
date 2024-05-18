package main

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// TODO font size

func ValidateFontProgram() error {
	_, err := exec.LookPath(fontListCmd)
	if err != nil {
		return err
	}
	return nil
}

func FontsInstalled() ([]string, error) {
	var fonts []string
	cmd := exec.Command(fontListCmd, fontListCompleteArgs...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	split := strings.Split(string(out), "\n")
	for _, v := range split {
		if v == "" {
			continue
		}
		fonts = append(fonts, v)
	}
	sort.Strings(fonts)
	fonts = RemoveDuplicate(fonts)
	return fonts, nil
}

// list fonts on system using custom fontListCmd (results depend on this)
func FontList(args []string) ([]string, error) {
	var fonts []string
	if len(args) == 0 {
		args = fontListDefaultArgs
	}

	cmd := exec.Command(fontListCmd, args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	split := strings.Split(string(out), "\n")
	for _, v := range split {
		if v == "" {
			continue
		}
		fonts = append(fonts, v)
	}
	sort.Strings(fonts)
	fonts = RemoveDuplicate(fonts)
	return fonts, nil
}

// compare name of font to entire list, error if no match is found
func FontMatch(list []string, name string) error {
	name = strings.ToLower(name)
	for _, font := range list {
		if strings.Compare(strings.ToLower(font), name) == 0 {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("font '%v' not found", name))
}

// return the current font configuration (return error if can't reach end of loop)
func GetCurrentFontBytes(paths []string) ([]byte, error) {
	var bytes []byte
	for _, path := range paths {
		_, err := os.Stat(filepath.Join(path, fontFile))
		if err != nil {
			continue
		}
		bytes, err = os.ReadFile(filepath.Join(path, fontFile))
		if err != nil {
			continue
		}
		fontFileExists = true
		// fmt.Println("font file found")
		return bytes, nil
	}
	return bytes, errors.New("font file not found")
}
// TODO separate fonts per type
// pass current font bytes and name of font to switch to
func FontValues(bytes []byte, font string) ([]byte, error) {
	var f FontStruct
	if fontFileExists {
		err := toml.Unmarshal(bytes, &f)
		if err != nil {
			return nil, err
		}
	}
	f.Font.Normal.Family = font
	f.Font.Bold.Family = font
	f.Font.Italic.Family = font
	f.Font.BoldItalic.Family = font

	bytes, err := toml.Marshal(f)
	bytes = ToLowerCase(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// used the passed font name to create and write font configuration to font file
func Font(paths []string, bytes []byte, name string) error {
	fontBytes, err := FontValues(bytes, name)
	if err != nil {
		return err
	}
	for _, path := range paths {
		err := WriteToFile(path, fontFile, fontBytes)
		if err != nil {
			continue
		}
	}
	return nil
}
