package borders

import (
	_ "embed"
)

const (
	BorderSize = 6
)

var (
	//go:embed borders.png
	Borders_png []byte
)
