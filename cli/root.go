package cli

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/divolgin/helmsplain/pkg/log"
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
			v := viper.GetViper()

			log.SetDebug(v.GetBool("debug"))

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

	cmd.Flags().Bool("debug", false, "set to true to enable debug output")

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
		return nil
	}

	if isTGZ(arg) {
		err := lookInTGZ(arg)
		if err != nil {
			return errors.Wrapf(err, "look in tgz %s", arg)
		}
		return nil
	}

	printValuesInFile(arg, "")

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

		printValuesInFile(path, root)

		return nil
	})

	return errors.Wrap(err, "walk dir")
}

func printValuesInFile(filename string, filenameMask string) {
	refs, err := values.GetFromFiles(filename)
	if err != nil {
		fmt.Printf("%s\n", filename)
		fmt.Println("    ", err)
		return
	}

	if len(refs) == 0 {
		return
	}

	fmt.Printf("%s\n", strings.TrimPrefix(filename, filenameMask))
	for _, ref := range refs {
		fmt.Println("    ", ref)
	}
}

func isTGZ(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return false
	}
	defer gr.Close()

	return true
}

func lookInTGZ(archiveName string) error {
	tmpDir, err := ioutil.TempDir("", "helmsplain-")
	if err != nil {
		return errors.Wrap(err, "create temp dir")
	}
	os.RemoveAll(tmpDir)

	fileReader, err := os.Open(archiveName)
	if err != nil {
		return errors.Wrap(err, "open archive")
	}
	defer fileReader.Close()

	gzipReader, err := gzip.NewReader(fileReader)
	if err != nil {
		return errors.Wrap(err, "new gzip reader")
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Wrap(err, "read tar header")
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		err = func() error {
			outFilename := filepath.Join(tmpDir, header.Name)
			outPath := filepath.Dir(outFilename)
			err = os.MkdirAll(outPath, 0755)
			if err != nil {
				return errors.Wrap(err, "create file path")
			}

			outFile, err := os.Create(outFilename)
			if err != nil {
				return errors.Wrap(err, "create output file")
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				return errors.Wrapf(err, "copy file %s", header.Name)
			}

			return nil
		}()

		if err != nil {
			return err
		}
	}

	if err := walkDir(tmpDir); err != nil {
		return errors.Wrapf(err, "walk temp dir")
	}

	return nil
}
