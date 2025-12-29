// Package main is the entry point for the SonarQube client generator.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/boxboxjason/sonarqube-client-go/pkg/api"
	"github.com/boxboxjason/sonarqube-client-go/pkg/generate"
	"github.com/fatih/color"
	glog "github.com/magicsong/color-glog"
)

var (
	//nolint:gochecknoglobals // Flags
	help bool
	//nolint:gochecknoglobals // Flags
	JSONPath string
	//nolint:gochecknoglobals // Flags
	OutputPath string
	//nolint:gochecknoglobals // Flags
	PackageName string
	//nolint:gochecknoglobals // Flags
	Endpoint string
	//nolint:gochecknoglobals // Flags
	Username string
	//nolint:gochecknoglobals // Flags
	Password string
)

//nolint:gochecknoinits // Required for flags
func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&JSONPath, "f", "", "specify location of api file(only support local file)")
	flag.StringVar(&OutputPath, "o", ".", "specify the destination dir, default is current workspace")
	flag.StringVar(&PackageName, "n", "sonarqube", "specify the name of generated package,default is \"sonarqube\"")
	flag.StringVar(&Endpoint, "e", "", "specify the web url of sonarqube")
	flag.StringVar(&Username, "u", "admin", "specify the username,default is \"admin\"")
	flag.StringVar(&Password, "p", "admin", "specify the password ,default is \"admin\"")

	flag.Usage = usage
}

func validate() error {
	reg := regexp.MustCompile(`[a-z]+`)
	if !reg.MatchString(PackageName) {
		return errors.New("illegal package name:" + PackageName)
	}

	if JSONPath == "" {
		return errors.New("must specify the json location,please add -f [filepath]")
	}

	_, err := os.Stat(JSONPath)
	if err != nil {
		glog.Errorln(err)

		return errors.New("no such api file")
	}

	return nil
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()

		return
	}

	err := run()
	if err != nil {
		glog.Fatal(err)
	}

	color.Green("Go files generated successfully")
}

func run() error {
	err := validate()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	//nolint:gosec // G304: User provided path
	file, err := os.Open(JSONPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer func() { _ = file.Close() }()

	decoder := json.NewDecoder(file)
	myapi := new(api.API)

	err = decoder.Decode(myapi)
	if err != nil {
		return fmt.Errorf("cannot decode api file: %w", err)
	}

	err = generate.Build(PackageName, OutputPath, Endpoint, Username, Password, myapi)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, ` generate-go-for-sonarqube version: 0.0.1
Usage: main.go [-h] -f jsonpath  -e endpoint [-n packagename] [-o outputpath]  [-u username] [-p password]
`)
}
