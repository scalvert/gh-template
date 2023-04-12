package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/cli/cli/v2/git"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

const CACHE_DIR string = ".pages"

type TemplateOptions struct {
	template string
	outDir   string
}

func main() {
	templateOpts := &TemplateOptions{}

	templateCmd := &cobra.Command{
		Use:   "template <template name>",
		Short: "Add a file template",
		Args:  cmdutil.ExactArgs(1, "Missing argument for template"),
		PreRun: func(cmd *cobra.Command, args []string) {
			cwd, _ := os.Getwd()

			if templateOpts.outDir == "." {
				templateOpts.outDir = cwd
			} else {
				templateOpts.outDir = filepath.Join(cwd, templateOpts.outDir)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			templateOpts.template = args[0]

			return templateRun(templateOpts)
		},
	}

	templateCmd.Flags().StringVarP(&templateOpts.outDir, "out-dir", "o", ".", "Output directory path")

	rootCmd := &cobra.Command{
		Use:   "pages",
		Short: "GitHub Pages commands",
	}

	rootCmd.AddCommand(templateCmd)

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

func getTemplateDirPath(cacheDir string) string {
	return filepath.Join(cacheDir, "README-Template")
}

func cloneRepo(ctx context.Context, repoURL string, clonePath string) {
	client := &git.Client{}
	fmt.Println(clonePath)
	target, err := client.Clone(ctx, repoURL, []string{clonePath})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(target)
}

func ensureTemplateCache(ctx context.Context, cacheDir string, templateName string) {
	templateDir := getTemplateDirPath(cacheDir)

	_, err := os.Stat(templateDir)

	if os.IsNotExist(err) {
		cloneRepo(ctx, "https://linkedin.ghe.com/managed/README-Template", templateDir)
	}
}

func readTemplate(templateName string) (string, error) {
	templatePath := filepath.Join(getTemplateDirPath(ensureCacheDir()), templateName)
	content, err := os.ReadFile(templatePath)

	if err != nil {
		return "", err
	}

	return string(content), nil

}

func writeTemplate(outDir string, templateName string, templateContent string) (string, error) {
	templatePath := filepath.Join(outDir, templateName)

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)

	return templatePath, err
}

func templateRun(templateOpts *TemplateOptions) error {
	ctx := context.Background()

	cacheDir := ensureCacheDir()

	ensureTemplateCache(ctx, cacheDir, templateOpts.template)

	template, err := readTemplate(templateOpts.template)

	if err != nil {
		fmt.Println("No template found with the template name:", templateOpts.template)
		return nil
	}

	templatePath, _ := writeTemplate(templateOpts.outDir, templateOpts.template, template)

	fmt.Println("Template written to:", templatePath)
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
