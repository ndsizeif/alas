# alas

alas is a utility to set your [Alacritty](https://github.com/alacritty/alarcritty) appearance.

## Why

The Alacritty terminal emulator previous to `0.13.0` used `yaml` files for
configuration.  There are large collections of colorschemes and configuration
files available online.  I previously used a single file to store all of my
colorschemes and used a bash script to list colorschemes with simple regex
pattern matching. I changed the active colorscheme with search/replace.

The migration to `toml` configuration, while a definite improvement, left my
previous strategy for colorscheme selection defunct. The built-in Alacritty
migration tools, are effective at converting `yaml` to `toml`. But now my
colorscheme management tools need to marshal/unmarshal `toml`, deal with
duplicate keys, and covert table names.

While implementing this, I found there were other aspects that wanted to also
improve along the way; supporting all Alacritty directory locations, changing
the font, setting bell color, setting tty color strings, etc.  This utility
can perform several basic functions an Alacritty user might want.

## Why Not

Many colorschemes are distributed as individual files containing a single
colorscheme, such as `dracula.toml`. These files will use the table `[colors]`.
In most cases, it's much less painful to simply to change the import path to the
individual file in `alacritty.toml`. A shell script would be more than adequate.

## Installation

### Release

Navigate to the `Releases` section, download and run the latest binary for your
system architecture. Alternatively, build `alas` for your system from the source
tarball. Install [Go](https://go.dev/doc/install) if not already present on your system. 

### Clone

You may also clone the main branch of the project, and build the project that way. Place the resulting binary
in your `$PATH` or run locally.

<details>
    <summary>code statistics</summary>
```
===============================================================================
 Language            Files        Lines         Code     Comments       Blanks
===============================================================================
 Go                      8         1089          938           51          100
===============================================================================
 Total                   8         1089          938           51          100
===============================================================================
```
</details>

## Usage 

<details>
    <summary>alas --help</summary>
Usage: alas
  -l, -list
        return a list of available colorschemes
  -r, -random
        apply a random colorscheme
  -b, -bell
        set bell color to a base-8 color or colorscheme property
  -p, -print
        return string data from <colorscheme> (no input returns current)
  -f, -font
        apply <font>
  -F, -fonts
        return a list of available fonts
  -t, -tty
        convert colorscheme into sourceable shell script for tty colors
  -m, -mode
        return if colorscheme is a light or dark mode scheme
  -x, -clear
        clear current color settings and use default
  -h, -help
        print help for alas

Example: "alas <colorscheme>" to apply a colorscheme
</details>

### File Location

The program will read all `toml` files inside the subdirectory `/colorschemes` in
any valid Alacritty configuration directory. For example, these files will all be
read if an `alacritty.toml` file is present in the parent directory. If you have
existing colorschemes in `toml` files, they should be placed in a `/colorschemes`.

```
$HOME/.config/alacritty/colorschemes/myColors.toml
$HOME/.config/alacritty/colorschemes/oldColors.toml
$HOME/.config/colorschemes/otherColors.toml
$HOME/colorschemes/crazyColors.toml
```

### Alacritty.toml

The user must have an `alacritty.toml` file present in valid Alacritty
configuration directory.  This utility will not edit `alacritty.toml` directly.
Instead it will create or write to `colors.toml` and `fonts.toml` in the same
location. The `alacritty.toml` file must import these files for changes to take
effect.

Add the two files 
```toml
# alacritty.toml

import = [ 
	"~/.config/alacritty/colors.toml",
	"~/.config/alacritty/fonts.toml",
]
```

### Colorschemes
 
Each colorscheme should use its own unique name for the `toml` table. This is the
case when using the Alacritty migration tool on a yaml-based colorscheme. The
table `[colors]` should not be used in colorscheme files, as it is used by
Alacritty to set the active colorscheme. Doing so will return a "duplicate
table" error. 

<details>
    <summary>example colorscheme</summary>

```toml
[midboxlight.bright]
black = "0x928374"
blue = "0x076678"
cyan = "0x427b58"
green = "0x79740e"
magenta = "0x8f3f71"
red = "0x9d0006"
white = "0x3c3836"
yellow = "0xb57614"

[midboxlight.normal]
black = "0xE6D8AD"
blue = "0x458588"
cyan = "0x689d6a"
green = "0x98971a"
magenta = "0xb16286"
red = "0xcc241d"
white = "0x7c6f64"
yellow = "0xd79921"

[midboxlight.primary]
background = "0xfbf1c7"
foreground = "0x3c3836"
```

</details>

Alacritty's migration tool will not automatically handle duplicate colorscheme table names.
This can occur if a file, or multiple files have the same colorscheme, or at
least the same colorscheme name string when converted to `toml`.

<details>
    <summary>file tree</summary>

My `$HOME/.config/alacritty/`. I have two `toml` files that contain colorschemes
`custom.toml` for those I create and `internet.toml` for others that I find
posted elsewhere. You can have any number of `toml` files in that subdirectory.

```
├── alacritty.toml
├── colorschemes
│   ├── custom.toml
│   └── internet.toml
├── colors.toml
├── fonts.toml
├── keybindings.toml
└── ttyscheme.sh
```
</details>

## Tips

List available colorschemes.  
Pipe the colorscheme list into fzf and set the chosen scheme.  
Pipe the colorscheme list into fzf, set the chosen scheme, preview selection.  
```sh
alas -l
alas $(alas -l | fzf)
alas $(alas -l | fzf --preview='alas {}')
```

Pipe the font list into fzf and set the chosen font (requires fc-list).
```sh
alas -f $(alas -F | fzf)
```

Make Alacritty bell color red, or match it to another property.
```
alas -b red
alas -b cursor cursor
alas -b primary foreground 
```

Set your tty colors to match your alacritty theme.
```sh
alas -tty 
```

Source the created bash script in your `.bashrc` so tty colors match Alacritty.
```sh
if [ "$TERM" == "linux" ]; then
	if [ -f "$HOME/.config/alacritty/ttyscheme.sh" ]; then
		source "$HOME/.config/alacritty/ttyscheme.sh"
	fi
fi
```


## Contributing

Bug reports, or any form of constructive feedback is appreciated.

<details>
    <summary>here's a colorscheme for making all the way to the end!</summary>

```toml
[Plumbus.bright]
black = "#2E1A31"
blue = "#7269B8"
cyan = "#9062C4"
green = "#CD67C6"
magenta = "#A30061"
red = "#73002D"
white = "#E5C9E9"
yellow = "#A000BA"

[Plumbus.cursor]
cursor = "#736E7D"
text = "#050014"

[Plumbus.normal]
black = "#1F1720"
blue = "#77617B"
cyan = "#C1AEC4"
green = "#770E87"
magenta = "#9F82A3"
red = "#502659"
white = "#E7E0E8"
yellow = "#564559"

[Plumbus.primary]
background = "#130E14"
foreground = "#D2C5D4"
```

</details>
