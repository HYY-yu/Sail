package ui

import (
	"embed"
)

//go:embed template
var TemplateFs embed.FS

//go:embed static
var StaticFs embed.FS
