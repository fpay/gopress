package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var viewIsPartial bool

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Generate view files (v)",
	Long:  `This command helps generating view template files.`,
	Example: `
  # generete login view file. It is located at ./views directory by default
  gopress make view login

  # use v for short
  gopress make v login

  # generate multiple viewss in a single command
  gopress make view login users/avatar posts/article

  # generate view files in directory other than ./views
  gopress make view login --dir=templates`,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		dir := viper.Sub("view").GetString("dir")
		if moduleRootPath != "" {
			dir = moduleRootPath
		}

		if viewIsPartial {
			dir += "/_partials"
		}

		modules := make([]commonTemplateData, len(args))

		for i, m := range args {
			m = validateModuleFileName(m)
			mt := validateModuleTypeName(m)
			packageName, fileName := parseModulePath(dir, m, fileExtensionHandlebars)
			if err := checkFileName(fileName); err != nil {
				er(err)
			}
			m = validateViewName(dir, m)
			modules[i] = commonTemplateData{packageName, fileName, m, mt}
		}

		for _, data := range modules {
			err := createViewFile(
				data.packageName,
				data.fileName,
				data.moduleName,
				data.moduleTypeName,
			)
			if err != nil {
				er(err)
			}
			fmt.Printf("View file created: %s\n", data.fileName)
		}
	},
}

func init() {
	makeCmd.AddCommand(viewCmd)

	viewCmd.PersistentFlags().BoolVarP(&viewIsPartial, "partial", "p", false, "If the views are partials")

	viper.SetDefault("view.dir", "views")
}

func validateViewName(root, module string) string {
	return strings.TrimPrefix(module, root)
}

func createViewFile(packageName, fileName, moduleName, moduleTypeName string) error {
	var template string

	if viewIsPartial {
		template = `<!-- partial: {{.moduleName}} -->
<div>This is partial {{.moduleName}}</div>`
	} else {
		template = `<!-- template: {{.moduleName}} -->
<div>This is template {{.moduleName}}</div>`
	}

	data := map[string]interface{}{
		"packageName":    packageName,
		"moduleName":     moduleName,
		"moduleTypeName": moduleTypeName,
	}

	viewScript, err := render(template, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(fileName, viewScript)
	if err != nil {
		er(err)
	}

	return nil
}
