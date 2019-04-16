package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type configmap struct {
	Kind       string            `yaml:"kind"`
	APIVersion string            `yaml:"apiVersion"`
	Metadata   metadata          `yaml:"metadata"`
	BinaryData map[string]string `yaml:"binaryData"`
}

type metadata struct {
	CreationTimestamp time.Time `yaml:"creationTimestamp"`
	Name              string    `yaml:"name"`
	Namespace         string    `yaml:"namespace"`
}

var args struct {
	Name      string
	Namespace string
	Files     []string
}

func parseArgs() {
	for _, arg := range os.Args[1:] {
		switch true {
		case args.Name == "":
			args.Name = arg
		case args.Namespace == "":
			args.Namespace = arg
		default:
			if _, err := os.Lstat(arg); err == nil {
				args.Files = append(args.Files, arg)
			}
		}
	}
}

func main() {
	parseArgs()

	cm := configmap{
		Kind:       "ConfigMap",
		APIVersion: "v1",
		Metadata: metadata{
			CreationTimestamp: time.Now().UTC(),
			Name:              args.Name,
			Namespace:         args.Namespace,
		},
		BinaryData: make(map[string]string),
	}

	b64 := base64.StdEncoding
	for _, f := range args.Files {
		buf, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		_, name := path.Split(f)
		cm.BinaryData[name] = b64.EncodeToString(buf)
	}

	if len(cm.BinaryData) == 0 {
		log.Println("Nothing to save.")
		os.Exit(1)
	}

	buf, err := yaml.Marshal(&cm)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))
}
