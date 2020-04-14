package excel

import (
	"encoding/json"

	"bangseller.com/lib/exception"
	"github.com/tealeg/xlsx"
)

type Excel struct {
	*xlsx.File
	header      []string //  Excel的索引
	cSheetIndex int      //  当前Excel的索引
}

// New用于新建立一个Excle
func New(fileName string) *Excel {
	return new(Excel)
}

// ParseHeader 用于设置 获取excel中文件的头部
func (e *Excel) ParseHeader(row int) {
	e.header = e.GetRow(row)
}

// GetAllRows 用于获取所有的行的数据
func (e *Excel) GetAllRows() (ret [][]string) {
	curSheet := e.File.Sheets[e.cSheetIndex]
	for _, row := range curSheet.Rows {
		var inst []string
		for _, cell := range row.Cells {
			inst = append(inst, cell.Value)
		}
		ret = append(ret, inst)
	}
	return ret
}

// GetHeader 用于获取Excle zh Header
func (e *Excel) GetHeader() []string {
	return e.header
}

// SetSheetIndex 用于设置excle的索引
func (e *Excel) SetSheetIndex(index int) {
	e.cSheetIndex = index
}

// GetRowMap 根据 Header 为 key 值 获取行单元的map
func (e *Excel) GetRowMap(row int) map[string]string {
	m, curSheet := map[string]string{}, e.File.Sheets[e.cSheetIndex]
	for i, cell := range curSheet.Row(row).Cells {
		m[e.header[i]] = cell.Value
	}
	return m
}

// GetRowMaps 用于获取所有的行
func (e *Excel) GetRowMaps() []map[string]string {
	var m []map[string]string
	curSheet := e.File.Sheets[e.cSheetIndex]
	for index := 0; index < len(curSheet.Rows); index++ {
		m = append(m, e.GetRowMap(index))
	}
	return m
}

// GetRowsJsonStructs 用于转换Excle到结构体数组
func (e *Excel) GetRowsJsonStructs(pv interface{}) {
	ms := e.GetRowMaps()
	data, err := json.Marshal(ms)
	exception.CheckError(err)
	err = json.Unmarshal(data, pv)
	exception.CheckError(err)
}

// GetRowStruct 获取行的结构体
func (e *Excel) GetRowStruct(row int, s interface{}) {
	m := e.GetRowMap(row)
	data, err := json.Marshal(m)
	exception.CheckError(err)
	err = json.Unmarshal(data, s)
	exception.CheckError(err)
}

// GetRow 获取行单元数组
func (e *Excel) GetRow(row int) []string {
	values := []string{}
	for _, cell := range e.Sheets[e.cSheetIndex].Row(row).Cells {
		values = append(values, cell.Value)
	}
	return values
}
