package webui

import "embed"

// Files will contain the production web UI after the Vite build output is
// copied into internal/webui/dist during release packaging.
//
//go:embed dist/.gitkeep
var Files embed.FS
