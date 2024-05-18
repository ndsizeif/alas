package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
)

const (
	programName     = "alas"
	alacritty       = "alacritty"
	alacrittyConfig = "alacritty.toml"
	colorFile       = "colors.toml"
	fontFile        = "fonts.toml"
	colorschemeDir  = "colorschemes"
	ttyFile         = "ttyscheme.sh"
	modeEnvVariable = "ALAS_MODE"
	fontListCmd     = "fc-list"
)

var (
	colorFileExists         bool                                                                       // assume the file does not exist until read without error
	fontFileExists          bool                                                                       // assume the file does not exist until read without error
	fontListDefaultArgs     = []string{":lang=en:style=Regular:scalable=True", "--format=%{family}\n"} // subset of all installed fonts
	fontListCompleteArgs    = []string{"--format=%{family}\n"}                                         // match against all installed fonts
	currentColorschemeBytes []byte                                                                     // hold current (previous) colorscheme bytes
	currentFontBytes        []byte                                                                     // hold current (previous) font bytes
)

// get list of possible configuration file paths
func GetConfigPaths() ([]string, error) {
	var paths []string
	user, err := user.Current()
	if err == nil {
		paths = append(paths, user.HomeDir)
	}
	config, err := os.UserConfigDir()
	if err == nil {
		paths = append(paths, config)
		paths = append(paths, filepath.Join(config, alacritty))
	}
	if len(paths) == 0 {
		err = errors.New("No valid configuration paths found")
		return paths, err
	}
	return paths, nil
}

// get list of configuration file paths that have alacritty.toml
func GetAlacrittyPaths(root []string) ([]string, error) {
	var paths []string
	for _, path := range root {
		_, err := os.Stat(filepath.Join(path, alacrittyConfig))
		if err == nil {
			paths = append(paths, path)
		}
	}
	if len(paths) == 0 {
		err := errors.New("alacritty.toml not found")
		return paths, err
	}
	return paths, nil
}

// return list of non-directory file info objects for a given directory
func GetDirectoryFiles(dir string) ([]fs.FileInfo, error) {
	file, err := os.Open(dir)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	files, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}
	var info []fs.FileInfo
	for _, v := range files {
		if !v.IsDir() {
			info = append(info, v)
		}
	}
	return info, nil
}

// for each potential path (with an alacritty.toml), find the colorscheme directory and extract the bytes from each file therein
func GetColorschemeBytes(paths []string) ([]byte, error) {
	var bytes []byte
	for _, path := range paths {
		dir := filepath.Join(path, colorschemeDir)
		info, err := GetDirectoryFiles(dir)
		if err != nil {
			continue
		}
		for _, file := range info {
			fileBytes, err := os.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				continue
			}
			bytes = append(bytes, fileBytes...)
		}
	}
	return bytes, nil
}

// return the bytes of the current alacritty colorscheme (return error if can't reach end of loop)
func GetCurrentColorschemeBytes(paths []string) ([]byte, error) {
	var bytes []byte
	for _, path := range paths {
		_, err := os.Stat(filepath.Join(path, colorFile))
		if err != nil {
			continue
		}
		bytes, err = os.ReadFile(filepath.Join(path, colorFile))
		if err != nil {
			continue
		}
		colorFileExists = true
		return bytes, nil
	}
	return bytes, errors.New("color file not found")
}

// (over)write bytes to colorFile file in given directory
func WriteToFile(path, basename string, data []byte) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}
	filename := filepath.Join(path, basename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
