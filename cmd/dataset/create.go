package dataset

import (
	"github.com/data-preservation-programs/go-singularity/cmd/cliutil"
	"github.com/data-preservation-programs/go-singularity/database"
	"github.com/data-preservation-programs/go-singularity/handler/dataset"
	"github.com/urfave/cli/v2"
)

var CreateCmd = &cli.Command{
	Name:      "create",
	Usage:     "Create a new dataset",
	ArgsUsage: "DATASET_NAME",
	Description: "DATASET_NAME must be a unique identifier for a dataset\n" +
		"The dataset is a top level object to distinguish different dataset.\n",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "min-size",
			Aliases:  []string{"m"},
			Usage:    "Minimum size of the CAR files to be created",
			Value:    "20GiB",
			Category: "Preparation Parameters",
		},
		&cli.StringFlag{
			Name:     "max-size",
			Aliases:  []string{"M"},
			Usage:    "Maximum size of the CAR files to be created",
			Value:    "30GiB",
			Category: "Preparation Parameters",
		},
		&cli.StringFlag{
			Name:        "piece-size",
			Aliases:     []string{"s"},
			Usage:       "Target piece size of the CAR files used for piece commitment calculation",
			DefaultText: "inferred",
			Category:    "Preparation Parameters",
		},
		&cli.StringSliceFlag{
			Name:        "output-dir",
			Aliases:     []string{"o"},
			Usage:       "Output directory for CAR files",
			DefaultText: "not needed",
			Category:    "Inline Preparation",
		},
		&cli.StringSliceFlag{
			Name:     "encryption-recipient",
			Usage:    "Public key of the encryption recipient",
			Category: "Encryption",
		},
		&cli.StringFlag{
			Name:     "encryption-script",
			Usage:    "Script command to run for custom encryption",
			Category: "Encryption",
		},
	},
	Action: func(c *cli.Context) error {
		db := database.MustOpenFromCLI(c)
		dataset, err := dataset.CreateHandler(
			db,
			dataset.CreateRequest{
				Name:         c.Args().Get(0),
				MinSizeStr:   c.String("min-size"),
				MaxSizeStr:   c.String("max-size"),
				PieceSizeStr: c.String("piece-size"),
				OutputDirs:   c.StringSlice("output-dir"),
				Recipients:   c.StringSlice("encryption-recipients"),
				Script:       c.String("encryption-script")},
		)
		if err != nil {
			return err
		}
		cliutil.PrintToConsole(dataset, c.Bool("json"))
		return nil
	},
}
