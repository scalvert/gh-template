package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

const CACHE_DIR string = ".pages"

func main() {
	cmd := &cobra.Command{
		Use:   "template <template name>",
		Short: "Add a file template",
		Args:  cmdutil.ExactArgs(1, "Missing argument for template"),
		RunE:  runTemplate,
	}

	cmd.Flags().StringP("out-dir", "o", ".", "Output directory path")

	rootCmd := &cobra.Command{
		Use:   "pages",
		Short: "GitHub Pages commands",
	}

	rootCmd.AddCommand(cmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getUserHomeDir() (string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

func ensureCacheDir() string {
	homeDir, err := getUserHomeDir()

	if err != nil {
		return ""
	}

	cacheDir := filepath.Join(homeDir, CACHE_DIR)

	_, err = os.Stat(cacheDir)

	if os.IsNotExist(err) {
		os.Mkdir(cacheDir, 0755)
	}

	return cacheDir
}

func runTemplate(cmd *cobra.Command, args []string) error {
	// ctx := context.Background()

	// // Retrieve the template name argument.
	// templateName := args[0]

	// // Retrieve the output directory path option.
	// outDir, _ := cmd.Flags().GetString("out-dir")

	cacheDir := ensureCacheDir()

	fmt.Println(cacheDir)

	// // Decode the template file content.
	// templateContent, err := content.Decode()
	// if err != nil {
	// 	return errors.New("failed to decode template file content: " + err.Error())
	// }

	// // Write the template file to the output directory.
	// err = ioutil.WriteFile(filepath.Join(outDir, templateName), []byte(templateContent), 0644)
	// if err != nil {
	// 	return errors.New("failed to write template file to output directory: " + err.Error())
	// }

	return nil
}
