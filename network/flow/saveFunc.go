package flow

import (
    "encoding/csv"
    "os"
    "strconv"
    "time"
	"fmt"
	"github.com/xuri/excelize/v2"
)


func SaveCSV(file string, encaps []*Method) {
    needHeader := false
    if _, err := os.Stat(file); os.IsNotExist(err) {
        needHeader = true
    }
    f, _ := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    defer f.Close()
    w := csv.NewWriter(f)

    if needHeader {
        w.Write([]string{
            "TimeStamp","Method","TotalFlows",
            "StreamBytes","TSN_StreamCount","O1_Encap_Drop","O1_Decap_Drop","Delay_ms",
        })
    }
    now := time.Now().Format("2006-01-02 15:04:05")
    for _, m := range encaps {
        w.Write([]string{
            now,
            m.Method_Name,
            strconv.Itoa(len(m.CAN2TSNFlows)),
            strconv.Itoa(int(m.BytesSent)),
            strconv.Itoa(m.TSNFrameCount),
            strconv.Itoa(m.CAN2TSN_O1_Drop),
            strconv.Itoa(m.CAN_Area_O1_Drop),
            strconv.FormatFloat(m.CAN2TSN_Delay.Seconds()*1000,
                                'f', 6, 64),
        })
    }
    w.Flush()
}

func SaveExcel(file string, encaps []*Method) {
    const sheet = "history"

    var f *excelize.File
    var err error

    // 檔案存在就開啟，不存在就新建
    if _, err = os.Stat(file); err == nil {
        f, err = excelize.OpenFile(file)
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

    if err = f.SaveAs(file); err != nil {
        fmt.Println("save:", err)
    }
}