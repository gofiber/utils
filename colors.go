package utils

import (
	"fmt"
	"io"

	colorable "github.com/mattn/go-colorable"
)

const (
	colorBlack   = "\u001b[90m"
	colorRed     = "\u001b[91m"
	colorGreen   = "\u001b[92m"
	colorYellow  = "\u001b[93m"
	colorBlue    = "\u001b[94m"
	colorMagenta = "\u001b[95m"
	colorCyan    = "\u001b[96m"
	colorWhite   = "\u001b[97m"
	colorReset   = "\u001b[0m"
)

type Colors struct {
	base   string
	Output io.Writer
}

func Paint(base ...string) *Colors {
	c := &Colors{
		base:   colorReset,
		Output: colorable.NewColorableStdout(),
	}
	if len(base) > 0 {
		c.base = base[0]
	}
	return c
}

func (c *Colors) Default(color string) {
	c.base = color
}
func (c *Colors) Black(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorBlack, v[0], c.base)
	}
	return colorBlack
}
func (c *Colors) Red(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorRed, v[0], c.base)
	}
	return colorRed
}
func (c *Colors) Green(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorGreen, v[0], c.base)
	}
	return colorGreen
}
func (c *Colors) Yellow(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorYellow, v[0], c.base)
	}
	return colorYellow
}
func (c *Colors) Blue(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorBlue, v[0], c.base)
	}
	return colorBlue
}
func (c *Colors) Magenta(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorMagenta, v[0], c.base)
	}
	return colorMagenta
}
func (c *Colors) Cyan(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorCyan, v[0], c.base)
	}
	return colorCyan
}
func (c *Colors) White(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorWhite, v[0], c.base)
	}
	return colorWhite
}
func (c *Colors) Reset(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorReset, v[0], c.base)
	}
	return colorReset
}
