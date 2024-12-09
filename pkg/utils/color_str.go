package utils

import (
	"sync"

	"github.com/fatih/color"
)

var (
	colorInsMap = &sync.Map{}
)

func GetColorIns(fgType color.Attribute) *color.Color {
	var colorIns *color.Color
	if v, ok := colorInsMap.Load(fgType); ok {
		colorIns = v.(*color.Color)
	} else {
		colorIns = color.New(fgType)
		colorInsMap.Store(fgType, colorIns)
	}
	return colorIns
}
