package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/divolgin/helmsplain/pkg/values"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "helmsplain [path]",
		Short:        "List all values variables used in templates",
		Long:         `List all values variables used in templates`,
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// v := viper.GetViper()

			if len(args) == 0 {
				// TODO: implement
				panic(errors.New("reading from STDIN is not implemented"))
			}

			for _, arg := range args {
				err := processArg(arg)
				if err != nil {
					fmt.Println(err)
				}
			}

			return nil
		},
	}

	viper.BindPFlags(cmd.Flags())

	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func processArg(arg string) error {
	fileInfo, err := os.Stat(arg)
	if err != nil {
		return errors.Wrap(err, "stat file")
	}

	if fileInfo.IsDir() {
		err := walkDir(arg)
		if err != nil {
			return errors.Wrapf(err, "walk dir %s", arg)
		}
	} else {
		printValuesInFile(arg)
	}

	return nil
}

func walkDir(root string) error {
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "walk path %s", path)
		}

		if info.IsDir() {
			return nil
		}

		printValuesInFile(path)

		return nil
	})

	return errors.Wrap(err, "walk dir")
}

func printValuesInFile(filename string) {
	refs, err := values.GetFromFiles(filename)
	if err != nil {
		fmt.Printf("%s\n", filename)
		fmt.Println("    ", err)
		return
	}

	if len(refs) == 0 {
		return
	}

	fmt.Printf("%s\n", filename)
	for _, ref := range refs {
		fmt.Println("    ", ref)
	}
}
