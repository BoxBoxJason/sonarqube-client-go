package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"

	"github.com/boxboxjason/sonarqube-client-go/pkg/api"
	"github.com/boxboxjason/sonarqube-client-go/pkg/generate"
	glog "github.com/magicsong/color-glog"
)

var (
	h           bool
	JsonPath    string
	OutputPath  string
	PackageName string
	Endpoint    string
	Username    string
	Password    string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&JsonPath, "f", "", "specify location of api file(only support local file)")
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
		return errors.New("Ilegal package name:" + PackageName)
	}
	if JsonPath == "" {
		return errors.New("Must specify the json location,please add -f [filepath]")
	}
	_, err := os.Stat(JsonPath)
	if err != nil {
		glog.Errorln(err)
		return errors.New("No such api file")
	}
	return nil
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	err := validate()
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	file, err := os.Open(JsonPath)
	if err != nil {
		glog.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	myapi := new(api.API)
	err = decoder.Decode(myapi)
	if err != nil {
		glog.Error(err)
		glog.Errorln("cannot decode api file")
		os.Exit(1)
	}
	err = generate.Build(PackageName, OutputPath, Endpoint, Username, Password, myapi)
	if err != nil {
		glog.Fatal(err)
	} else {
		color.Green("Go files generated successfully")
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, ` generate-go-for-sonarqube version: 0.0.1
Usage: main.go [-h] -f jsonpath  -e endpoint [-n packagename] [-o outputpath]  [-u username] [-p password]

Options:
`)
	flag.PrintDefaults()
}
