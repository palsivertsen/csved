package col

import "github.com/urfave/cli/v2"

func Command() *cli.Command {
	return &cli.Command{
		Name:  "col",
		Usage: "Manipulate columns",
		Subcommands: []*cli.Command{
			{
				Name:   "rm",
				Action: RM,
				Flags: []cli.Flag{
					&cli.IntSliceFlag{
						Name:    "columns",
						Usage:   "Columns to remove",
						Aliases: []string{"c"},
					},
				},
			},
		},
	}
}

type csvReader interface {
	Read() ([]string, error)
}
