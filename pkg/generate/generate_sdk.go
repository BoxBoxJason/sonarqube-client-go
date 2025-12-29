package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/boxboxjason/sonarqube-client-go/pkg/api"
	"github.com/boxboxjason/sonarqube-client-go/pkg/response"
	"github.com/boxboxjason/sonarqube-client-go/pkg/util/strcase"
	jen "github.com/dave/jennifer/jen"
	glog "github.com/magicsong/color-glog"
)

// GenerateClient generates the client struct and constructor.
func (gen *Generator) GenerateClient() error {
	gen.client = jen.NewFile(gen.PackageName)
	gen.client.HeaderComment(generatedHeader)
	gen.client.Type().Id("Client").StructFunc(func(group *jen.Group) {
		group.Id("baseURL").Op("*").Qual("net/url", "URL")
		group.List(jen.Id("username"), jen.Id("password"), jen.Id("token")).String()
		group.Id("authType").Id("authType")
		group.Id("httpClient").Op("*").Qual(pkgNetHTTP, "Client")

		for _, service := range gen.services {
			group.Id(strcase.ToCamel(service)).Op("*").Id(strcase.ToCamel(service) + "Service")
		}
	})

	gen.client.Func().Id("NewClient").Params(jen.List(jen.Id("endpoint"), jen.Id("username"), jen.Id("password")).String()).Params(jen.Op("*").Id("Client"), jen.Error()).BlockFunc(func(group *jen.Group) {
		group.Id("c").Op(":=").Op("&Client{username: username, password: password, authType: basicAuth, httpClient: http.DefaultClient}")
		gen.generateClientInitBody(group)
	})

	// Generate NewClientWithToken function
	gen.client.Func().Id("NewClientWithToken").Params(jen.List(jen.Id("endpoint"), jen.Id("token")).String()).Params(jen.Op("*").Id("Client"), jen.Error()).BlockFunc(func(group *jen.Group) {
		group.Id("c").Op(":=").Op("&Client{token: token, authType: privateToken, httpClient: http.DefaultClient}")
		gen.generateClientInitBody(group)
	})

	err := gen.client.Save(gen.WorkingDir + "/" + GeneratedFilenamePrefix + "client.go")
	if err != nil {
		return fmt.Errorf("failed to save client file: %w", err)
	}

	return nil
}

// GenerateGoContent generates the Go code for the service.
func (gen *Generator) GenerateGoContent(packageName string, service *api.WebService) (*jen.File, error) {
	if service == nil {
		return nil, errors.New("service must not be nil")
	}

	if packageName == "" {
		return nil, errors.New("package name is illegal")
	}

	file := jen.NewFile(packageName)

	file.HeaderComment(generatedHeader)
	file.PackageComment(service.Description)
	file.ImportName("github.com/google/glog", "glog")

	name := service.Path[4:]
	// Create Service jen.Struct
	file.Type().Id(strcase.ToCamel(name) + "Service").Struct(
		jen.Id("client").Op("*").Id("Client"),
	).Line()

	// Process examples and generate structs
	gen.processExamples(file, service, name)

	// Create Methods
	for _, item := range service.Actions {
		file.Add(gen.generateServiceActionContent(name, &item))
	}

	return file, nil
}

func (gen *Generator) processExamples(file *jen.File, service *api.WebService, name string) {
	fetcher := response.NewExampleFetcher(gen.Endpoint, gen.Username, gen.Password)

	examples, err := fetcher.GetResponseExample(service)
	if err != nil {
		glog.Warningf("cannot fetch examples of <%s>: %v (continuing without examples)", service.Path, err)

		examples = make([]*response.WebservicesResponseExampleResp, 0)
	}

	// Always set default response types for actions that didn't get one from examples
	gen.setDefaultResponseTypes(file, service, name)

	for _, exam := range examples {
		if exam.Format == formatProto {
			glog.V(1).Infof("The response of action <%s> for api <%s> is proto, using []byte", exam.Name, name)
			gen.setProtoResponseType(service, exam.Name)

			continue
		}

		if exam.Format != formatJSON {
			glog.V(1).Infof("The response of action <%s> for api <%s> is %s, not json", exam.Name, name, exam.Format)

			continue
		}

		if exam.Example != "" {
			gen.generateStructFromExample(file, service, name, exam)
		}
	}
}

func (gen *Generator) setDefaultResponseTypes(file *jen.File, service *api.WebService, name string) {
	for idx := range service.Actions {
		if service.Actions[idx].ResponseType == "" {
			if service.Actions[idx].Post {
				service.Actions[idx].ResponseType = respNoContent
			} else {
				service.Actions[idx].ResponseType = formatJSON
				respName := strcase.ToCamel(name + "_" + service.Actions[idx].Key + "Object")
				file.Commentf("[TODO] cannot fetch response example for <%s>, struct needs to be filled manually", service.Actions[idx].Key)
				file.Type().Id(respName).Struct().Line()
			}
		}
	}
}

func (gen *Generator) setProtoResponseType(service *api.WebService, actionName string) {
	for idx := range service.Actions {
		if service.Actions[idx].Key == actionName {
			service.Actions[idx].ResponseType = formatProto

			break
		}
	}
}

func (gen *Generator) generateStructFromExample(file *jen.File, service *api.WebService, name string, exam *response.WebservicesResponseExampleResp) {
	respName := strcase.ToCamel(name + "_" + exam.Name + "Object")

	stru, err := ConvertStringToStruct(exam.Example, respName)
	if err != nil {
		glog.Warningf("cannot generate resp struct of <%s>,you should manual edit the file %s,esspecial method response", service.Path, gen.WorkingDir+"/"+GeneratedFilenamePrefix+name+"_service.go")
		glog.Errorln(err.Error())
		file.Commentf("[TODO] cannot generate resp struct of <%s>,you should manual edit the file %s,esspecial method response", service.Path, gen.WorkingDir+"/"+GeneratedFilenamePrefix+name+"_service.go")
		file.Type().Id(respName).Struct().Line()
	} else {
		file.Id(stru).Line()
		file.Line()
	}
}

func detectDeprecatedField(field *api.Param) bool {
	return strings.Contains(strings.ToLower(field.Description), "deprecated")
}

// generateServiceActionContent generate code of each service,include api method and related structs.
func (gen *Generator) generateServiceActionContent(serviceName string, action *api.Action) *jen.Statement {
	code := jen.Line()
	hasOption := len(action.Params) > 0
	optionName := strcase.ToCamel(serviceName + "_" + action.Key + "Option")

	if hasOption {
		gen.generateOptionStruct(code, action, optionName, serviceName)
	}

	gen.generateServiceMethod(code, serviceName, action, hasOption, optionName)

	return code
}

func (gen *Generator) generateOptionStruct(code *jen.Statement, action *api.Action, optionName, serviceName string) {
	code.Type().Id(optionName).StructFunc(func(group *jen.Group) {
		for _, field := range action.Params {
			if detectDeprecatedField(&field) {
				glog.V(1).Infof("Detected deprecated field <%s> in <action>:%s,description:%s\n", field.Key, action.Key, field.Description)

				continue
			}

			group.Id(strcase.ToCamel(field.Key)).String().Tag(map[string]string{"url": field.Key + ",omitempty"}).Commentf("Description:\"%s\",ExampleValue:\"%s\"", field.Description, field.ExampleValue)
		}
	}).Line()

	// create valid method
	gen.validation.Func().Params(jen.Id("s").Op("*").Id(strcase.ToCamel(serviceName) + "Service")).Id("Validate" + strcase.ToCamel(action.Key) + "Opt").Params(
		jen.Id("opt").Op("*").Id(optionName)).Params(jen.Error()).Block(
		jen.Return(jen.Nil()),
	)
}

func (gen *Generator) generateServiceMethod(code *jen.Statement, serviceName string, action *api.Action, hasOption bool, optionName string) {
	respName := strcase.ToCamel(serviceName + "_" + action.Key + "Object")
	method := methodGet
	noResp := false

	if action.Post {
		method = "POST"
	}

	code.Commentf("%s %s", strcase.ToCamel(action.Key), action.Description).Line()
	code.Func().Params(jen.Id("s").Op("*").Id(strcase.ToCamel(serviceName) + "Service")).Id(strcase.ToCamel(action.Key)).ParamsFunc(func(group *jen.Group) {
		if hasOption {
			group.Id("opt").Op("*").Id(optionName)
		}
	}).ParamsFunc(func(group *jen.Group) {
		switch action.ResponseType {
		case formatJSON:
			group.Id("v").Op("*").Id(respName)
		case formatTxt, formatLog, formatSvg, formatXML:
			group.Id("v").Op("*").String()

			respName = "string"
		case formatProto:
			// Protocol Buffer binary data
			group.Id("v").Op("[]").Byte()

			respName = "[]byte"
		case respNoContent, "":
			// No response body expected
			noResp = true
		default:
			if method == methodGet {
				group.Id("v").Op("*").Id(respName)
			} else {
				noResp = true
			}
		}

		group.Id("resp").Op("*").Qual(pkgNetHTTP, "Response")
		group.Err().Error()
	}).BlockFunc(func(group *jen.Group) {
		gen.generateServiceMethodBody(group, serviceName, action, hasOption, method, noResp, respName)
	})
}

func (gen *Generator) generateServiceMethodBody(group *jen.Group, serviceName string, action *api.Action, hasOption bool, method string, noResp bool, respName string) {
	if hasOption {
		group.Err().Op("=").Id("s").Dot("Validate" + strcase.ToCamel(action.Key) + "Opt").Call(jen.Id("opt"))
		ErrorHandlerHelper(group)
	}

	group.List(jen.Id("req"), jen.Id("err")).Op(":=").Id("s").Dot("client").Dot("NewRequest").CallFunc(func(reqGroup *jen.Group) {
		reqGroup.Lit(method)
		reqGroup.Lit(serviceName + "/" + action.Key)

		if hasOption {
			reqGroup.Id("opt")
		} else {
			reqGroup.Nil()
		}
	})
	ErrorHandlerHelper(group)

	switch {
	case noResp:
		group.List(jen.Id("resp"), jen.Err()).Op("=").Id("s").Dot("client").Dot("Do").Call(jen.Id("req"), jen.Nil())
		ErrorHandlerHelper(group)
	case action.ResponseType == formatProto:
		// jen.For proto responses, don't use jen.New() since v is []byte not a pointer
		group.List(jen.Id("resp"), jen.Err()).Op("=").Id("s").Dot("client").Dot("Do").Call(jen.Id("req"), jen.Op("&").Id("v"))
		group.If(
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return().List(jen.Nil(), jen.Id("resp"), jen.Err()),
		)
	default:
		group.Id("v").Op("=").New(jen.Id(respName))
		group.List(jen.Id("resp"), jen.Err()).Op("=").Id("s").Dot("client").Dot("Do").Call(jen.Id("req"), jen.Id("v"))
		group.If(
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return().List(jen.Nil(), jen.Id("resp"), jen.Err()),
		)
	}

	group.Return()
}

// ErrorHandlerHelper generates error handling code.
func ErrorHandlerHelper(group *jen.Group) {
	group.If(
		jen.Err().Op("!=").Nil(),
	).Block(
		jen.Return(),
	)
}

func (gen *Generator) generateClientInitBody(group *jen.Group) {
	group.If(
		jen.Id("endpoint").Op("==").Lit(""),
	).Block(
		jen.Id("c").Dot("SetBaseURL").Call(jen.Id("defaultBaseURL")),
	).Else().Block(
		jen.If(jen.Err().Op(" := c.SetBaseURL(endpoint); err != nil").Block(
			jen.Return(jen.Nil(), jen.Err()),
		)),
	)

	for _, service := range gen.services {
		group.Id("c").Dot(strcase.ToCamel(service)).Op("=&").Id(strcase.ToCamel(service) + "Service").Values(jen.Dict{jen.Id("client"): jen.Id("c")})
	}

	group.Return(jen.Id("c"), jen.Nil())
}
