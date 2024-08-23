package main

import (
	"aurascan-backend/internal"
	"github.com/urfave/cli/v2"
	"os"
)

// VERSION 版本号，可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "1.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "file coin assistant"
	app.Version = VERSION

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "c",
			Usage:       "-c=配置文件(.toml)",
			Required:    false,
			Value:       "./config.toml",
			DefaultText: "./config.toml",
			Destination: nil,
			HasBeenSet:  false,
		},
	}
	app.Action = func(ctx *cli.Context) error {
		return internal.Run(ctx.Context, ctx.String("c"))
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

}
