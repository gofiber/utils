package utils

import (
	"fmt"
	"io"

	"github.com/mattn/go-colorable"
)

type Painter struct {
	base   string
	output io.Writer

	Color Color
}

type Color struct {
	Black   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
	White   string
	Reset   string
}

func Colors() *Painter {
	return &Painter{
		output: colorable.NewColorableStdout(),
		Color: Color{
			Black:   "\u001b[90m",
			Red:     "\u001b[91m",
			Green:   "\u001b[92m",
			Yellow:  "\u001b[93m",
			Blue:    "\u001b[94m",
			Magenta: "\u001b[95m",
			Cyan:    "\u001b[96m",
			White:   "\u001b[97m",
			Reset:   "\u001b[0m",
		},
	}
}

func (p *Painter) Default(color ...string) {
	if len(color) > 0 {
		p.base = color[0]
	}
}
func (p *Painter) Black(v interface{}) string {
	return fmt.Sprintf("%s%v%s", p.Color.Black, v, p.base)
}
func (p *Painter) Red(v interface{}) string {
	return fmt.Sprintf("%s%v%s", p.Color.Red, v, p.base)
}
func (p *Painter) Green(v interface{}) string {
	return fmt.Sprintf("%s%v%s", p.Color.Green, v, p.base)
}
func (p *Painter) Yellow(v interface{}) string {
	return fmt.Sprintf("%s%v%s", p.Color.Yellow, v, p.base)
}
func (p *Painter) Blue(v interface{}) string {
	return fmt.Sprintf("%s%v%s", p.Color.Blue, v, p.base)
}
func (p *Painter) Magenta(v interface{}) string {
	return fmt.Sprintf("%s%v%s", p.Color.Magenta, v, p.base)
}
func (p *Painter) Cyan(v interface{}) string {
	return fmt.Sprintf("%s%v%s", p.Color.Cyan, v, p.base)
}
func (p *Painter) White(v interface{}) string {
	return fmt.Sprintf("%s%v%s", p.Color.White, v, p.base)
}
func (p *Painter) Reset(v interface{}) string {
	return p.Color.Reset
}
