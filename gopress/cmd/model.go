package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate model files (m)",
	Long:  `This command helps generating model structs.`,
	Example: `
  # generete user model file. It is located at ./models directory by default
  gopress make model user

  # use m for short
  gopress make m user

  # generate multiple models in a single command
  gopress make model user post

  # generate model files in directory other than ./models
  gopress make model user --dir=entities`,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"m"},
	Run: func(cmd *cobra.Command, args []string) {
		dir := viper.Sub("model").GetString("dir")
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
			err := createModelFile(
				data.packageName,
				data.fileName,
				data.moduleName,
				data.moduleTypeName,
			)
			if err != nil {
				er(err)
			}
			fmt.Printf("model file created: %s", data.fileName)
		}
	},
}

func init() {
	makeCmd.AddCommand(modelCmd)

	viper.BindPFlag("model.dir", makeCmd.PersistentFlags().Lookup("dir"))
	viper.SetDefault("model.dir", "models")
}

func createModelFile(packageName, fileName, moduleName, moduleTypeName string) error {
	template := `package {{.packageName}}

import "time"

type {{.moduleTypeName}} struct {
	ID        uint64 ` + "`" + `gorm:"column:id"` + "`" + `
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *{{.moduleTypeName}}) TableName() string {
	return "{{.moduleName}}s"
}
`

	data := map[string]interface{}{
		"packageName":    packageName,
		"moduleName":     moduleName,
		"moduleTypeName": moduleTypeName,
	}

	modelScript, err := render(template, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(fileName, modelScript)
	if err != nil {
		er(err)
	}

	return nil
}
