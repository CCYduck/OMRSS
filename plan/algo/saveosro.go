package algo

import (
	"fmt"
	"os"
	"time"
	"github.com/xuri/excelize/v2"
)

// SaveOSROExcel 會把多種封裝方法的 O1~O4 目標值附時間戳寫進 Excel。
// 如果檔案不存在就建立，否則在最後一列後面繼續追加。
func SaveOSROExcel(file string, sp []*Result, kp []*Result) {
	name := []string{"sp", "kp"}
	var f *excelize.File
	var err error

	// 確保 output 資料夾存在
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			fmt.Println("mkdir:", err)
			return
		}
	}
	fullPath := outputDir + "/" + file

	// ① 開啟或建立 Excel
	if _, err = os.Stat(fullPath); err == nil {
		f, err = excelize.OpenFile(fullPath)
		if err != nil {
			fmt.Println("open:", err)
			return
		}
	} else {
		f = excelize.NewFile()

		f.NewSheet("sp")
		f.NewSheet("kp")

		// 表頭
		headers := []string{
			"TimeStamp", "Method",
			"O1_Failed_TSN_CAN2TSN", "O2_Failed_AVB",
			"O3_Bytes", "O4_WCD_us", "Cost",
		}
		for _, n := range name{
			for i, h := range headers {
				cell, _ := excelize.CoordinatesToCellName(i+1, 1)
				f.SetCellValue(n, cell, h)
			}
		}
		
	}


	// ② 取得最後一列
	sp_rows, _ := f.GetRows("sp")
	sp_startRow := len(sp_rows) + 1
	now := time.Now().Format("2006-01-02 15:04:05")

	// ③ 逐筆寫入
	for _, r := range sp {
		row := []interface{}{
			now,
			r.Method,
			r.Obj[0], // O1-TSN Area
			r.Obj[1], // O2
			r.Obj[2], // O3
			r.Obj[3], // O4
			r.Cost,   // 你有存 Cost
		}
		ref := fmt.Sprintf("A%d", sp_startRow)
		f.SetSheetRow("sp", ref, &row)
		sp_startRow++
	}

	kp_rows, _ := f.GetRows("kp")
	kp_startRow := len(kp_rows) + 1
	now = time.Now().Format("2006-01-02 15:04:05")

	for _, r := range kp {
		row := []interface{}{
			now,
			r.Method,
			r.Obj[0], // O1-TSN Area
			r.Obj[1], // O2
			r.Obj[2], // O3
			r.Obj[3], // O4
			r.Cost,   // 你有存 Cost
		}
		ref := fmt.Sprintf("A%d", kp_startRow)
		f.SetSheetRow("kp", ref, &row)
		kp_startRow++
	}

	// ④ 儲存
	if err = f.SaveAs(fullPath); err != nil {
		fmt.Println("save:", err)
	}
}
