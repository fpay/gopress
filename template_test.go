package gopress

import (
	"bytes"
	"os"
	"sync"
	"testing"
)

var (
	testTemplateRoot           = "./tests/templates"
	testTemplateName           = "users/detail"
	testTemplateContent        = "<div>{{> users/avatar }} {{ name }}</div>"
	testTemplatePartialName    = "users/avatar"
	testTemplatePartialContent = "<div>User avatar here</div>"
)

func templateSetup(t *testing.T) {
	t.Helper()

	partialDir := testTemplateRoot + "/_partials/users"
	templateDir := testTemplateRoot + "/users"

	for _, dir := range []string{partialDir, templateDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}

	templateSetupWriteStringToFile(templateDir+"/detail."+handlebarsExtension, testTemplateContent)
	templateSetupWriteStringToFile(partialDir+"/avatar."+handlebarsExtension, testTemplatePartialContent)
	templateSetupWriteStringToFile(partialDir+"/avatar.go", testTemplatePartialContent)
}

func templateTeardown(t *testing.T) {
	t.Helper()
	os.RemoveAll("./tests")
}

func templateSetupWriteStringToFile(fileName, content string) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	f.WriteString(content)
	f.Chmod(0644)
	f.Close()
}

func TestParse(t *testing.T) {
	templateSetup(t)
	defer templateTeardown(t)

	name := "users/detail." + handlebarsExtension
	r := &TemplateRenderer{testTemplateRoot, new(sync.Map)}
	tpl, err := r.Parse(name)
	if err != nil {
		t.Errorf("expect template renderer parse the template: %s", err)
	}

	tpl2, _ := r.Parse(name)
	if tpl2 != tpl {
		t.Errorf("expect template from cache %#v, actual is %#v", tpl, tpl2)
	}

	if _, err := r.Parse("template/not/exists"); err == nil {
		t.Errorf("expect parse file error")
	}
}

func TestRender(t *testing.T) {
	templateSetup(t)
	defer templateTeardown(t)

	buf := new(bytes.Buffer)
	name := "users/detail"
	data := map[string]interface{}{"name": "gopress"}

	r := NewTemplateRenderer(testTemplateRoot)

	if err := r.Render(buf, "template/not/exists", data, nil); err == nil {
		t.Errorf("expect template not rendered")
	}

	if err := r.Render(buf, name, data, nil); err != nil {
		t.Errorf("expect template rendered, actual is %s", err)
	}

	if buf.Len() == 0 {
		t.Errorf("expect rendered template content writen to buffer")
	}
}
