package plan

import (
	// "fmt"
	"src/network"
	"src/plan/algo"
)

type Result struct {
    Method string
    Obj    [4]float64
	Cost	int
	Linkmap map[string]float64
}

// type OMACO struct {
// 	Network   *network.Network
// 	SMT       *algo.SMT
// 	MDTC      *algo.MDTC
// 	OSACO     *algo.OSACO
// 	OSACO_IAS *algo.OSACO
// }

// // Developing the OMACO plan
// func new_OMACO_Plan(network *network.Network, osaco_timeout int, osaco_K int, osaco_P float64) *OMACO {
// 	OMACO := &OMACO{Network: network}

// 	OMACO.SMT = &algo.SMT{}
// 	OMACO.MDTC = &algo.MDTC{}
// 	OMACO.OSACO = &algo.OSACO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: 0}
// 	OMACO.OSACO_IAS = &algo.OSACO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: 1}

// 	return OMACO
// }

type OSRO struct {
	Network   	*network.Network
	SP		  	*algo.SP
	KP			*algo.KP
	// SMT       	*algo.SMT
	// MDTC      	*algo.MDTC
	OSRO    	*algo.OSRO
	OSRO_IAS 	*algo.OSRO
	OSRO_method []*algo.OSRO
}

// Developing the OMACO plan
func new_OSRO_Plan(network *network.Network, osaco_timeout int, osaco_K int, osaco_P float64) *OSRO {
	OSRO := &OSRO{Network: network}
	OSRO.SP	 = &algo.SP{}
	OSRO.KP  = &algo.KP{}
	OSRO.OSRO = &algo.OSRO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: 0}
	OSRO.OSRO_IAS = &algo.OSRO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: 1}
	for i:=0; i<4 ;i++{
		OSRO.OSRO_method =append(OSRO.OSRO_method, &algo.OSRO{Timeout: osaco_timeout, K: osaco_K, P: osaco_P, Method_Number: 0})
		// fmt.Println(i)
	}
	

	return OSRO
}
//type plan2 struct {
//	Network *network.Network
//}

// Plan2

//type plan3 struct {
//	Network *network.Network
//}

// Plan3
// ...
