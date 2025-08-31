package spec

import (
	"embed"
)

//go:embed openapi.yaml
var File []byte

//go:embed ui/*
var UI embed.FS
