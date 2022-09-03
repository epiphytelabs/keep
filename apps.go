package keep

import "embed"

//go:embed apps/*.yml
var Apps embed.FS
