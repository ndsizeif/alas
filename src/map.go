package main

import "github.com/cmatthias/mapstructure"

type PrimaryColor struct {
	Background        string `toml:"background,omitempty"`
	Foreground        string `toml:"foreground,omitempty"`
	Dim_Foreground    string `toml:"dim_foreground,omitempty"`
	Bright_Foreground string `toml:"bright_foreground,omitempty"`
}
type NormalColor struct {
	Black   string `toml:"black,omitempty"`
	Red     string `toml:"red,omitempty"`
	Green   string `toml:"green,omitempty"`
	Yellow  string `toml:"yellow,omitempty"`
	Blue    string `toml:"blue,omitempty"`
	Magenta string `toml:"magenta,omitempty"`
	Cyan    string `toml:"cyan,omitempty"`
	White   string `toml:"white,omitempty"`
}
type BrightColor struct {
	Black   string `toml:"black,omitempty"`
	Red     string `toml:"red,omitempty"`
	Green   string `toml:"green,omitempty"`
	Yellow  string `toml:"yellow,omitempty"`
	Blue    string `toml:"blue,omitempty"`
	Magenta string `toml:"magenta,omitempty"`
	Cyan    string `toml:"cyan,omitempty"`
	White   string `toml:"white,omitempty"`
}
type DimColor struct {
	Black   string `toml:"black,omitempty"`
	Red     string `toml:"red,omitempty"`
	Green   string `toml:"green,omitempty"`
	Yellow  string `toml:"yellow,omitempty"`
	Blue    string `toml:"blue,omitempty"`
	Magenta string `toml:"magenta,omitempty"`
	Cyan    string `toml:"cyan,omitempty"`
	White   string `toml:"white,omitempty"`
}
type CursorColor struct {
	Cursor string `toml:"cursor,omitempty"`
	Text   string `toml:"text,omitempty"`
}
type ViCursorColor struct {
	Cursor string `toml:"cursor,omitempty"`
	Text   string `toml:"text,omitempty"`
}
type SearchColor struct {
	Matches struct {
		Background string `toml:"background,omitempty"`
		Foreground string `toml:"foreground,omitempty"`
	} `toml:"matches"`
	Focused struct {
		Background string `toml:"background,omitempty"`
		Foreground string `toml:"foreground,omitempty"`
	} `toml:"focused_match,omitempty"`
}
type HintsColor struct {
	Start struct {
		Background string `toml:"background,omitempty"`
		Foreground string `toml:"foreground,omitempty"`
	} `toml:"start"`
	End struct {
		Background string `toml:"background,omitempty"`
		Foreground string `toml:"foreground,omitempty"`
	} `toml:"end,omitempty"`
}
type LineColor struct {
	Background string `toml:"background,omitempty"`
	Foreground string `toml:"foreground,omitempty"`
}
type FooterColor struct {
	Background string `toml:"background,omitempty"`
	Foreground string `toml:"foreground,omitempty"`
}

type SelectionColor struct {
	Background string `toml:"Background,omitempty"`
	Text       string `toml:"text,omitempty"`
}

// Struct use to encode map[string]interface to scheme-compatible struct
// use color field when marshalling
type ColorScheme struct {
	Colors struct {
		Primary   PrimaryColor   `toml:"primary,omitempty"`
		Normal    NormalColor    `toml:"normal,omitempty"`
		Bright    BrightColor    `toml:"bright,omitempty"`
		Dim       DimColor       `toml:"dim,omitempty"`
		Cursor    CursorColor    `toml:"cursor,omitempty"`
		Vi_Mode_Cursor  ViCursorColor  `toml:"vi_mode_cursor,omitempty"`
		Search    SearchColor    `toml:"search,omitempty"`
		Hints     HintsColor     `toml:"hints,omitempty"`
		Line      LineColor      `toml:"line_indicator,omitempty"`
		Footer    FooterColor    `toml:"footer_bar,omitempty"`
		Selection SelectionColor `toml:"selection,omitempty"`
	} //`mapstructure:"colors"`
}

// create a ColorScheme struct from an interface
func BuildColorStruct(data map[string]interface{}) (ColorScheme, error) {
	var structure ColorScheme
	err := mapstructure.Decode(data, &structure)
	if err != nil {
		return structure, err
	}
	return structure, nil
}

/* // full Bell Structure for reference
type Bell struct {
	Animation string `toml:"animation, omitempty"`
	Command string `toml:"command, omitempty"`
	Color string `toml:"color, omitempty"`
	Duration int `toml:"duration, omitempty"`
}
*/

// we only care about the color
type BellStruct struct {
	Bell struct {
		Color string `toml:"color, omitempty"`
	} `toml:"bell, omitempty"`
}

type FontStruct struct {
	Font struct {
		Size float64 `toml:"size, omitempty"`
		BoxDrawing bool `toml:"builtin_box_drawing, omitempty"`
		Offset struct {
			x int `toml:"x,omitempty"`
			y int `toml:"y,omitempty"`
		} //`toml:"offset,omitempty"`
		Glyph_Offset struct {
			x int `toml:"x"`
			y int `toml:"y"`
		} `toml:"glyph_offset"`
		Normal struct {
			Family string `toml:"family, omitempty"`
			Style string `toml:"style, omitempty"`
		} `toml:"normal,omitempty"`
		Bold struct {
			Family string `toml:"family, omitempty"`
			Style string  `toml:"style, omitempty"`
		} `toml:"bold,omitempty"`
		Italic struct {
			Family string `toml:"family, omitempty"`
			Style string  `toml:"style, omitempty"`
		} `toml:"italic,omitempty"`
		BoldItalic struct {
			Family string `toml:"family, omitempty"`
			Style string  `toml:"style, omitempty"`
		} `toml:"bold_italic,omitempty"`
	} `toml:"font, omitempty"`
}
