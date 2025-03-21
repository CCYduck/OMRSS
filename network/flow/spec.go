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
	unimportantCAN_period_arr := []int{50000, 60000, 70000, 80000, 90000, 100000}
	unimportantCAN_deadline := []int{10000, 12000, 14000, 16000, 18000, 20000}
	period_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(unimportantCAN_period_arr))))
	deadline_rng, _ := rand.Int(rand.Reader, big.NewInt(int64(len(unimportantCAN_deadline))))

	return unimportantCAN_period_arr[period_rng.Int64()], unimportantCAN_deadline[deadline_rng.Int64()]
}


func Random_Devices(Nnode int) (int, []int) {
	// Talker
	source, _ := rand.Int(rand.Reader, big.NewInt(int64(Nnode)))
	// Listener
	destinations := []int{}
	for i := 0; i < Nnode; i++ {
    	if i != int(source.Int64()) {
        destinations = append(destinations, i+2000)
    	}
	}
	randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(destinations))))
	selectedDestination := destinations[int(randIndex.Int64())]
	
	return int(source.Int64()) + 1000, []int{selectedDestination}
}

func Random_CANDevices(Cnode int) (int, []int) {
	// Talker
	source, _ := rand.Int(rand.Reader, big.NewInt(int64(Cnode)))
	// Listener
	destinations := []int{}
	for i := 0; i < Cnode; i++ {
    	if i != int(source.Int64()) {
        destinations = append(destinations, i+4000)
    	}
	}
	randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(destinations))))
	selectedDestination := destinations[int(randIndex.Int64())]
	
	return int(source.Int64()) + 1000, []int{selectedDestination}
}
