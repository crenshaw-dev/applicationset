package generators

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/argoproj-labs/applicationset/api/v1alpha1"
)

func TestGetRelevantGenerators(t *testing.T) {
	requestedGenerator := &v1alpha1.ApplicationSetGenerator{
		ApplicationSetTerminalGenerator: &v1alpha1.ApplicationSetTerminalGenerator{
			List: &v1alpha1.ListGenerator{},
		},
	}
	allGenerators := map[string]Generator{
		"List": NewListGenerator(),
	}
	relevantGenerators := GetRelevantGenerators(requestedGenerator, allGenerators)

	for _, generator := range relevantGenerators {
		if generator == nil {
			t.Fatal(`GetRelevantGenerators produced a nil generator`)
		}
	}

	numRelevantGenerators := len(relevantGenerators)
	if numRelevantGenerators != 1 {
		t.Fatalf(`GetRelevantGenerators produced %d generators instead of the expected 1`, numRelevantGenerators)
	}
}

func TestNoGeneratorNilReferenceError(t *testing.T) {
	generators := []Generator{
		&ClusterGenerator{},
		&DuckTypeGenerator{},
		&GitGenerator{},
		&ListGenerator{},
		&MatrixGenerator{},
		&MergeGenerator{},
		&PullRequestGenerator{},
		&SCMProviderGenerator{},
	}

	applicationSetGenerator := &v1alpha1.ApplicationSetGenerator{
		ApplicationSetTerminalGenerator: nil,
		Matrix:                          nil,
		Merge:                           nil,
	}

	for _, generator := range generators {
		testCaseCopy := generator // since tests may run in parallel

		generatorName := reflect.TypeOf(testCaseCopy).Elem().Name()
		t.Run(fmt.Sprintf("%s does not throw a nil reference error when ApplicationSetTerminalGenerator is nil", generatorName), func(t *testing.T) {
			_, err := generator.GenerateParams(applicationSetGenerator, &v1alpha1.ApplicationSet{})

			assert.NoError(t, err)
		})
	}
}