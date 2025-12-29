package generate

import (
	"fmt"

	"github.com/boxboxjason/sonarqube-client-go/pkg/api"
	"github.com/boxboxjason/sonarqube-client-go/pkg/util/strcase"
	jen "github.com/dave/jennifer/jen"
)

// AddServiceTestFile generates unit tests for a service using httptest mock server.
func (gen *Generator) AddServiceTestFile(service *api.WebService) error {
	if service == nil || len(service.Actions) == 0 {
		return nil
	}

	serviceName := service.Path[4:]
	serviceNameCamel := strcase.ToCamel(serviceName)

	file := jen.NewFile(gen.PackageName)
	file.HeaderComment(generatedHeader)
	file.ImportName(pkgNetHTTP, "http")
	file.ImportName("net/http/httptest", "httptest")
	file.ImportName("testing", "testing")
	file.ImportName("encoding/json", "json")

	// Generate test for each action
	for _, action := range service.Actions {
		gen.generateActionTest(file, serviceNameCamel, action)
	}

	err := file.Save(gen.WorkingDir + "/" + GeneratedFilenamePrefix + serviceName + "_service_test.go")
	if err != nil {
		return fmt.Errorf("failed to save service test file: %w", err)
	}

	return nil
}

func (gen *Generator) generateActionTest(file *jen.File, serviceNameCamel string, action api.Action) {
	actionName := strcase.ToCamel(action.Key)
	hasOption := len(action.Params) > 0

	// Determine if the method returns a response value (3 return values) or just (resp, err)
	// This must match the logic in GenerateGoContent
	method := methodGet
	if action.Post {
		method = "POST"
	}

	noResp := false

	switch action.ResponseType {
	case formatJSON, formatTxt, formatLog, formatSvg, formatXML, formatProto:
		noResp = false
	case respNoContent, "":
		noResp = true
	default:
		// jen.For unknown response types, GET returns response, POST doesn't
		if method != methodGet {
			noResp = true
		}
	}

	hasResp := action.HasResponseExample || !action.Post
	isProto := action.ResponseType == formatProto

	testFuncName := "Test" + serviceNameCamel + "_" + actionName

	file.Func().Id(testFuncName).Params(jen.Id("t").Op("*").Qual("testing", "T")).BlockFunc(func(group *jen.Group) {
		// Create mock server
		gen.generateMockServer(group, method, isProto, noResp)
		gen.generateTestExecution(group, serviceNameCamel, actionName, hasOption, hasResp, noResp)
	})
	file.Line()
}

//nolint:nestif
func (gen *Generator) generateTestExecution(group *jen.Group, serviceNameCamel, actionName string, hasOption, hasResp, noResp bool) {
	// Create client pointing to mock server
	group.Comment("Create client pointing to mock server")
	group.List(jen.Id("client"), jen.Id("err")).Op(":=").Id("NewClient").Call(jen.Id("ts").Dot("URL").Op("+").Lit("/api/"), jen.Lit("user"), jen.Lit("pass"))
	group.If(jen.Err().Op("!=").Nil()).Block(
		jen.Id("t").Dot("Fatalf").Call(jen.Lit("failed to create client: %v"), jen.Err()),
	)

	// jen.Call the service method
	group.Comment("jen.Call service method")

	if hasOption {
		group.Id("opt").Op(":=").Op("&").Id(serviceNameCamel + strcase.ToCamel(actionName) + "Option").Values()

		if hasResp {
			group.List(jen.Id("_"), jen.Id("resp"), jen.Id("err")).Op(":=").Id("client").Dot(serviceNameCamel).Dot(actionName).Call(jen.Id("opt"))
		} else {
			group.List(jen.Id("resp"), jen.Id("err")).Op(":=").Id("client").Dot(serviceNameCamel).Dot(actionName).Call(jen.Id("opt"))
		}
	} else {
		if hasResp {
			group.List(jen.Id("_"), jen.Id("resp"), jen.Id("err")).Op(":=").Id("client").Dot(serviceNameCamel).Dot(actionName).Call()
		} else {
			group.List(jen.Id("resp"), jen.Id("err")).Op(":=").Id("client").Dot(serviceNameCamel).Dot(actionName).Call()
		}
	}

	group.If(jen.Err().Op("!=").Nil()).Block(
		jen.Id("t").Dot("Fatalf").Call(jen.Lit(actionName+" failed: %v"), jen.Err()),
	)
	// Check expected status code
	expectedStatus := statusOK
	if noResp {
		expectedStatus = statusNoContent
	}

	group.If(jen.Id("resp").Dot("StatusCode").Op("!=").Lit(expectedStatus)).Block(
		jen.Id("t").Dot("Errorf").Call(jen.Lit(fmt.Sprintf("expected status %d, got %%d", expectedStatus)), jen.Id("resp").Dot("StatusCode")),
	)
}

func (gen *Generator) generateMockServer(group *jen.Group, method string, isProto, noResp bool) {
	group.Comment("Create mock server")
	group.Id("ts").Op(":=").Qual("net/http/httptest", "NewServer").Call(
		jen.Qual(pkgNetHTTP, "HandlerFunc").Call(
			jen.Func().Params(jen.Id("w").Qual(pkgNetHTTP, "ResponseWriter"), jen.Id("r").Op("*").Qual(pkgNetHTTP, "Request")).BlockFunc(func(handlerGroup *jen.Group) {
				handlerGroup.Comment("Verify request method")
				handlerGroup.If(jen.Id("r").Dot("Method").Op("!=").Lit(method)).Block(
					jen.Id("t").Dot("Errorf").Call(jen.Lit("expected method "+method+", got %s"), jen.Id("r").Dot("Method")),
				)
				handlerGroup.Comment("jen.Return mock response")

				switch {
				case isProto:
					// jen.For proto responses, return 200 with valid JSON empty array
					// The Do function uses json.Decode for []byte type
					handlerGroup.Id("w").Dot("Header").Call().Dot("Set").Call(jen.Lit("Content-Type"), jen.Lit("application/json"))
					handlerGroup.Id("w").Dot("WriteHeader").Call(jen.Lit(statusOK))
					handlerGroup.Id("w").Dot("Write").Call(jen.Index().Byte().Call(jen.Lit("[]")))
				case noResp:
					// No response body expected
					handlerGroup.Id("w").Dot("WriteHeader").Call(jen.Lit(statusNoContent))
				default:
					// Use null for JSON which is valid for any pointer/slice type
					handlerGroup.Id("w").Dot("Header").Call().Dot("Set").Call(jen.Lit("Content-Type"), jen.Lit("application/json"))
					handlerGroup.Id("w").Dot("WriteHeader").Call(jen.Lit(statusOK))
					handlerGroup.Id("w").Dot("Write").Call(jen.Index().Byte().Call(jen.Lit("null")))
				}
			}),
		),
	)
	group.Defer().Id("ts").Dot("Close").Call()
}
