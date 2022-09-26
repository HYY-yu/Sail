package ui

import (
	"embed"
)

//go:embed template
var templateFs embed.FS

//go:embed static
var staticFs embed.FS
