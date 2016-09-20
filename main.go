package main

import (
	"fmt"
	"github.com/UniversityRadioYork/alias-go/generator"
	"github.com/UniversityRadioYork/alias-go/utils"
	"github.com/urfave/cli"
	"os"
)

func main() {

	var configfilepath string
	var outfile string
	var writeexample string

	app := cli.NewApp()
	app.Name = "alias-go"
	app.HideVersion = true
	app.Usage = "Generates mailing lists"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Chris Taylor",
			Email: "christhebaron@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config-file, config, c",
			Usage:       "Load configuration from `FILE` (required)",
			Destination: &configfilepath,
		},
		cli.StringFlag{
			Name:        "out-filename, out, o",
			Usage:       "Write aliases to `FILE`",
			Value:       "aliases",
			Destination: &outfile,
		},
		cli.StringFlag{
			Name:        "example-config, example, e",
			Usage:       "Write an example config to `FILE`",
			Destination: &writeexample,
		},
	}

	// Before the application runs, let's just do some validation
	app.Before = func(c *cli.Context) error {
		if "" == writeexample {
			if "" == configfilepath {
				return cli.NewExitError("Config file is required", 1)
			}
			if _, err := os.Stat(configfilepath); os.IsNotExist(err) {
				return cli.NewExitError("Invalid config file", 1)
			}
		}
		return nil
	}

	app.After = func(c *cli.Context) error {
		fmt.Fprintln(c.App.Writer, "All done!")
		return nil
	}

	// Now we have passed validation we can get on with it
	app.Action = func(c *cli.Context) error {
		if "" != writeexample {
			err := utils.WriteExampleConfigToFile(writeexample)
			if err != nil {
				cli.NewExitError(err.Error(), 1)
			}
		} else {
			config, err := utils.NewConfigFromFile(configfilepath)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			ury, err := utils.NewURY(config.ApiKey)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			aliases, err := generator.GenerateAliases(ury, config)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			err = utils.WriteAliasesToFile(aliases, outfile)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}
		return nil
	}

	app.Run(os.Args)
}
