package services_test

import (
	"encoding/json"
	"export-service/internal/core/domain"
	"export-service/internal/services"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type TestCase struct {
	Name      string                      `yaml:"name"`
	SpecValue domain.PresentationSpecSpec `yaml:"spec"`
	Expected  map[string]any              `yaml:"expected"`
}

func getTestData() map[string]any {
	dataBytes, err := os.ReadFile("testdata/company.json")
	if err != nil {
		panic(err)
	}
	var data map[string]any
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		panic(err)
	}
	return data

}
func getTestCases() []TestCase {
	casesBytes, err := os.ReadFile("testdata/test_presenter.yaml")
	if err != nil {
		panic(err)
	}
	var cases []TestCase
	err = yaml.Unmarshal(casesBytes, &cases)
	if err != nil {
		panic(err)
	}
	return cases
}
func TestPresentSingle(t *testing.T) {
	data := getTestData()
	testCases := getTestCases()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			presented, err := services.PresentSingle(data, tc.SpecValue)
			require.NoError(t, err)
			require.Equal(t, tc.Expected, presented)
		})
	}
}
