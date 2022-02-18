package main

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/stern/stern/stern"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"os"
	"regexp"
	"text/template"
	"time"
)

func runStern() error {
	selector, _ := labels.Parse("carto.run/workload-name=tanzu-java-web-app")
	containerQuery := regexp.MustCompile(".*")

	t := "{{color .ContainerColor .PodName}}{{color .PodColor \"[\"}}{{color .PodColor .ContainerName}}{{color .PodColor \"]\"}} {{.Message}}\n"
	funs := map[string]interface{}{
		"json": func(in interface{}) (string, error) {
			b, err := json.Marshal(in)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
		"color": func(color color.Color, text string) string {
			return color.SprintFunc()(text)
		},
	}
	template, _ := template.New("log").Funcs(funs).Parse(t)

	configStern := stern.Config{
		KubeConfig:     "",
		ContextName:    "",
		Namespaces:     []string{"default"},
		Timestamps:     false,
		Location:       time.Local,
		LabelSelector:  selector,
		ContainerQuery: containerQuery,
		ContainerStates: []stern.ContainerState{
			stern.RUNNING,
			stern.TERMINATED,
		},
		InitContainers: true,
		Since:          3600000000000,

		// PodQuery and FieldSelector are required, but we use LabelSelector instead
		PodQuery:      regexp.MustCompile(""),
		FieldSelector: fields.Everything(),

		Template: template,
		Out:      os.Stdout,
		ErrOut:   os.Stderr,
	}


	//context.Context
	return stern.Run(context.Background(), &configStern)
}

func main() {
	runStern()
}