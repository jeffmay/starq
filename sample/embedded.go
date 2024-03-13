package sample

import (
	_ "embed"
)

//go:embed input/petstore-openapi.json
var PetstoreOpenAPIspecJSON string
