package generate

import (
	"fmt"
	"os"
	"strings"

	"github.com/boxboxjason/sonarqube-client-go/pkg/api"
	"github.com/boxboxjason/sonarqube-client-go/pkg/util/strcase"
	jen "github.com/dave/jennifer/jen"
)

// AddIntegrationFile generates integration tests for the service.
func (gen *Generator) AddIntegrationFile(service *api.WebService) error {
	fileName := "integration_testing/" + GeneratedFilenamePrefix + service.Path[4:] + "_service_test.go"
	// Always regenerate integration tests
	// if _, err := os.Stat(fileName); err == nil {
	// 	return nil
	// }

	serviceName := strcase.ToCamel(service.Path[4:])
	// set up

	file := jen.NewFile(gen.PackageName + "_test")
	file.HeaderComment(generatedHeader)
	file.ImportName(gen.CurrentRepo, "")
	// construct

	file.Var().Id("_").Op("=").Qual(pkgGinkgo, "Describe").Call(jen.Lit("SonarCLI integration test"), jen.Func().Call().BlockFunc(func(group *jen.Group) {
		group.Qual(pkgGinkgo, "BeforeEach").Call(jen.Func().Call().Block())
		group.Qual(pkgGinkgo, "JustBeforeEach").Call(jen.Func().Call().Block())

		for _, action := range service.Actions {
			gen.generateIntegrationTest(group, serviceName, service.Path, action)
		}
	}))

	err := file.Save(fileName)
	if err != nil {
		return fmt.Errorf("failed to save integration test file: %w", err)
	}

	return nil
}

func (gen *Generator) generateIntegrationTest(group *jen.Group, serviceName, servicePath string, action api.Action) {
	actionName := strcase.ToCamel(action.Key)
	actionKey := strings.ToLower(serviceName) + "." + strings.ToLower(action.Key)

	if gen.skipActions[actionKey] {
		group.Line().Comment(fmt.Sprintf("// SKIPPED: Test %s - action is dangerous/disruptive", actionName))

		return
	}

	hasResp := gen.hasResponse(action)
	hasOption := len(action.Params) > 0

	// Check if any required parameters have no test value - if so, skip this test
	if gen.shouldSkipTest(serviceName, action) {
		group.Line().Comment(fmt.Sprintf("// SKIPPED: Test %s - required param has no test value", actionName))

		return
	}

	// integration files
	group.Qual(pkgGinkgo, "Describe").Call(jen.Lit("Test "+actionName+" in "+servicePath), jen.Func().Call().BlockFunc(func(g1 *jen.Group) {
		g1.Qual(pkgGinkgo, "It").Call(jen.Lit("Should be ok"), jen.Func().Call().BlockFunc(func(group *jen.Group) {
			gen.generateGinkgoTestBody(group, serviceName, actionName, action, hasOption, hasResp)
		}))
	}))
}

func (gen *Generator) hasResponse(action api.Action) bool {
	switch action.ResponseType {
	case formatJSON, formatTxt, formatLog, formatSvg, formatXML, formatProto:
		return true
	case respNoContent, "":
		return false
	default:
		return !action.Post
	}
}

func (gen *Generator) shouldSkipTest(serviceName string, action api.Action) bool {
	for _, param := range action.Params {
		if param.Required && !detectDeprecatedField(&param) {
			val := gen.GetTestValue(serviceName, action.Key, param.Key)
			if val == "" {
				return true
			}
		}
	}

	return false
}

func (gen *Generator) generateGinkgoTestBody(group *jen.Group, serviceName, actionName string, action api.Action, hasOption, hasResp bool) {
	if hasOption {
		// Use dynamic values for integration tests
		group.Id("opt").Op(":= &").Qual(gen.CurrentRepo, strcase.ToCamel(serviceName+"_"+action.Key+"Option")).Values(jen.DictFunc(gen.generateOptionValues(serviceName, action)))
		group.ListFunc(func(grp *jen.Group) {
			if hasResp {
				grp.Id("v")
			}

			grp.Id("resp")
			grp.Err()
		}).Op(":=").Id("client").Dot(serviceName).Dot(actionName).Call(jen.Id("opt"))
	} else {
		group.ListFunc(func(grp *jen.Group) {
			if hasResp {
				grp.Id("v")
			}

			grp.Id("resp")
			grp.Err()
		}).Op(":=").Id("client").Dot(serviceName).Dot(actionName).Call()
	}

	group.Qual(pkgGomega, "Expect").Call(jen.Err()).Dot("ShouldNot").Call(jen.Qual(pkgGomega, "HaveOccurred").Call())

	if action.ResponseType == respNoContent {
		group.Qual(pkgGomega, "Expect").Call(jen.Id("resp").Dot("StatusCode")).Dot("To").Call(jen.Qual(pkgGomega, "Equal").Call(jen.Lit(statusNoContent)))
	} else {
		group.Qual(pkgGomega, "Expect").Call(jen.Id("resp").Dot("StatusCode")).Dot("To").Call(jen.Qual(pkgGomega, "Equal").Call(jen.Lit(statusOK)))
	}

	if hasResp {
		if action.ResponseType == formatProto {
			group.Qual(pkgGomega, "Expect").Call(jen.Id("v")).Dot("To").Call(jen.Qual(pkgGomega, "Not").Call(jen.Qual(pkgGomega, "BeEmpty").Call()))
		} else {
			group.Qual(pkgGomega, "Expect").Call(jen.Id("v")).Dot("To").Call(jen.Qual(pkgGomega, "Not").Call(jen.Qual(pkgGomega, "BeNil").Call()))
		}
	} else {
		group.Qual(pkgGomega, "Expect").Call(jen.Id("resp").Dot("ContentLength")).Dot("To").Call(jen.Qual(pkgGomega, "Equal").Call(jen.Id("int64").Call(jen.Lit(0))))
	}
}

func (gen *Generator) generateOptionValues(serviceName string, action api.Action) func(jen.Dict) {
	return func(dict jen.Dict) {
		for _, param := range action.Params {
			if detectDeprecatedField(&param) {
				continue
			}

			// Use GetTestValue to get the value
			val := gen.GetTestValue(serviceName, action.Key, param.Key)
			if param.Required {
				if val != "" {
					dict[jen.Id(strcase.ToCamel(param.Key))] = jen.Lit(val)
				} else {
					dict[jen.Id(strcase.ToCamel(param.Key))] = jen.Lit("MUST_EDIT_IT")
				}
			} else {
				if val != "" {
					dict[jen.Id(strcase.ToCamel(param.Key))] = jen.Lit(val)
				}
			}
		}
	}
}

// AddIntegrationSuiteFile generates the suite test file for integration tests.
func (gen *Generator) AddIntegrationSuiteFile() error {
	suiteBody := strings.ReplaceAll(TestSuiteConst, "\r\n", "\n")
	suiteBody = strings.ReplaceAll(suiteBody, "\r", "\n")

	// Replace package name with aliased import
	importStmt := fmt.Sprintf("sonargo \"%s\"", gen.CurrentRepo)
	suiteBody = strings.ReplaceAll(suiteBody, "{REPLACE_PACKAGENAME}", importStmt)

	content := fmt.Sprintf(headerFmt, generatedHeader, "sonargo_test", suiteBody)

	err := os.WriteFile("integration_testing/"+GeneratedFilenamePrefix+"suite_test.go", []byte(content), permFile)
	if err != nil {
		return fmt.Errorf("failed to write suite_test.go: %w", err)
	}

	return nil
}
