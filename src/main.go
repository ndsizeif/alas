package main

func main() {
	paths, err := GetConfigPaths()
	ExitOnError(err)

	alacrittyPaths, err := GetAlacrittyPaths(paths)
	ExitOnError(err)

	bytes, err := GetColorschemeBytes(alacrittyPaths)
	ExitOnError(err)

	currentColorschemeBytes, err = GetCurrentColorschemeBytes(alacrittyPaths)
	currentFontBytes, err = GetCurrentFontBytes(alacrittyPaths)

	input := UserInput()
	HandleFlags(bytes, alacrittyPaths, input)
}
