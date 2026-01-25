package sonargo

const (
	// MaxPageSize is the maximum allowed page size for pagination.
	MaxPageSize = 500
	// MinPageSize is the minimum allowed page size for pagination.
	MinPageSize = 1

	// MaxLinkNameLength is the maximum length for a project link name.
	MaxLinkNameLength = 128
	// MaxLinkURLLength is the maximum length for a project link URL.
	MaxLinkURLLength = 2048
	// MaxTokenNameLength is the maximum length for a user token name.
	MaxTokenNameLength = 100
	// MaxBranchNameLength is the maximum length for a branch name.
	MaxBranchNameLength = 255
)

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedLanguages is the set of supported programming languages.
	allowedLanguages = map[string]struct{}{
		"azureresourcemanager": {},
		"cloudformation":       {},
		"cs":                   {},
		"css":                  {},
		"docker":               {},
		"flex":                 {},
		"go":                   {},
		"ipynb":                {},
		"java":                 {},
		"js":                   {},
		"json":                 {},
		"jsp":                  {},
		"kotlin":               {},
		"kubernetes":           {},
		"php":                  {},
		"py":                   {},
		"ruby":                 {},
		"rust":                 {},
		"scala":                {},
		"secrets":              {},
		"terraform":            {},
		"text":                 {},
		"ts":                   {},
		"vbnet":                {},
		"web":                  {},
		"xml":                  {},
		"yaml":                 {},
	}

	// allowedSeverities is the set of supported severity levels.
	allowedSeverities = map[string]struct{}{
		"BLOCKER":  {},
		"CRITICAL": {},
		"MAJOR":    {},
		"MINOR":    {},
		"INFO":     {},
	}

	// allowedImpactSeverities is the set of supported impact severity levels.
	allowedImpactSeverities = map[string]struct{}{
		"BLOCKER": {},
		"HIGH":    {},
		"MEDIUM":  {},
		"LOW":     {},
		"INFO":    {},
	}

	// allowedCleanCodeAttributesCategories is the set of supported Clean Code attribute categories.
	allowedCleanCodeAttributesCategories = map[string]struct{}{
		"ADAPTABLE":   {},
		"CONSISTENT":  {},
		"INTENTIONAL": {},
		"RESPONSIBLE": {},
	}

	// allowedCleanCodeAttributes is the set of supported Clean Code attributes.
	allowedCleanCodeAttributes = map[string]struct{}{
		"CONVENTIONAL": {},
		"FORMATTED":    {},
		"IDENTIFIABLE": {},
		"CLEAR":        {},
		"COMPLETE":     {},
		"EFFICIENT":    {},
		"LOGICAL":      {},
		"DISTINCT":     {},
		"FOCUSED":      {},
		"MODULAR":      {},
		"TESTED":       {},
		"LAWFUL":       {},
		"RESPECTFUL":   {},
		"TRUSTWORTHY":  {},
	}

	// allowedImpactSoftwareQualities is the set of supported impact software qualities.
	allowedImpactSoftwareQualities = map[string]struct{}{
		"MAINTAINABILITY": {},
		"RELIABILITY":     {},
		"SECURITY":        {},
	}

	// allowedInheritanceTypes is the set of supported inheritance types.
	allowedInheritanceTypes = map[string]struct{}{
		"NONE":       {},
		"INHERITED":  {},
		"OVERRIDDES": {},
	}

	// allowedOwaspCategories is the set of supported OWASP categories.
	allowedOwaspCategories = map[string]struct{}{
		"a1":  {},
		"a2":  {},
		"a3":  {},
		"a4":  {},
		"a5":  {},
		"a6":  {},
		"a7":  {},
		"a8":  {},
		"a9":  {},
		"a10": {},
	}

	// allowedOwaspMobileCategories is the set of supported OWASP Mobile categories.
	allowedOwaspMobileCategories = map[string]struct{}{
		"m1":  {},
		"m2":  {},
		"m3":  {},
		"m4":  {},
		"m5":  {},
		"m6":  {},
		"m7":  {},
		"m8":  {},
		"m9":  {},
		"m10": {},
	}

	// allowedRulesStatuses is the set of supported statuses.
	allowedRulesStatuses = map[string]struct{}{
		"READY":      {},
		"DEPRECATED": {},
		"REMOVED":    {},
		"BETA":       {},
	}

	// allowedRulesExistingStatuses is the set of supported existing statuses.
	allowedRulesExistingStatuses = map[string]struct{}{
		"READY":      {},
		"DEPRECATED": {},
		"BETA":       {},
	}

	// allowedRulesTypes is the set of supported rule types.
	allowedRulesTypes = map[string]struct{}{
		"CODE_SMELL":       {},
		"BUG":              {},
		"VULNERABILITY":    {},
		"SECURITY_HOTSPOT": {},
	}

	// allowedSansTop25Categories is the set of supported SANS Top 25 categories.
	allowedSansTop25Categories = map[string]struct{}{
		"insecure-interaction": {},
		"risky-resource":       {},
		"porous-defenses":      {},
	}

	// allowedSelectedFilters is the set of supported selected filters.
	allowedSelectedFilters = map[string]struct{}{
		"all":        {},
		"selected":   {},
		"deselected": {},
	}

	// allowedIssueTypes is the set of supported issue types.
	allowedIssueTypes = map[string]struct{}{
		"CODE_SMELL":       {},
		"BUG":              {},
		"VULNERABILITY":    {},
		"SECURITY_HOTSPOT": {},
	}

	// allowedIssueTransitions is the set of supported issue transitions.
	allowedIssueTransitions = map[string]struct{}{
		"confirm":           {},
		"unconfirm":         {},
		"reopen":            {},
		"resolve":           {},
		"falsepositive":     {},
		"wontfix":           {},
		"accept":            {},
		"close":             {},
		"resolveasreviewed": {},
		"resetastoreview":   {},
	}

	// allowedIssueStatuses is the set of supported issue statuses.
	allowedIssueStatuses = map[string]struct{}{
		"OPEN":           {},
		"CONFIRMED":      {},
		"FALSE_POSITIVE": {},
		"ACCEPTED":       {},
		"FIXED":          {},
		"IN_SANDBOX":     {},
	}

	// allowedIssueResolutions is the set of supported issue resolutions.
	allowedIssueResolutions = map[string]struct{}{
		"FIXED":          {},
		"REMOVED":        {},
		"FALSE-POSITIVE": {},
		"WONTFIX":        {},
	}

	// allowedIssueScopes is the set of supported issue scopes.
	allowedIssueScopes = map[string]struct{}{
		"MAIN": {},
		"TEST": {},
	}
)

// PaginationArgs contains common pagination parameters for API requests.
type PaginationArgs struct {
	// Page is the response page number. Must be strictly greater than 0.
	Page int64 `url:"p,omitempty"`
	// PageSize is the response page size. Must be greater than 0 and less than or equal to 500.
	PageSize int64 `url:"ps,omitempty"`
}

// Validate validates the pagination arguments.
func (p *PaginationArgs) Validate() error {
	return ValidatePagination(p.Page, p.PageSize)
}
