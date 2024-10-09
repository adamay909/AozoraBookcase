package main

import (
	"embed"

	aozorafs "github.com/adamay909/AozoraBookcase/aozoraFS"
)

//go:embed resources/*
var templateFiles embed.FS

func SetTemplates(lib *aozorafs.Library) {

	lib.ImportTemplates(templateFiles)

}
