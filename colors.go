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

func Colors(base ...string) *Colors {
	c := &Colors{
		base:   colorReset,
		Output: colorable.NewColorableStdout(),
	}
	if len(base) > 0 {
		c.base = base
	}
	return c
}

func (p *Painter) Default(color string) {
	p.base = color[0]
}
func (p *Painter) Black(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorBlack, v[0], p.base)
	}
	return colorBlack
}
func (p *Painter) Red(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorRed v[0], p.base)
	}
	return colorRed
}
func (p *Painter) Green(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorGreen, v[0], p.base)
	}
	return colorGreen
}
func (p *Painter) Yellow(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorYellow, v[0], p.base)
	}
	return colorYellow
}
func (p *Painter) Blue(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorBlue, v[0], p.base)
	}
	return colorBlue
}
func (p *Painter) Magenta(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorMagenta, v[0], p.base)
	}
	return colorMagenta
}
func (p *Painter) Cyan(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorCyan, v[0], p.base)
	}
	return colorCyan
}
func (p *Painter) White(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorWhite, v[0], p.base)
	}
	return colorWhite
}
func (p *Painter) Reset(v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v%s", colorReset, v[0], p.base)
	}
	return colorReset
}
