package algo

import (
	"fmt"
	"os"
	"time"
	"github.com/xuri/excelize/v2"
)

// SaveOSROExcel 會把多種封裝方法的 O1~O4 目標值附時間戳寫進 Excel。
// 如果檔案不存在就建立，否則在最後一列後面繼續追加。
func SaveOSROExcel(file string, results []*Result) {
	const sheet = "history"

	var f *excelize.File
	var err error

	// ① 開啟或建立 Excel
	if _, err = os.Stat(file); err == nil {
		f, err = excelize.OpenFile(file)
		if err != nil {
			fmt.Println("open:", err)
			return
		}
	} else {
		f = excelize.NewFile()
		f.NewSheet(sheet)
		// 表頭
		headers := []string{
			"TimeStamp", "Method",
			"O1_Failed_TSN_CAN2TSN", "O2_Failed_AVB",
			"O3_Bytes", "O4_WCD_us", "Cost",
		}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
		}
	}

	// ② 取得最後一列
	rows, _ := f.GetRows(sheet)
	startRow := len(rows) + 1
	now := time.Now().Format("2006-01-02 15:04:05")

	// ③ 逐筆寫入
	for _, r := range results {
		row := []interface{}{
			now,
			r.Method,
			r.Obj[0], // O1-TSN Area
			r.Obj[1], // O2
			r.Obj[2], // O3
			r.Obj[3], // O4
			r.Cost,   // 你有存 Cost
		}
		ref := fmt.Sprintf("A%d", startRow)
		f.SetSheetRow(sheet, ref, &row)
		startRow++
	}

	// ④ 儲存
	if err = f.SaveAs(file); err != nil {
		fmt.Println("save:", err)
	}
}