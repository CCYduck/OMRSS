package memorizer

import "time"

type OMACO_Memorizer struct {
	average_obj_smt        [4]float64    // {o1, o2, o3, o4}
	average_obj_mdt        [4]float64    // {o1, o2, o3, o4}
	average_objs_osaco     [5][4]float64 // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	average_objs_osaco_ias [5][4]float64 // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	average_time_mdt       time.Duration
	average_time_osaco     [5]time.Duration // 200ms{time} 400ms{time} 600ms{time}, 800ms{time}, 1000ms{time}
	average_time_osaco_ias [5]time.Duration // 200ms{time} 400ms{time} 600ms{time}, 800ms{time}, 1000ms{time}

}

func new_OMACO_Memorizer() *OMACO_Memorizer {
	return &OMACO_Memorizer{}
}

type OSRO_Memorizer struct {
	can2tsn_stats [4][6]float64 // fifo/priority/obo/wat × {Flows,Size,Count,Drop,DelayMs,Testcases}
}

func new_OSRO_Memorizer() *OSRO_Memorizer {
	return &OSRO_Memorizer{}
}

//type Plan3_Computer struct {

//}

//func New_plan3_Memorizer() *Memorizer3 {
//	return
//}
