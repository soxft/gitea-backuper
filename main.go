package main

import (
	"github.com/soxft/gitea-backuper/config"
	"github.com/soxft/gitea-backuper/core"
)

func main() {
	config.Init()
	core.Run()
}
