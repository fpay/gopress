package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/viper"
)

type commonTemplateData struct {
	packageName    string
	fileName       string
	moduleName     string
	moduleTypeName string
}

func validateModuleFileName(module string) string {
	var res string = module
	res = strings.Replace(res, "-", "_", -1)
	res = strings.Replace(res, ".", "_", -1)
	return res
}

func validateModuleTypeName(module string) string {
	s := strings.Split(module, "/")
	module = s[len(s)-1]

	i := 0
	l := len(module)
	output := ""

	for i < l {
		if module[i] == '_' {
			i++
			continue
		}

		if i == 0 || module[i-1] == '_' {
			output += string(unicode.ToUpper(rune(module[i])))
			i++
			continue
		}

		output += string(unicode.ToLower(rune(module[i])))
		i++
	}

	return output
}

func parseModulePath(base, module, ext string) (packageName, fileName string) {
	fileName = filepath.Join(viper.GetString("root"), base, module) + ext
	parts := strings.Split(fileName, "/")
	packageName = parts[len(parts)-2]
	return
}

func checkFileName(fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("File %s already exists.", fileName)
}

func er(msg interface{}) {
	fmt.Println("Error: ", msg)
	os.Exit(1)
}

func render(tplStr string, data interface{}) (string, error) {
	tpl, err := template.New("").Parse(tplStr)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	return buf.String(), err
}

func writeStringToFile(path, content string) error {
	base := filepath.Dir(path)
	err := os.MkdirAll(base, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Chmod(0644)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}
