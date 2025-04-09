package flow

import (
	"crypto/rand"
	"math/big"
)

func TSN_stream() *TSN {
	t_period, t_datasize := tsn_random()
	tsn := new_TSN(t_period, t_datasize)

	return tsn
}

func AVB_stream() *AVB {
	a_datasize := avb_random()
	avb := new_AVB(a_datasize)

	return avb
}

func importantCAN_stream() *importantCAN {
	importantcan := new_importantCAN()
	
	return importantcan
}

func unimportantCAN_stream() *unimportantCAN {
	uc_period,  uc_deadline := unimportantCAN_random()
	unimportantcan := new_unimportantCAN(uc_period,  uc_deadline)

	return unimportantcan
}

func tsn_random() (int, float64) {
	tsn_period_arr := []int{100, 500, 1000, 1500, 2000}
	tsn_datasize_arr := []float64{30., 40., 50., 60., 70., 80., 90., 100.}
	period_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(tsn_period_arr))))
	datasize_rng, _ := rand.Int(rand.Reader, big.NewInt(int64((len(tsn_datasize_arr)))))

	return tsn_period_arr[period_rng.Int64()], tsn_datasize_arr[datasize_rng.Int64()]
}

func avb_random() float64 {
	avb_datasize_arr := []float64{1000., 1100., 1200., 1300., 1400., 1500.}
	datasize_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(avb_datasize_arr))))

	return avb_datasize_arr[datasize_rng.Int64()]
}

//因為ImpCAN 是固定的無需做random生成
func unimportantCAN_random() (int, int) {
	unimportantCAN_period_arr := []int{50000, 100000, 150000}
	unimportantCAN_deadline := []int{10000, 12000, 14000, 16000, 18000, 20000}
	period_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(unimportantCAN_period_arr))))
	deadline_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(unimportantCAN_deadline))))

	return unimportantCAN_period_arr[period_rng.Int64()], unimportantCAN_deadline[deadline_rng.Int64()]
}


func Random_Devices(Nnode int) (int, int) {

	sourceBig, _ := rand.Int(rand.Reader, big.NewInt(int64(Nnode)))
    sourceIndex := int(sourceBig.Int64())

    var destIndex int
    for {
        destBig, _ := rand.Int(rand.Reader, big.NewInt(int64(Nnode)))
        temp := int(destBig.Int64())
        if temp != sourceIndex {
            destIndex = temp
            break
        }
    }

    // 回傳來源 (加1000) 和目的 (加2000)
    return sourceIndex + 1000, destIndex + 2000
}

func Random_CANDevices(CAN_Node_Set []int) (int, int) {
	
	 // 1. 隨機選取來源
	 sourceIndexBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(CAN_Node_Set))))
	 sourceIndex := int(sourceIndexBig.Int64())
	 sourceNode := CAN_Node_Set[sourceIndex]
 
	 // 2. 移除來源
	 tempSet := make([]int, 0, len(CAN_Node_Set)-1)
	 tempSet = append(tempSet, CAN_Node_Set[:sourceIndex]...)
	 tempSet = append(tempSet, CAN_Node_Set[sourceIndex+1:]...) 
 
	 // 3. 從新切片中選目的
	 destIndexBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(tempSet))))
	 destIndex := int(destIndexBig.Int64())
	 destinationNode := tempSet[destIndex]

	 return sourceNode-2000, destinationNode-1000
}
