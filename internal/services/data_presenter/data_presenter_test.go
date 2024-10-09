package data_presenter_test

import (
	"encoding/json"
	"export-service/internal/core/domain"
	"export-service/internal/services/data_presenter"
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

func getTestData(file_path string) map[string]any {
	dataBytes, err := os.ReadFile(file_path)
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
func getArrayTestData(file_path string) []map[string]any {
	dataBytes, err := os.ReadFile(file_path)
	if err != nil {
		panic(err)
	}
	var data []map[string]any
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		panic(err)
	}
	return data
}

func getTestCases(file_path string) []TestCase {
	casesBytes, err := os.ReadFile(file_path)
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
	data := getTestData("testdata/company.json")
	testCases := getTestCases("testdata/test_presenter.yaml")

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			presented, err := data_presenter.PresentSingle(data, tc.SpecValue)
			require.NoError(t, err)
			require.Equal(t, tc.Expected, presented)
		})
	}
}

func TestPresentLinkedin(t *testing.T) {
	data := getArrayTestData("testdata/linkedin.json")
	testCases := getTestCases("testdata/linkedin_test_presenter.yaml")

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			presented, err := data_presenter.PresentMultiple(data, tc.SpecValue)
			require.NoError(t, err)
			require.Equal(t, tc.Expected, presented)
		})
	}
}
