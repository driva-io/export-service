package writers

import (
	"export-service/internal/core/domain"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tealeg/xlsx/v3"
)

func TestExcelWriter_Write(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		data := []map[string]any{{
			"RFB": map[string]any{
				"CNPJ":   "111111",
				"Titulo": "RAZAO SOCIAL",
			},
			"Telefones": []map[string]any{
				{"CNPJ": "111111", "Telefone": "123456", "WhatsApp": "SIM"},
				{"CNPJ": "111111", "Telefone": "564565", "WhatsApp": "SIM"},
			},
		}, {
			"RFB": map[string]any{
				"CNPJ":   "222222",
				"Titulo": "RAZAO SOCIAL 222",
			},
			"Telefones": []map[string]any{
				{"CNPJ": "222222", "Telefone": "9999999", "WhatsApp": "SIM"},
			},
		}}

		spec := domain.PresentationSpec{
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "RFB",
					ActiveColumns: []string{"CNPJ", "Titulo"},
					Position:      1,
					ShouldExplode: false,
				},
				{
					Key:           "Telefones",
					ActiveColumns: []string{"CNPJ", "Telefone", "WhatsApp"},
					Position:      2,
					ShouldExplode: true,
				},
			},
		}

		ew := ExcelWriter{}
		path, err := ew.Write(data, spec)
		require.NoError(t, err)
		assert.NotEmpty(t, path)

		wb, err := xlsx.OpenFile(path)
		require.NoError(t, err)

		expectedSheets := []string{"RFB", "Telefones"}
		assert.Lenf(t, wb.Sheets, len(expectedSheets), "Expected %d sheets, got %d", len(expectedSheets), len(wb.Sheets))
		for i, sh := range wb.Sheets {
			assert.Equalf(t, expectedSheets[i], sh.Name, "Sheet %d should be %s", i, expectedSheets[i])
		}

		err = wb.Sheet["RFB"].ForEachRow(func(r *xlsx.Row) error {
			if r.GetCoordinate() == 0 {
				assert.Equalf(t, "CNPJ", r.GetCell(0).Value, "Expected 'CNPJ' as first column header, got %s", r.GetCell(0).Value)
				assert.Equalf(t, "Titulo", r.GetCell(1).Value, "Expected 'Titulo' as second column header, got %s", r.GetCell(1).Value)
				return nil
			}

			expected := data[r.GetCoordinate()-1]["RFB"].(map[string]any)
			assert.Equal(t, expected["CNPJ"], r.GetCell(0).Value)
			assert.Equal(t, expected["Titulo"], r.GetCell(1).Value)
			return nil
		})
		require.NoError(t, err)

		header, err := wb.Sheet["Telefones"].Row(0)
		require.NoError(t, err)
		assert.Equalf(t, "CNPJ", header.GetCell(0).Value, "Expected 'CNPJ' as first column header, got %s", header.GetCell(0).Value)
		assert.Equalf(t, "Telefone", header.GetCell(1).Value, "Expected 'Telefone' as second column header, got %s", header.GetCell(1).Value)
		assert.Equalf(t, "WhatsApp", header.GetCell(2).Value, "Expected 'WhatsApp' as third column header, got %s", header.GetCell(2).Value)

		expectedPhones := []string{"123456", "564565", "9999999"}
		err = wb.Sheet["Telefones"].ForEachRow(func(r *xlsx.Row) error {
			if r.GetCoordinate() == 0 {
				return nil
			}

			assert.Equal(t, expectedPhones[r.GetCoordinate()-1], r.GetCell(1).Value)
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("Should order sheets positions", func(t *testing.T) {
		spec := domain.PresentationSpec{
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "Telefones",
					ActiveColumns: []string{"CNPJ", "Telefone", "WhatsApp"},
					Position:      2,
				},
				{
					Key:           "RFB",
					ActiveColumns: []string{"CNPJ", "Titulo"},
					Position:      1,
				},
			},
		}

		ew := ExcelWriter{}
		path, err := ew.Write([]map[string]any{}, spec)
		require.NoError(t, err)
		assert.NotEmpty(t, path)

		wb, err := xlsx.OpenFile(path)
		require.NoError(t, err)

		expectedSheets := []string{"RFB", "Telefones"}
		assert.Lenf(t, wb.Sheets, len(expectedSheets), "Expected %d sheets, got %d", len(expectedSheets), len(wb.Sheets))
		for i, sh := range wb.Sheets {
			assert.Equalf(t, expectedSheets[i], sh.Name, "Sheet %d should be %s", i, expectedSheets[i])
		}
	})

	t.Run("Should handle duplicated positions", func(t *testing.T) {
		spec := domain.PresentationSpec{
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "Emails",
					ActiveColumns: []string{"CNPJ", "Email", "WhatsApp"},
					Position:      3,
				},
				{
					Key:           "RFB",
					ActiveColumns: []string{"CNPJ", "Titulo"},
					Position:      1,
					ShouldExplode: false,
				},
				{
					Key:           "Telefones",
					ActiveColumns: []string{"CNPJ", "Telefone", "WhatsApp"},
					Position:      1,
					ShouldExplode: true,
				},
			},
		}

		ew := ExcelWriter{}
		path, err := ew.Write([]map[string]any{}, spec)
		require.NoError(t, err)
		assert.NotEmpty(t, path)

		wb, err := xlsx.OpenFile(path)
		require.NoError(t, err)

		expectedSheets := []string{"RFB", "Telefones", "Emails"}
		assert.Lenf(t, wb.Sheets, len(expectedSheets), "Expected %d sheets, got %d", len(expectedSheets), len(wb.Sheets))
		for i, sh := range wb.Sheets {
			assert.Equalf(t, expectedSheets[i], sh.Name, "Sheet %d should be %s", i, expectedSheets[i])
		}
	})

	t.Run("Should handle empty values for sheets", func(t *testing.T) {
		data := []map[string]any{{
			"RFB": map[string]any{
				"CNPJ":   "111111",
				"Titulo": "RAZAO SOCIAL",
			},
		}, {
			"RFB": map[string]any{
				"CNPJ":   "222222",
				"Titulo": "RAZAO SOCIAL 222",
			},
		}}

		spec := domain.PresentationSpec{
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "RFB",
					ActiveColumns: []string{"CNPJ", "Titulo"},
					Position:      1,
					ShouldExplode: false,
				},
				{
					Key:           "Telefones",
					ActiveColumns: []string{"CNPJ", "Telefone"},
					Position:      2,
					ShouldExplode: true,
				},
			},
		}

		ew := ExcelWriter{}
		path, err := ew.Write(data, spec)
		require.NoError(t, err)
		assert.NotEmpty(t, path)

		wb, err := xlsx.OpenFile(path)
		require.NoError(t, err)

		expectedSheets := []string{"RFB"}
		assert.Lenf(t, wb.Sheets, len(expectedSheets), "Expected %d sheets, got %d", len(expectedSheets), len(wb.Sheets))
		for i, sh := range wb.Sheets {
			assert.Equalf(t, expectedSheets[i], sh.Name, "Sheet %d should be %s", i, expectedSheets[i])
		}
	})

	t.Run("Should handle missing fields in data", func(t *testing.T) {
		data := []map[string]any{{
			"RFB": map[string]any{
				"CNPJ": "111111",
			},
		}, {
			"RFB": map[string]any{
				"Titulo": "RAZAO SOCIAL",
			},
		}}

		spec := domain.PresentationSpec{
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "RFB",
					ActiveColumns: []string{"CNPJ", "Titulo"},
					Position:      1,
				},
			},
		}

		ew := ExcelWriter{}
		path, err := ew.Write(data, spec)
		require.NoError(t, err)
		assert.NotEmpty(t, path)
		fmt.Println(path)

		wb, err := xlsx.OpenFile(path)
		require.NoError(t, err)

		sh := wb.Sheet["RFB"]
		first, err := sh.Row(1)
		require.NoError(t, err)
		assert.Equalf(t, data[0]["RFB"].(map[string]any)["CNPJ"], first.GetCell(0).Value, "Should have CNPJ for first entry")
		assert.Emptyf(t, first.GetCell(1).Value, "Should not have Titulo for first entry")

		second, err := sh.Row(2)
		require.NoError(t, err)
		assert.Emptyf(t, second.GetCell(0).Value, "Should not have CNPJ for second entry")
		assert.Equalf(t, data[1]["RFB"].(map[string]any)["Titulo"], second.GetCell(1).Value, "Should have Titulo for second entry")
	})

	t.Run("Should handle invalid sheet names", func(t *testing.T) {
		spec := domain.PresentationSpec{
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "",
					ActiveColumns: []string{"CNPJ"},
				},
				{
					Key:           "Maximum sheet name exceeded here!!!!!!!!!!!!!!!!!",
					ActiveColumns: []string{"CNPJ"},
				},
				{
					Key:           "Invalid characters :  / ? * [ ]",
					ActiveColumns: []string{"CNPJ"},
				},
				{
					Key:           "Valid",
					ActiveColumns: []string{"CNPJ"},
				},
			},
		}

		ew := ExcelWriter{}
		path, err := ew.Write([]map[string]any{}, spec)
		require.NoError(t, err)
		assert.NotEmpty(t, path)

		wb, err := xlsx.OpenFile(path)
		require.NoError(t, err)

		expectedSheets := []string{"Valid"}
		assert.Lenf(t, wb.Sheets, len(expectedSheets), "Expected %d sheets, got %d", len(expectedSheets), len(wb.Sheets))
		for i, sh := range wb.Sheets {
			assert.Equalf(t, expectedSheets[i], sh.Name, "Sheet %d should be %s", i, expectedSheets[i])
		}
	})

	t.Run("Should handle wrong type for not exploding", func(t *testing.T) {
		data := []map[string]any{{
			"RFB": map[string]any{
				"CNPJ": "111111",
			},
		}, {
			"RFB": true,
		}, {
			"RFB": []map[string]any{
				{
					"CNPJ": "111111",
				},
			},
		}}

		spec := domain.PresentationSpec{
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "RFB",
					ActiveColumns: []string{"CNPJ"},
					Position:      1,
					ShouldExplode: false,
				},
			},
		}

		ew := ExcelWriter{}
		path, err := ew.Write(data, spec)
		require.NoError(t, err)
		assert.NotEmpty(t, path)

		wb, err := xlsx.OpenFile(path)
		require.NoError(t, err)

		expectedSheets := []string{"RFB"}
		assert.Lenf(t, wb.Sheets, len(expectedSheets), "Expected %d sheets, got %d", len(expectedSheets), len(wb.Sheets))
		for i, sh := range wb.Sheets {
			assert.Equalf(t, expectedSheets[i], sh.Name, "Sheet %d should be %s", i, expectedSheets[i])
		}

		err = wb.Sheet["RFB"].ForEachRow(func(r *xlsx.Row) error {
			if r.GetCoordinate() == 0 {
				assert.Equalf(t, "CNPJ", r.GetCell(0).Value, "Expected 'CNPJ' as first column header, got %s", r.GetCell(0).Value)
				return nil
			}

			expected := data[r.GetCoordinate()-1]["RFB"].(map[string]any)
			assert.Equal(t, expected["CNPJ"], r.GetCell(0).Value)
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("Should handle wrong type for exploding", func(t *testing.T) {
		data := []map[string]any{{
			"Telefones": []map[string]any{
				{"Telefone": "123456"},
			},
		}, {
			"Telefones": map[string]any{
				"Telefone": "123456",
			},
		}, {
			"Telefones": true,
		}}

		spec := domain.PresentationSpec{
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "Telefones",
					ActiveColumns: []string{"Telefone"},
					Position:      1,
					ShouldExplode: true,
				},
			},
		}

		ew := ExcelWriter{}
		path, err := ew.Write(data, spec)
		require.NoError(t, err)
		assert.NotEmpty(t, path)

		wb, err := xlsx.OpenFile(path)
		require.NoError(t, err)

		expectedSheets := []string{"Telefones"}
		assert.Lenf(t, wb.Sheets, len(expectedSheets), "Expected %d sheets, got %d", len(expectedSheets), len(wb.Sheets))
		for i, sh := range wb.Sheets {
			assert.Equalf(t, expectedSheets[i], sh.Name, "Sheet %d should be %s", i, expectedSheets[i])
		}

		err = wb.Sheet["Telefones"].ForEachRow(func(r *xlsx.Row) error {
			if r.GetCoordinate() == 0 {
				assert.Equalf(t, "Telefone", r.GetCell(0).Value, "Expected 'Telefone' as first column header, got %s", r.GetCell(0).Value)
				return nil
			}

			expected := data[r.GetCoordinate()-1]["Telefones"].([]map[string]any)[0]
			assert.Equal(t, expected["Telefone"], r.GetCell(0).Value)
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("Should handle not string values", func(t *testing.T) {
		data := []map[string]any{{
			"RFB": map[string]any{
				"CNPJ": "111111",
			},
		}, {
			"RFB": map[string]any{
				"CNPJ": true,
			},
		}, {
			"RFB": map[string]any{
				"CNPJ": 123,
			},
		}, {
			"RFB": map[string]any{
				"CNPJ": 1.2,
			},
		}}

		spec := domain.PresentationSpec{
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "RFB",
					ActiveColumns: []string{"CNPJ"},
					Position:      1,
					ShouldExplode: false,
				},
			},
		}

		ew := ExcelWriter{}
		path, err := ew.Write(data, spec)
		require.NoError(t, err)
		assert.NotEmpty(t, path)

		wb, err := xlsx.OpenFile(path)
		require.NoError(t, err)

		expectedSheets := []string{"RFB"}
		assert.Lenf(t, wb.Sheets, len(expectedSheets), "Expected %d sheets, got %d", len(expectedSheets), len(wb.Sheets))
		for i, sh := range wb.Sheets {
			assert.Equalf(t, expectedSheets[i], sh.Name, "Sheet %d should be %s", i, expectedSheets[i])
		}

		expectedValues := []string{"111111", "true", "123", "1.2"}
		err = wb.Sheet["RFB"].ForEachRow(func(r *xlsx.Row) error {
			if r.GetCoordinate() == 0 {
				assert.Equalf(t, "CNPJ", r.GetCell(0).Value, "Expected 'CNPJ' as first column header, got %s", r.GetCell(0).Value)
				return nil
			}

			expected := expectedValues[r.GetCoordinate()-1]
			assert.Equal(t, expected, r.GetCell(0).Value)
			return nil
		})
		require.NoError(t, err)
	})
}
