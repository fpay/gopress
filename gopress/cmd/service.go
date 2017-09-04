package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Generate service files (svc)",
	Long:  `This command helps generating basic service structs.`,
	Example: `
  # generete cache service file. It is located at ./services directory by default
  gopress make service cache

  # use svc for short
  gopress make svc cache

  # generate multiple services in a single command
  gopress make service cache queue

  # generate service files in directory other than ./services
  gopress make service cache --dir=libs`,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"svc"},
	Run: func(cmd *cobra.Command, args []string) {
		dir := viper.Sub("service").GetString("dir")
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
			err := createServiceFile(
				data.packageName,
				data.fileName,
				data.moduleName,
				data.moduleTypeName,
			)
			if err != nil {
				er(err)
			}
			fmt.Printf("Service file created: %s\n", data.fileName)
		}
	},
}

func init() {
	makeCmd.AddCommand(serviceCmd)

	viper.BindPFlag("service.dir", makeCmd.PersistentFlags().Lookup("dir"))
	viper.SetDefault("service.dir", "services")
}

func createServiceFile(packageName, fileName, moduleName, moduleTypeName string) error {
	template := `package {{.packageName}}

import (
	"github.com/fpay/gopress"
)

const (
	// {{.moduleTypeName}}Name is the identity of {{.moduleName}} service
	{{.moduleTypeName}}Name = "{{.moduleName}}"
)

// {{.moduleTypeName}} type
type {{.moduleTypeName}} struct {
	// Uncomment this line if this service has dependence on other services in the container
	// c *gopress.Container
}

// New{{.moduleTypeName}} returns instance of {{.moduleName}} service
func New{{.moduleTypeName}}() *{{.moduleTypeName}} {
	return new({{.moduleTypeName}})
}

// ServiceName is used to implements gopress.Service
func (s *{{.moduleTypeName}}) ServiceName() string {
	return {{.moduleTypeName}}Name
}

// RegisterContainer is used to implements gopress.Service
func (s *{{.moduleTypeName}}) RegisterContainer(c *gopress.Container) {
	// Uncomment this line if this service has dependence on other services in the container
	// s.c = c
}

func (s *{{.moduleTypeName}}) SampleMethod() {
}
`

	data := map[string]interface{}{
		"packageName":    packageName,
		"moduleName":     moduleName,
		"moduleTypeName": moduleTypeName + "Service",
	}

	serviceScript, err := render(template, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(fileName, serviceScript)
	if err != nil {
		er(err)
	}

	return nil
}
