package memorizer

// import(
// 	"os"
// 	"encoding/csv"
// 	"strconv"
// 	"fmt"

// )
	


// type Result struct {
//     TimeStamp      string  // 例如 "2025-05-06 15:04:05"
//     Method         string
//     TotalFlows     int
//     StreamSize     int
//     StreamCount    int
//     O1Drop         int
//     DelayMs        float64
// }

// func appendCSV(res []Result, csvFile string) error {
//     // 如果檔案不存在就連同標頭一起建立
//     needHeader := false
//     if _, err := os.Stat(csvFile); os.IsNotExist(err) {
//         needHeader = true
//     }

//     f, err := os.OpenFile(csvFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
//     if err != nil { return err }
//     defer f.Close()

//     w := csv.NewWriter(f)
//     if needHeader {
//         w.Write([]string{
//             "TimeStamp","Method","TotalFlows","StreamSize",
//             "StreamCount","O1Drop","DelayMs",
//         })
//     }
//     for _, r := range res {
//         w.Write([]string{
//             r.TimeStamp, r.Method,
//             strconv.Itoa(r.TotalFlows),
//             strconv.Itoa(r.StreamSize),
//             strconv.Itoa(r.StreamCount),
//             strconv.Itoa(r.O1Drop),
//             fmt.Sprintf("%.6f", r.DelayMs),
//         })
//     }
//     w.Flush()
//     return w.Error()
// }