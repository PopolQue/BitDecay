package assets

import "embed"

//go:embed sfx/*.wav fonts/*.otf
var FS embed.FS
