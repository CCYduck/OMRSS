package flow

import (
    "os"
    "time"
	"fmt"
	"github.com/xuri/excelize/v2"
)


func SaveExcel(file string, encaps []*Method) {
    const sheet = "history"

    var f *excelize.File
    var err error

    // 確保 output 資料夾存在
	outputDir := "En_output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			fmt.Println("mkdir:", err)
			return
		}
	}
	fullPath := outputDir + "/" + file

    // 檔案存在就開啟，不存在就新建
    if _, err = os.Stat(fullPath ); err == nil {
        f, err = excelize.OpenFile(fullPath )
        if err != nil { fmt.Println("open:", err); return }
    } else {
        f = excelize.NewFile()
        f.NewSheet(sheet)
        // 表頭
        headers := []string{
            "TimeStamp", "Method", "TotalFlows",
            "StreamBytes", "TSN_StreamCount","O1_Encap_Drop","O1_Decap_Drop", "Delay_ms",
        }
        for i, h := range headers {
            cell, _ := excelize.CoordinatesToCellName(i+1, 1)
            f.SetCellValue(sheet, cell, h)
        }
    }

    // 找目前最後一列
    rows, _ := f.GetRows(sheet)
    startRow := len(rows) + 1

    now := time.Now().Format("2006-01-02 15:04:05")
    for _, m := range encaps {
        row := []interface{}{
            now,
            m.Method_Name,
            len(m.CAN2TSNFlows),
            int(m.BytesSent),
            m.TSNFrameCount,
            m.CAN2TSN_O1_Drop,
            m.CAN_Area_O1_Drop,
            m.CAN2TSN_Delay.Seconds() * 1000,
        }
        f.SetSheetRow(sheet, fmt.Sprintf("A%d", startRow), &row)
        startRow++
    }

    if err = f.SaveAs(fullPath ); err != nil {
        fmt.Println("save:", err)
    }
}