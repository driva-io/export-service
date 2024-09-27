package writers

import (
	"export-service/internal/core/domain"
	"export-service/internal/core/ports"
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"log"
	"os"
)

type ExcelWriter struct {
}

var _ ports.DataWriter = (*ExcelWriter)(nil)

func (e *ExcelWriter) Write(data []map[string]any, spec domain.PresentationSpec) (string, error) {
	wb := xlsx.NewFile()

	options := spec.GetOrderedSheetOptions()
	for _, sheetOption := range options {
		addSheet(data, wb, sheetOption)
	}

	dir, err := os.MkdirTemp("", "sheets")
	if err != nil {
		log.Println("Error creating temp dir", err)
		return "", err
	}

	path := fmt.Sprintf("%s/sheets.xlsx", dir)
	err = wb.Save(path)
	if err != nil {
		log.Println("Error saving file", err)
		return "", err
	}
	return path, nil
}

func addSheet(data []map[string]any, wb *xlsx.File, sheetOption domain.PresentationSpecSheetOptions) {
	sh, err := wb.AddSheet(sheetOption.Key)
	if err != nil {
		log.Println("Error adding sheet", sheetOption.Key, err)
		return
	}

	addHeaders(sh, sheetOption)

	insertData(data, sh, sheetOption)

	for i := range sheetOption.ActiveColumns {
		err = sh.SetColAutoWidth(i+1, xlsx.DefaultAutoWidth)
		if err != nil {
			log.Println("Error setting column width", err)
		}
	}
}

func insertData(data []map[string]any, sh *xlsx.Sheet, sheetOption domain.PresentationSpecSheetOptions) {
	for _, d := range data {
		values, ok := d[sheetOption.Key] // verifica se o campo com os valores da aba atual existe
		if !ok {
			continue
		}

		if !sheetOption.ShouldExplode {
			insertDataAsMap(values, sh, sheetOption)
		} else {
			insertDataAsList(values, sh, sheetOption)
		}
	}
}

func insertDataAsMap(values any, sh *xlsx.Sheet, sheetOption domain.PresentationSpecSheetOptions) {
	if valMap, ok := values.(map[string]any); ok {
		addRow(valMap, sh, sheetOption.ActiveColumns)
	} else {
		log.Println("Wrong type for key", sheetOption.Key)
	}
}

func insertDataAsList(values any, sh *xlsx.Sheet, sheetOption domain.PresentationSpecSheetOptions) {
	if valList, ok := values.([]map[string]any); ok {
		for _, v := range valList {
			addRow(v, sh, sheetOption.ActiveColumns)
		}
	} else {
		log.Println("Wrong type for key", sheetOption.Key)
	}
}

func addRow(data map[string]any, sh *xlsx.Sheet, activeColumns []string) {
	row := sh.AddRow()
	for _, c := range activeColumns {
		cell := row.AddCell() // adiciona célula mesmo que não tenha o valor para pular a coluna
		if value, ok := data[c]; ok {
			cell.Value = fmt.Sprintf("%v", value)
		}
		cell.SetStyle(cellStyle)
	}
}

func addHeaders(sh *xlsx.Sheet, sheetOption domain.PresentationSpecSheetOptions) {
	headers := sh.AddRow()
	headers.SetHeight(20)
	for _, c := range sheetOption.ActiveColumns {
		cell := headers.AddCell()
		cell.Value = c
		cell.SetStyle(headerStyle)
	}
}

var headerStyle = &xlsx.Style{
	Font:           xlsx.Font{Color: "FFFFFFFF", Bold: true, Size: 12, Name: "Calibri", Family: 2},
	Alignment:      xlsx.Alignment{Horizontal: "centerContinuous", Vertical: "center"},
	Fill:           xlsx.Fill{FgColor: "FF407AD6", PatternType: "solid"},
	Border:         xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"},
	ApplyFill:      true,
	ApplyAlignment: true,
	ApplyFont:      true,
	ApplyBorder:    true,
}

var cellStyle = &xlsx.Style{
	Font: xlsx.Font{Color: "FF000000", Size: 11, Name: "Calibri", Family: 2},
}
