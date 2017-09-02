package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// middlewareCmd represents the middleware command
var middlewareCmd = &cobra.Command{
	Use:   "middleware",
	Short: "Generate middleware files (mw)",
	Long:  `This command helps generating middleware functions.`,
	Example: `
  # generete auth middleware file. It is located at ./middlewares directory by default
  gopress make middleware auth

  # use mw for short
  gopress make mw auth

  # generate multiple middlewares in a single command
  gopress make middleware auth throttle

  # generate middleware files in directory other than ./middlewares
  gopress make middleware auth --dir=interceptors`,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"mw"},
	Run: func(cmd *cobra.Command, args []string) {
		dir := viper.Sub("middleware").GetString("dir")
		if moduleRootPath != "" {
			dir = moduleRootPath
		}

		modules := make([]commonTemplateData, len(args))

		for i, m := range args {
			m = validateModuleFileName(m)
			mt := validateModuleTypeName(m)
			packageName, fileName := parseModulePath(dir, m, fileExtensionGolang)
			if err := checkFileName(fileName); err != nil {
				er(err)
			}
			modules[i] = commonTemplateData{packageName, fileName, m, mt}
		}

		for _, data := range modules {
			err := createMiddlewareFile(
				data.packageName,
				data.fileName,
				data.moduleName,
				data.moduleTypeName,
			)
			if err != nil {
				er(err)
			}
			fmt.Printf("Middleware file created: %s\n", data.fileName)
		}
	},
}

func init() {
	makeCmd.AddCommand(middlewareCmd)

	viper.SetDefault("middleware.dir", "middlewares")
}

func createMiddlewareFile(packageName, fileName, moduleName, moduleTypeName string) error {
	template := `package {{.packageName}}

import (
	"github.com/fpay/gopress"
)

// New{{.moduleTypeName}}
func New{{.moduleTypeName}}() gopress.MiddlewareFunc {
	return func(next gopress.HandlerFunc) gopress.HandlerFunc {
		return func(c gopress.Context) error {
			// 如果需要在中间件中使用service，可以通过此方式获取App实例，从App中获取服务容器
			// services := c.(*gopress.AppContext).App().Services()
			return next(c)
		}
	}
}
`

	data := map[string]interface{}{
		"packageName":    packageName,
		"moduleName":     moduleName,
		"moduleTypeName": moduleTypeName + "Middleware",
	}

	middlewareScript, err := render(template, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(fileName, middlewareScript)
	if err != nil {
		er(err)
	}

	return nil
}
