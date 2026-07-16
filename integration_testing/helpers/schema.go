package helpers

import (
	"slices"
	"sync"

	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

//nolint:gochecknoglobals // process-wide recorder shared by every e2e client created via NewClient
var (
	schemaMismatchesMu   sync.Mutex
	schemaMismatchesSeen = make(map[string]struct{})
	schemaMismatches     []sonar.SchemaMismatch
)

// RecordSchemaMismatches is a sonar.SchemaObserver that accumulates every
// distinct SchemaMismatch observed across the whole e2e run. It is wired into
// every client returned by NewClient/NewDefaultClient so that any API call
// made by any e2e test contributes to strict schema validation, without each
// test needing to assert on it individually. Retrieve the accumulated results
// with SchemaMismatches.
func RecordSchemaMismatches(_ string, mismatches []sonar.SchemaMismatch) {
	if len(mismatches) == 0 {
		return
	}

	schemaMismatchesMu.Lock()
	defer schemaMismatchesMu.Unlock()

	for _, mismatch := range mismatches {
		key := mismatch.GoType + "|" + mismatch.Path
		if _, seen := schemaMismatchesSeen[key]; seen {
			continue
		}

		schemaMismatchesSeen[key] = struct{}{}

		schemaMismatches = append(schemaMismatches, mismatch)
	}
}

// SchemaMismatches returns every distinct SchemaMismatch recorded so far
// across the whole e2e run.
func SchemaMismatches() []sonar.SchemaMismatch {
	schemaMismatchesMu.Lock()
	defer schemaMismatchesMu.Unlock()

	return slices.Clone(schemaMismatches)
}
