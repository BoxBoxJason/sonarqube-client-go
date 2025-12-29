package generate

import "strings"

const (
	defaultTestProject = "test-project"
	testUser           = "test-user"
	testPassword       = "TestPass123!"
	testEmail          = "test@example.com"
	testGroup          = "test-group"
	testProfile        = "test-profile"
	testQualityGate    = "test-quality-gate"
	valTrue            = "true"
	valFalse           = "false"
	valJava            = "java"
	valCoverage        = "coverage"
	valPublic          = "public"
	valOpen            = "OPEN"
	valFixed           = "FIXED"
	valMajor           = "MAJOR"
	valBug             = "BUG"
	valName            = "name"
	dateStart          = "2023-01-01"
	dateEnd            = "2023-12-31"
	testTag            = "test-tag"
	testDescription    = "test-description"
	testComponent      = "test-component"
	valAdmin           = "admin"
	ruleJavaS101       = "java:S101"
	testCategory       = "test-category"
	valGithub          = "github"
	testValue          = "test-value"
	valTRK             = "TRK"
)

// GetTestValue returns a valid test value for the given service, action, and parameter.
// Returns empty string for optional parameters we don't want to set.
func (gen *Generator) GetTestValue(service, action, param string) string {
	paramLower := strings.ToLower(param)
	serviceLower := strings.ToLower(service)
	actionLower := strings.ToLower(action)

	// Service+Action+Param specific overrides first
	key := serviceLower + "." + actionLower + "." + paramLower
	if val, ok := gen.getSpecificValue(key); ok {
		return val
	}

	// Service+Param specific overrides
	key = serviceLower + "." + paramLower
	if val, ok := gen.getSpecificValue(key); ok {
		return val
	}

	// Action+Param specific overrides
	key = actionLower + "." + paramLower
	if val, ok := gen.getSpecificValue(key); ok {
		return val
	}

	// Check if this param should be skipped
	if gen.shouldSkipParam(paramLower) {
		return ""
	}

	// Check explicit mappings
	if val, ok := gen.getExplicitMapping(paramLower); ok {
		return val
	}

	// General rules based on parameter name patterns
	return getPatternBasedValue(paramLower)
}

func (gen *Generator) getSpecificValue(key string) (string, bool) {
	if val, exists := gen.specificValues[key]; exists {
		return val, true
	}

	return "", false
}

func (gen *Generator) shouldSkipParam(param string) bool {
	return gen.skipParams[param]
}

func (gen *Generator) getExplicitMapping(param string) (string, bool) {
	if val, ok := gen.explicitMappings[param]; ok {
		return val, true
	}

	return "", false
}

func getPatternBasedValue(param string) string {
	if val := getProjectRelatedValue(param); val != "" {
		return val
	}

	if val := getUserRelatedValue(param); val != "" {
		return val
	}

	if val := getCommonValue(param); val != "" {
		return val
	}

	return ""
}

func getProjectRelatedValue(param string) string {
	if strings.Contains(param, "projectkey") || (strings.Contains(param, "project") && strings.Contains(param, "key")) {
		return defaultTestProject
	}

	if strings.Contains(param, "project") {
		return defaultTestProject
	}

	if strings.Contains(param, "key") {
		return defaultTestProject
	}

	if strings.Contains(param, "name") {
		return "Test Project"
	}

	if strings.Contains(param, "component") {
		return defaultTestProject
	}

	return ""
}

func getUserRelatedValue(param string) string {
	if strings.Contains(param, "login") {
		return testUser
	}

	if strings.Contains(param, "password") {
		return testPassword
	}

	if strings.Contains(param, "email") {
		return testEmail
	}

	if strings.Contains(param, "group") {
		return testGroup
	}

	if strings.Contains(param, "permission") {
		return "scan"
	}

	return ""
}

func getCommonValue(param string) string {
	if val := getCommonValuePart1(param); val != "" {
		return val
	}

	return getCommonValuePart2(param)
}

func getCommonValuePart1(param string) string {
	if strings.Contains(param, "metric") {
		return valCoverage
	}

	if strings.Contains(param, "language") {
		return valJava
	}

	if strings.Contains(param, "rule") {
		return "java:S106"
	}

	if strings.Contains(param, "profile") {
		return testProfile
	}

	if strings.Contains(param, "organization") {
		return "default-organization"
	}

	return ""
}

func getCommonValuePart2(param string) string {
	if strings.Contains(param, "template") {
		return "default_template"
	}

	if strings.Contains(param, "visibility") {
		return valPublic
	}

	if strings.Contains(param, "gate") {
		return testQualityGate
	}

	if strings.Contains(param, "local") {
		return valTrue
	}

	if strings.Contains(param, "url") {
		return "https://example.com"
	}

	return ""
}

func loadSpecificValues() map[string]string {
	values := make(map[string]string)
	loadUsersValues(values)
	loadProjectsValues(values)
	loadQualityGatesValues(values)
	loadQualityProfilesValues(values)
	loadMiscValues(values)

	return values
}

func loadUsersValues(values map[string]string) {
	values["users.create.login"] = "test-user-new"
	values["users.create.name"] = "Test User New"
	values["users.create.password"] = testPassword
	values["users.changepassword.login"] = testUser
	values["users.changepassword.password"] = "TestNewPass123!"
	values["users.sethomepage.type"] = "PROJECT"
	values["users.dismissnotice.notice"] = "educationPrinciples"
	values["users.updateidentityprovider.newexternalprovider"] = ""
	values["users.updatelogin.newlogin"] = ""
}

func loadProjectsValues(values map[string]string) {
	values["projects.create.project"] = "test-project-new"
	values["projects.create.name"] = "Test Project New"
	values["projects.updatedefaultvisibility.projectvisibility"] = valPublic
	values["projects.updatevisibility.visibility"] = valPublic
	values["projects.updatekey.from"] = defaultTestProject
	values["projects.updatekey.to"] = "test-project-renamed"
	values["projectlinks.create.projectkey"] = defaultTestProject
	values["projectlinks.search.projectkey"] = defaultTestProject
	values["projectlinks.create.url"] = "https://example.com/link"
	values["projecttags.set.project"] = defaultTestProject
	values["projecttags.set.tags"] = testTag
}

func loadQualityGatesValues(values map[string]string) {
	values["qualitygates.create.name"] = "test-quality-gate-new"
	values["qualitygates.copy.name"] = "test-quality-gate-copy"
	values["qualitygates.rename.name"] = "test-quality-gate-renamed"
	values["qualitygates.select.gatename"] = testQualityGate
	values["qualitygates.deselect.gatename"] = testQualityGate
	values["qualitygates.destroy.name"] = testQualityGate
	values["qualitygates.setasdefault.name"] = testQualityGate
	values["qualitygates.show.name"] = testQualityGate
	values["qualitygates.getbyproject.project"] = defaultTestProject
	values["qualitygates.createcondition.gatename"] = testQualityGate
	values["qualitygates.createcondition.metric"] = valCoverage
	values["qualitygates.createcondition.op"] = "LT"
	values["qualitygates.createcondition.error"] = "80"
	values["qualitygates.updatecondition.id"] = ""
	values["qualitygates.deletecondition.id"] = ""
}

func loadQualityProfilesValues(values map[string]string) {
	values["qualityprofiles.create.language"] = valJava
	values["qualityprofiles.create.name"] = "test-profile-new"
	values["qualityprofiles.copy.fromkey"] = testProfile
	values["qualityprofiles.copy.toname"] = "test-profile-copy"
	values["qualityprofiles.changeparent.language"] = valJava
	values["qualityprofiles.changeparent.qualityprofile"] = testProfile
	values["qualityprofiles.compare.leftkey"] = ""
	values["qualityprofiles.compare.rightkey"] = ""
	values["qualityprofiles.restore.backup"] = ""
}

func loadMiscValues(values map[string]string) {
	loadMiscValuesPart1(values)
	loadMiscValuesPart2(values)
	loadMiscValuesPart3(values)
}

func loadMiscValuesPart1(values map[string]string) {
	values["notifications.add.type"] = "NewIssues"
	values["notifications.remove.type"] = "NewIssues"
	values["notifications.add.project"] = defaultTestProject
	values["notifications.remove.project"] = defaultTestProject
	values["newcodeperiods.set.type"] = "NUMBER_OF_DAYS"
	values["newcodeperiods.set.value"] = "30"
	values["webhooks.create.name"] = "test-webhook-new"
	values["webhooks.create.url"] = "https://example.com/webhook"
	values["webhooks.create.project"] = defaultTestProject
	values["settings.set.key"] = "sonar.core.serverbaseurl"
	values["settings.set.value"] = "http://127.0.0.1:9000"
	values["settings.reset.keys"] = "sonar.core.serverbaseurl"
	values["system.changeloglevel.level"] = "INFO"
	values["dismissmessage.dismiss.messagetype"] = "GLOBAL_NCD_PAGE_90"
	values["dismissmessage.check.messagetype"] = "GLOBAL_NCD_PAGE_90"
}

func loadMiscValuesPart2(values map[string]string) {
	// Empty values for skipped params
	values["almsettings.createazure.key"] = ""
	values["almsettings.createbitbucket.key"] = ""
	values["almsettings.createbitbucketcloud.key"] = ""
	values["almsettings.creategithub.key"] = ""
	values["almsettings.creategitlab.key"] = ""
	values["almsettings.updateazure.key"] = ""
	values["almsettings.updatebitbucket.key"] = ""
	values["almsettings.updatebitbucketcloud.key"] = ""
	values["almsettings.updategithub.key"] = ""
	values["almsettings.updategitlab.key"] = ""
	values["almsettings.validate.key"] = ""
	values["almsettings.list.project"] = ""
	values["almsettings.getbinding.project"] = ""
	values["sources.index.resource"] = ""
	values["sources.lines.key"] = ""
	values["sources.raw.key"] = ""
	values["sources.scm.key"] = ""
	values["sources.show.key"] = ""
	values["sources.issuesnippets.issuekey"] = ""
	values["ce.task.id"] = ""
	values["ce.cancel.id"] = ""
	values["ce.submit.report"] = ""
	values["ce.dismissmessage.id"] = ""
	values["webhooks.update.webhook"] = ""
	values["webhooks.delete.webhook"] = ""
	values["webhooks.delivery.deliveryid"] = ""
	values["webhooks.deliveries.webhook"] = ""
}

func loadMiscValuesPart3(values map[string]string) {
	values["rules.create.customkey"] = ""
	values["rules.create.templatekey"] = ""
	values["rules.create.markdowndescription"] = ""
	values["settings.encrypt.value"] = ""
	values["emails.send.message"] = ""
	values["emails.send.to"] = ""
	values["emails.send.subject"] = ""
	values["plugins.download.plugin"] = ""
	values["plugins.install.key"] = ""
	values["plugins.uninstall.key"] = ""
	values["plugins.update.key"] = ""
	values["permissions.deletetemplate.templateid"] = ""
	values["permissions.deletetemplate.templatename"] = ""
	values["permissions.updatetemplate.id"] = ""
	values["measures.search.metrickeys"] = valCoverage
	values["project_branches.set_automatic_deletion_protection.value"] = valTrue
	values["projectbranches.set_automatic_deletion_protection.value"] = valTrue
	values["components.search_projects.facets"] = valCoverage
	values["components.search_projects.f"] = "analysisDate"
	values["components.suggestions.more"] = valTRK
	values["projectanalyses.createevent.analysis"] = ""
	values["projectanalyses.delete.analysis"] = ""
	values["projectanalyses.deleteevent.event"] = ""
	values["projectanalyses.setevent.event"] = ""
	values["projectbranches.delete.branch"] = ""
	values["projectbranches.rename.name"] = ""
	values["projectbranches.setnewcodeperiod.branch"] = ""
	values["push.sonarlint.languages"] = ""
	values["push.sonarlint.projectkeys"] = ""
}

func loadSkipParams() map[string]bool {
	return map[string]bool{
		"backup":                true,
		"report":                true,
		"profile":               true,
		"project_analyses":      true,
		"events":                true,
		"issues":                true,
		"components":            true,
		"metrics":               true,
		"rules":                 true,
		"users":                 true,
		"groups":                true,
		"permissions":           true,
		"quality_gates":         true,
		"quality_profiles":      true,
		"settings":              true,
		"webhooks":              true,
		"notifications":         true,
		"new_code_periods":      true,
		"alm_settings":          true,
		"projects":              true,
		"project_links":         true,
		"project_tags":          true,
		"project_branches":      true,
		"project_pull_requests": true,
		"hotspots":              true,
		"favorites":             true,
		"duplications":          true,
		"developers":            true,
		"ce":                    true,
		"authentication":        true,
		"analysis_cache":        true,
		"analysis_reports":      true,
		"editions":              true,
		"health":                true,
		"l10n":                  true,
		"languages":             true,
		"measures":              true,
		"monitoring":            true,
		"navigation":            true,
		"plugins":               true,
		"server":                true,
		"sources":               true,
		"system":                true,
		"update_center":         true,
		"user_groups":           true,
		"user_tokens":           true,
		"webservices":           true,
	}
}

func loadExplicitMappings() map[string]string {
	values := make(map[string]string)
	loadCommonMappings(values)
	loadSearchMappings(values)
	loadDateMappings(values)
	loadFilterMappings(values)
	loadMiscMappings(values)

	return values
}

func loadCommonMappings(values map[string]string) {
	values["login"] = testUser
	values["project"] = defaultTestProject
	values["projectKey"] = defaultTestProject
	values["key"] = defaultTestProject
	values["name"] = "test-name"
	values["description"] = testDescription
	values["password"] = "test-password"
	values["email"] = testEmail
	values["scmAccount"] = "test-scm-account"
	values["organization"] = "test-org"
	values["component"] = testComponent
	values["componentKey"] = testComponent
	values["branch"] = "main"
	values["pullRequest"] = "1"
	values["analysisId"] = "test-analysis-id"
	values["taskId"] = "test-task-id"
	values["webhookId"] = "test-webhook-id"
	values["groupId"] = "test-group-id"
	values["permission"] = valAdmin
	values["templateId"] = "test-template-id"
	values["profileKey"] = "test-profile-key"
	values["gateId"] = "test-gate-id"
	values["metric"] = valCoverage
	values["language"] = valJava
	values["rule"] = ruleJavaS101
	values["severity"] = valMajor
	values["type"] = valBug
	values["status"] = valOpen
	values["resolution"] = valFixed
	values["format"] = "json"
	values["visibility"] = valPublic
	values["qualifier"] = valTRK
	values["strategy"] = "static"
}

func loadSearchMappings(values map[string]string) {
	values["workerCount"] = "1"
	values["timeout"] = "30"
	values["page"] = "1"
	values["pageSize"] = "100"
	values["p"] = "1"
	values["ps"] = "100"
	values["asc"] = valTrue
	values["s"] = valName
	values["f"] = valName
	values["q"] = "test-query"
	values["text"] = "test-text"
}

func loadDateMappings(values map[string]string) {
	values["from"] = dateStart
	values["to"] = dateEnd
	values["date"] = dateStart
	values["createdAfter"] = dateStart
	values["createdBefore"] = dateEnd
	values["createdAt"] = dateStart
	values["startedAfter"] = dateStart
	values["startedBefore"] = dateEnd
	values["executedAfter"] = dateStart
	values["executedBefore"] = dateEnd
	values["analyzedBefore"] = dateEnd
	values["analyzedAfter"] = dateStart
}

func loadFilterMappings(values map[string]string) {
	values["onProvisionedOnly"] = valFalse
	values["selected"] = "all"
	values["more"] = valFalse
	values["onlyMine"] = valFalse
	values["resolved"] = valFalse
	values["assigned"] = valFalse
	values["planned"] = valFalse
	values["tags"] = testTag
	values["types"] = valBug
	values["severities"] = valMajor
	values["statuses"] = valOpen
	values["resolutions"] = valFixed
	values["facets"] = "types"
	values["fields"] = "key,name"
	values["additionalFields"] = "metrics"
	values["metricKeys"] = valCoverage
	values["ruleKeys"] = ruleJavaS101
	values["componentKeys"] = testComponent
	values["projectKeys"] = defaultTestProject
	values["gateNames"] = "test-gate"
	values["profileNames"] = testProfile
	values["groupNames"] = testGroup
	values["logins"] = testUser
	values["emails"] = testEmail
	values["languages"] = valJava
	values["repositories"] = "test-repo"
	values["versions"] = "1.0"
	values["categories"] = testCategory
}

func loadMiscMappings(values map[string]string) {
	loadMiscMappingsPart1(values)
	loadMiscMappingsPart2(values)
	loadMiscMappingsPart3(values)
	loadMiscMappingsPart4(values)
}

func loadMiscMappingsPart1(values map[string]string) {
	values["cwe"] = "123"
	values["owaspTop10"] = "a1"
	values["sansTop25"] = "insecure-interaction"
	values["sonarsourceSecurity"] = "sql-injection"
	values["cleanCodeAttributeCategories"] = "adaptability"
	values["impactSoftwareQualities"] = "security"
	values["impactSeverities"] = "high"
	values["issueStatuses"] = valOpen
	values["prioritizedRule"] = valTrue
	values["isNewCode"] = valTrue
	values["inNewCodePeriod"] = valTrue
	values["hotspots"] = "test-hotspot"
	values["files"] = "test-file"
	values["authors"] = "test-author"
	values["assignees"] = "test-assignee"
	values["admins"] = "test-admin"
	values["users"] = testUser
	values["groups"] = testGroup
	values["permissions"] = valAdmin
	values["tokens"] = "test-token"
	values["notifications"] = "test-notification"
	values["webhooks"] = "test-webhook"
	values["applications"] = "test-app"
	values["portfolios"] = "test-portfolio"
	values["subportfolios"] = "test-subportfolio"
	values["views"] = "test-view"
	values["subviews"] = "test-subview"
	values["devopsPlatforms"] = valGithub
	values["alm"] = valGithub
}

func loadMiscMappingsPart2(values map[string]string) {
	values["configuration"] = "test-config"
	values["definition"] = "test-def"
	values["value"] = testValue
	values["values"] = testValue
	values["defaultValue"] = "test-default"
	values["description"] = testDescription
	values["category"] = testCategory
	values["subCategory"] = "test-subcategory"
	values["options"] = "test-option"
	values["multiValues"] = valTrue
	values["fieldValues"] = testValue
	values["property"] = "test-prop"
	values["setting"] = "test-setting"
	values["componentId"] = "test-component-id"
	values["projectId"] = "test-project-id"
	values["organizationId"] = "test-org-id"
	values["developerId"] = "test-dev-id"
	values["userId"] = "test-user-id"
	values["userLogin"] = testUser
	values["userEmail"] = testEmail
	values["userName"] = "Test User"
	values["groupName"] = "Test Group"
	values["permissionName"] = valAdmin
	values["templateName"] = "Test Template"
	values["profileName"] = "Test Profile"
	values["gateName"] = "Test Gate"
	values["metricName"] = valCoverage
	values["languageName"] = "Java"
	values["ruleName"] = "Test Rule"
	values["ruleKey"] = ruleJavaS101
	values["ruleRepo"] = valJava
}

func loadMiscMappingsPart3(values map[string]string) {
	values["severityName"] = valMajor
	values["typeName"] = valBug
	values["statusName"] = valOpen
	values["resolutionName"] = valFixed
	values["formatName"] = "json"
	values["visibilityName"] = valPublic
	values["qualifierName"] = "TRK"
	values["strategyName"] = "static"
	values["workerCountName"] = "1"
	values["timeoutName"] = "30"
	values["pageName"] = "1"
	values["pageSizeName"] = "100"
	values["pName"] = "1"
	values["psName"] = "100"
	values["ascName"] = valTrue
	values["sName"] = valName
	values["fName"] = valName
	values["qName"] = "test-query"
	values["textName"] = "test-text"
	values["fromName"] = dateStart
	values["toName"] = dateEnd
	values["dateName"] = dateStart
	values["createdAfterName"] = dateStart
	values["createdBeforeName"] = dateEnd
	values["createdAtName"] = dateStart
	values["startedAfterName"] = dateStart
	values["startedBeforeName"] = dateEnd
	values["executedAfterName"] = dateStart
	values["executedBeforeName"] = dateEnd
	values["analyzedBeforeName"] = dateEnd
	values["analyzedAfterName"] = dateStart
}

func loadMiscMappingsPart4(values map[string]string) {
	values["onProvisionedOnlyName"] = valFalse
	values["selectedName"] = "all"
	values["moreName"] = valFalse
	values["onlyMineName"] = valFalse
	values["resolvedName"] = valFalse
	values["assignedName"] = valFalse
	values["plannedName"] = valFalse
	values["tagsName"] = testTag
	values["typesName"] = valBug
	values["severitiesName"] = valMajor
	values["statusesName"] = valOpen
	values["resolutionsName"] = valFixed
	values["facetsName"] = "types"
	values["fieldsName"] = "key,name"
	values["additionalFieldsName"] = "metrics"
	values["metricKeysName"] = valCoverage
	values["ruleKeysName"] = ruleJavaS101
	values["componentKeysName"] = testComponent
	values["projectKeysName"] = defaultTestProject
	values["gateNamesName"] = "test-gate"
	values["profileNamesName"] = testProfile
	values["groupNamesName"] = testGroup
	values["loginsName"] = testUser
	values["emailsName"] = testEmail
	values["languagesName"] = valJava
	values["repositoriesName"] = "test-repo"
	values["versionsName"] = "1.0"
	values["categoriesName"] = testCategory
	values["cweName"] = "123"
	values["owaspTop10Name"] = "a1"
	values["sansTop25Name"] = "insecure-interaction"
	values["sonarsourceSecurityName"] = "sql-injection"
	values["cleanCodeAttributeCategoriesName"] = "adaptability"
	values["impactSoftwareQualitiesName"] = "security"
	values["impactSeveritiesName"] = "high"
	loadMiscMappingsPart5(values)
}

func loadMiscMappingsPart5(values map[string]string) {
	values["issueStatusesName"] = valOpen
	values["prioritizedRuleName"] = valTrue
	values["isNewCodeName"] = valTrue
	values["inNewCodePeriodName"] = valTrue
	values["hotspotsName"] = "test-hotspot"
	values["filesName"] = "test-file"
	values["authorsName"] = "test-author"
	values["assigneesName"] = "test-assignee"
	values["adminsName"] = "test-admin"
	values["usersName"] = testUser
	values["groupsName"] = testGroup
	values["permissionsName"] = valAdmin
	values["tokensName"] = "test-token"
	values["notificationsName"] = "test-notification"
	values["webhooksName"] = "test-webhook"
	values["applicationsName"] = "test-app"
	values["portfoliosName"] = "test-portfolio"
	values["subportfoliosName"] = "test-subportfolio"
	values["viewsName"] = "test-view"
	values["subviewsName"] = "test-subview"
	values["devopsPlatformsName"] = valGithub
	values["almName"] = valGithub
	values["configurationName"] = "test-config"
	values["definitionName"] = "test-def"
	values["valueName"] = testValue
	values["valuesName"] = testValue
	values["defaultValueName"] = "test-default"
	values["descriptionName"] = testDescription
	values["categoryName"] = testCategory
	values["subCategoryName"] = "test-subcategory"
	values["optionsName"] = "test-option"
	values["multiValuesName"] = valTrue
	values["fieldValuesName"] = testValue
	values["propertyName"] = "test-prop"
	values["settingName"] = "test-setting"
}

func loadSkipActions() map[string]bool {
	actions := make(map[string]bool)
	loadSkipActionsPart1(actions)
	loadSkipActionsPart2(actions)
	loadSkipActionsPart3(actions)
	loadSkipActionsPart4(actions)

	return actions
}

func loadSkipActionsPart1(actions map[string]bool) {
	actions["system.restart"] = true
	actions["system.change_log_level"] = true
	actions["system.liveness"] = true
	actions["system.logs"] = true
	actions["plugins.install"] = true
	actions["plugins.uninstall"] = true
	actions["plugins.update"] = true
	actions["plugins.cancel_all"] = true
	actions["users.deactivate"] = true
	actions["users.anonymize"] = true
	actions["projects.bulk_delete"] = true
	actions["projects.delete"] = true
	actions["qualitygates.destroy"] = true
	actions["qualityprofiles.delete"] = true
	actions["usergroups.delete"] = true
	actions["permissions.delete_template"] = true
}

func loadSkipActionsPart2(actions map[string]bool) {
	actions["sources.index"] = true
	actions["sources.lines"] = true
	actions["sources.raw"] = true
	actions["sources.scm"] = true
	actions["sources.show"] = true
	actions["sources.issue_snippets"] = true
	actions["analysiscache.get"] = true
	actions["analysiscache.clear"] = true
	actions["analysisreports.is_queue_empty"] = true
	actions["h.file"] = true
	actions["h.index"] = true
	actions["h.project"] = true
	actions["push.sonarlint"] = true
	actions["push.sonarlint_events"] = true
	actions["rules.list"] = true
	actions["almsettings.create_azure"] = true
	actions["almsettings.create_bitbucket"] = true
	actions["almsettings.create_bitbucketcloud"] = true
	actions["almsettings.create_github"] = true
	actions["almsettings.create_gitlab"] = true
	actions["almsettings.update_azure"] = true
	actions["almsettings.update_bitbucket"] = true
	actions["almsettings.update_bitbucketcloud"] = true
	actions["almsettings.update_github"] = true
	actions["almsettings.update_gitlab"] = true
	actions["almsettings.validate"] = true
	actions["almsettings.delete"] = true
	actions["almsettings.list"] = true
	actions["almsettings.list_definitions"] = true
	actions["almsettings.get_binding"] = true
	actions["almsettings.count_binding"] = true
}

func loadSkipActionsPart3(actions map[string]bool) {
	actions["almintegrations.check_pat"] = true
	actions["almintegrations.get_github_client_id"] = true
	actions["almintegrations.import_azure_project"] = true
	actions["almintegrations.import_bitbucketcloud_repo"] = true
	actions["almintegrations.import_bitbucketserver_project"] = true
	actions["almintegrations.import_github_project"] = true
	actions["almintegrations.import_gitlab_project"] = true
	actions["almintegrations.list_azure_projects"] = true
	actions["almintegrations.list_bitbucketserver_projects"] = true
	actions["almintegrations.list_github_organizations"] = true
	actions["almintegrations.list_github_repositories"] = true
	actions["almintegrations.search_azure_repos"] = true
	actions["almintegrations.search_bitbucketcloud_repos"] = true
	actions["almintegrations.search_bitbucketserver_repos"] = true
	actions["almintegrations.search_gitlab_repos"] = true
	actions["almintegrations.set_pat"] = true
	actions["githubprovisioning.check"] = true
	actions["githubprovisioning.create_role_mapping"] = true
	actions["githubprovisioning.delete_role_mapping"] = true
	actions["githubprovisioning.get_status"] = true
	actions["githubprovisioning.role_mappings"] = true
	actions["githubprovisioning.update_role_mapping"] = true
	actions["emails.send"] = true
	actions["settings.encrypt"] = true
	actions["monitoring.metrics"] = true
	actions["projectlinks.create"] = true
	actions["projectlinks.search"] = true
	actions["projectlinks.delete"] = true
	actions["projectanalyses.create_event"] = true
	actions["projectanalyses.delete"] = true
	actions["projectanalyses.delete_event"] = true
	actions["projectanalyses.set_event"] = true
	actions["projectanalyses.search"] = true
	actions["projectbranches.delete"] = true
	actions["projectbranches.list"] = true
	actions["projectbranches.rename"] = true
	actions["projectbranches.set_new_code_period"] = true
}

func loadSkipActionsPart4(actions map[string]bool) {
	actions["hotspots.add_comment"] = true
	actions["hotspots.assign"] = true
	actions["hotspots.change_status"] = true
	actions["hotspots.delete_comment"] = true
	actions["hotspots.edit_comment"] = true
	actions["hotspots.search"] = true
	actions["hotspots.show"] = true
	actions["hotspots.list"] = true
	actions["hotspots.pull"] = true
	actions["issues.add_comment"] = true
	actions["issues.anticipated_transitions"] = true
	actions["issues.assign"] = true
	actions["issues.bulk_change"] = true
	actions["issues.changelog"] = true
	actions["issues.delete_comment"] = true
	actions["issues.do_transition"] = true
	actions["issues.edit_comment"] = true
	actions["issues.list"] = true
	actions["issues.pull"] = true
	actions["issues.pull_taint"] = true
	actions["issues.reindex"] = true
	actions["issues.search"] = true
	actions["issues.set_severity"] = true
	actions["issues.set_tags"] = true
	actions["issues.set_type"] = true
	actions["issues.tags"] = true
	actions["issues.authors"] = true
	actions["issues.component_tags"] = true
	loadSkipActionsPart5(actions)
}

func loadSkipActionsPart5(actions map[string]bool) {
	actions["measures.component"] = true
	actions["measures.component_tree"] = true
	actions["measures.search_history"] = true
	actions["ce.activity"] = true
	actions["ce.activity_status"] = true
	actions["ce.analysis_status"] = true
	actions["ce.cancel"] = true
	actions["ce.cancel_all"] = true
	actions["ce.component"] = true
	actions["ce.dismiss_analysis_warning"] = true
	actions["ce.indexation_status"] = true
	actions["ce.info"] = true
	actions["ce.pause"] = true
	actions["ce.resume"] = true
	actions["ce.submit"] = true
	actions["ce.task"] = true
	actions["components.app"] = true
	actions["components.search"] = true
	actions["components.show"] = true
	actions["components.tree"] = true
	actions["duplications.show"] = true
	actions["projectdump.export"] = true
	actions["projectdump.status"] = true
	actions["navigation.component"] = true
	actions["navigation.global"] = true
	loadSkipActionsPart6(actions)
}

func loadSkipActionsPart6(actions map[string]bool) {
	actions["qualitygates.create"] = true
	actions["qualitygates.create_condition"] = true
	actions["qualitygates.project_status"] = true
	actions["qualitygates.add_group"] = true
	actions["qualitygates.add_user"] = true
	actions["qualitygates.copy"] = true
	actions["qualitygates.deselect"] = true
	actions["qualitygates.remove_group"] = true
	actions["qualitygates.remove_user"] = true
	actions["qualitygates.rename"] = true
	actions["qualitygates.search"] = true
	actions["qualitygates.search_groups"] = true
	actions["qualitygates.search_users"] = true
	actions["qualitygates.select"] = true
	actions["qualitygates.set_as_default"] = true
	actions["qualitygates.show"] = true
	actions["qualitygates.update_condition"] = true
	actions["qualitygates.delete_condition"] = true
	actions["qualitygates.get_by_project"] = true
	actions["qualitygates.list"] = true
	loadSkipActionsPart6b(actions)
}

func loadSkipActionsPart6b(actions map[string]bool) {
	actions["qualityprofiles.create"] = true
	actions["qualityprofiles.activate_rule"] = true
	actions["qualityprofiles.activate_rules"] = true
	actions["qualityprofiles.add_group"] = true
	actions["qualityprofiles.add_project"] = true
	actions["qualityprofiles.add_user"] = true
	actions["qualityprofiles.backup"] = true
	actions["qualityprofiles.change_parent"] = true
	actions["qualityprofiles.changelog"] = true
	actions["qualityprofiles.compare"] = true
	actions["qualityprofiles.copy"] = true
	actions["qualityprofiles.deactivate_rule"] = true
	actions["qualityprofiles.deactivate_rules"] = true
	actions["qualityprofiles.exporters"] = true
	actions["qualityprofiles.export"] = true
	actions["qualityprofiles.inheritance"] = true
	actions["qualityprofiles.projects"] = true
	actions["qualityprofiles.remove_group"] = true
	actions["qualityprofiles.remove_project"] = true
	actions["qualityprofiles.remove_user"] = true
	actions["qualityprofiles.rename"] = true
	actions["qualityprofiles.restore"] = true
	actions["qualityprofiles.search"] = true
	actions["qualityprofiles.search_groups"] = true
	actions["qualityprofiles.search_users"] = true
	actions["qualityprofiles.set_default"] = true
	actions["qualityprofiles.show"] = true
	loadSkipActionsPart7(actions)
}

func loadSkipActionsPart7(actions map[string]bool) {
	actions["dismissmessage.check"] = true
	actions["dismissmessage.dismiss"] = true
	actions["settings.set"] = true
	actions["settings.reset"] = true
	actions["settings.values"] = true
	actions["settings.list_definitions"] = true
	actions["settings.check_secret_key"] = true
	actions["settings.generate_secret_key"] = true
	actions["users.create"] = true
	actions["users.change_password"] = true
	actions["users.dismiss_notice"] = true
	actions["users.set_homepage"] = true
	actions["users.search"] = true
	actions["users.groups"] = true
	actions["users.update_identity_provider"] = true
	actions["users.update_login"] = true
	actions["users.current"] = true
	actions["users.identity_providers"] = true
	actions["usertokens.generate"] = true
	actions["usertokens.revoke"] = true
	actions["usertokens.search"] = true
	actions["usergroups.create"] = true
	actions["usergroups.add_user"] = true
	actions["usergroups.remove_user"] = true
	actions["usergroups.update"] = true
	actions["usergroups.search"] = true
	actions["usergroups.users"] = true
	actions["favorites.add"] = true
	actions["favorites.remove"] = true
	actions["favorites.search"] = true
	loadSkipActionsPart8(actions)
}

func loadSkipActionsPart8(actions map[string]bool) {
	actions["permissions.add_group"] = true
	actions["permissions.add_group_to_template"] = true
	actions["permissions.add_project_creator_to_template"] = true
	actions["permissions.add_user"] = true
	actions["permissions.add_user_to_template"] = true
	actions["permissions.apply_template"] = true
	actions["permissions.bulk_apply_template"] = true
	actions["permissions.create_template"] = true
	actions["permissions.groups"] = true
	actions["permissions.remove_group"] = true
	actions["permissions.remove_group_from_template"] = true
	actions["permissions.remove_project_creator_from_template"] = true
	actions["permissions.remove_user"] = true
	actions["permissions.remove_user_from_template"] = true
	actions["permissions.search_global_groups"] = true
	actions["permissions.search_project_groups"] = true
	actions["permissions.search_templates"] = true
	actions["permissions.set_default_template"] = true
	actions["permissions.update_template"] = true
	actions["permissions.users"] = true
	actions["permissions.template_groups"] = true
	actions["permissions.template_users"] = true
	actions["projects.create"] = true
	actions["projects.search"] = true
	actions["projects.search_my_projects"] = true
	actions["projects.search_my_scannable_projects"] = true
	actions["projects.update_default_visibility"] = true
	actions["projects.update_key"] = true
	actions["projects.update_visibility"] = true
	actions["projects.export_findings"] = true
	actions["projecttags.set"] = true
	actions["projecttags.search"] = true
	actions["projectbadges.measure"] = true
	actions["projectbadges.quality_gate"] = true
	actions["projectbadges.token"] = true
	actions["projectbadges.renew_token"] = true
	loadSkipActionsPart9(actions)
}

func loadSkipActionsPart9(actions map[string]bool) {
	actions["rules.app"] = true
	actions["rules.create"] = true
	actions["rules.delete"] = true
	actions["rules.repositories"] = true
	actions["rules.search"] = true
	actions["rules.show"] = true
	actions["rules.tags"] = true
	actions["rules.update"] = true
	actions["newcodeperiods.list"] = true
	actions["newcodeperiods.set"] = true
	actions["newcodeperiods.show"] = true
	actions["newcodeperiods.unset"] = true
	actions["notifications.add"] = true
	actions["notifications.list"] = true
	actions["notifications.remove"] = true
	actions["webhooks.create"] = true
	actions["webhooks.delete"] = true
	actions["webhooks.deliveries"] = true
	actions["webhooks.delivery"] = true
	actions["webhooks.list"] = true
	actions["webhooks.update"] = true
	loadSkipActionsPart10(actions)
}

func loadSkipActionsPart10(actions map[string]bool) {
	actions["webservices.list"] = true
	actions["webservices.response_example"] = true
	actions["developers.search_events"] = true
	actions["l10n.index"] = true
	actions["languages.list"] = true
	actions["metrics.search"] = true
	actions["metrics.types"] = true
	actions["server.version"] = true
	actions["features.list"] = true
	actions["plugins.available"] = true
	actions["plugins.download"] = true
	actions["plugins.installed"] = true
	actions["plugins.pending"] = true
	actions["plugins.updates"] = true
	actions["authentication.login"] = true
	actions["authentication.logout"] = true
	actions["authentication.validate"] = true
	actions["system.db_migration_status"] = true
	actions["system.health"] = true
	actions["system.info"] = true
	actions["system.migrate_db"] = true
	actions["system.ping"] = true
	actions["system.status"] = true
	actions["system.upgrades"] = true
}
