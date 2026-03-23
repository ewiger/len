package configassets

import "embed"

// DefaultFS contains the bootstrap files used to initialize a LEN config dir.
//go:embed template.toml styles/*.yaml
var DefaultFS embed.FS