package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/wy-z/tspec/tspec"
)

type cliOpts struct {
	PkgPath       string
	TypeExpr      string
	RefPrefix     string
	IgnoreJSONTag bool
}

//Run runs tspec
func Run(version string) {
	app := cli.NewApp()
	app.Name = "TSpec"
	app.Version = version
	app.Usage = "Parse golang data structure into json schema."

	opts := new(cliOpts)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "package, p",
			Usage:       "package dir `PKG`",
			Value:       ".",
			Destination: &opts.PkgPath,
		},
		cli.StringFlag{
			Name:        "expression, expr",
			Usage:       "type expression `EXPR`",
			Destination: &opts.TypeExpr,
		},
		cli.StringFlag{
			Name:        "ref-prefix, rp",
			Usage:       "the prefix of ref url `PREFIX`",
			Value:       tspec.DefaultRefPrefix,
			Destination: &opts.RefPrefix,
		},
		cli.BoolFlag{
			Name:        "ignore-json-tag",
			Usage:       "ignore json tag",
			Destination: &opts.IgnoreJSONTag,
		},
	}
	app.Action = func(c *cli.Context) (err error) {
		if c.NArg() > 0 {
			opts.TypeExpr = c.Args().Get(0)
		}
		if opts.TypeExpr == "" {
			cli.ShowAppHelp(c)
			return
		}

		parser := tspec.NewParser()
		parserOpts := tspec.DefaultParserOptions
		if opts.RefPrefix != "" {
			parserOpts.RefPrefix = opts.RefPrefix
		}
		parserOpts.IgnoreJSONTag = opts.IgnoreJSONTag
		parser.Options(parserOpts)

		pkg, err := parser.Import(opts.PkgPath)
		if err != nil {
			msg := fmt.Sprintf("failed to import pkg %s: %s", pkg.Name, err)
			err = cli.NewExitError(msg, 1)
			return
		}
		_, err = parser.Parse(pkg, opts.TypeExpr)
		if err != nil {
			msg := fmt.Sprintf("failed to parse type expr %s: %s", opts.TypeExpr, err)
			err = cli.NewExitError(msg, 1)
			return
		}
		defs := parser.Definitions()
		bytes, err := json.MarshalIndent(defs, "", "\t")
		if err != nil {
			msg := fmt.Sprintf("failed to marshal definitions: %s", err)
			err = cli.NewExitError(msg, 1)
			return
		}
		fmt.Println(string(bytes))
		return
	}

	app.Run(os.Args)
}
