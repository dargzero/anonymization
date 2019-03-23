package anondb

import (
	"github.com/dargzero/anonymization/anonmodel"
	"log"
)

//MakeTypeConversionTable   asdasd
func MakeTypeConversionTable(datasetName string) (table map[string]anonmodel.TypeConversionFunc, err error) {
	table = make(map[string]anonmodel.TypeConversionFunc, 0)
	dataset, e := GetDataset(datasetName)
	log.Println(e)
	for _, field := range dataset.Fields {
		if field.Type == "coords" {
			table[field.Name] = PreprocessCoord
		}
		if err != nil {
			return
		}
	}
	return
}
