package flow

import "time"

type TSN struct {
	Period   int     // 100~2000us up 500us
	Deadline int     // Period = Deadline
	DataSize float64 // 30~100bytes up 10bytes 最多可以封裝12個CAN封包
}

func new_TSN(t_period int, t_datasize float64) *TSN {
	return &TSN{
		Period:   t_period,
		Deadline: t_period,
		DataSize: t_datasize,
	}
}

type AVB struct {
	Period   int     // 125us
	Deadline int     // 2000us
	DataSize float64 // 1000~1500bytes  up 100bytes
}

func new_AVB(a_datasize float64) *AVB {
	return &AVB{
		Period:   125,
		Deadline: 2000,
		DataSize: a_datasize,
	}
}

type importantCAN struct {
	Period   int
	Deadline int
	DataSize float64
}

func new_importantCAN() *importantCAN {
	return &importantCAN{
		Period:   5000, // 5000us
		Deadline: 5000, // Period = Deadline
		DataSize: 8,    // 8bytes
	}
}

type unimportantCAN struct {
	Period   int
	Deadline int
	DataSize float64
}

func new_unimportantCAN(uc_period int, uc_deadline int) *unimportantCAN {
	return &unimportantCAN{
		Period:   uc_period,   // 50000~150000us up 50000us
		Deadline: uc_deadline, // 10000~20000us up 2000us
		DataSize: 8,           // 8bytes
	}
}

type Stream struct {
	Name        string
	ArrivalTime int
	SendTime	int
	DataSize    float64
	Deadline    int
	FinishTime  int
}

func new_TTStream(name string, arrivalTime int, datasize float64, deadline int, finishTime int) *Stream {
	return &Stream{
		Name:        name,
		ArrivalTime: arrivalTime,
		DataSize:    datasize,
		Deadline:    deadline,
		FinishTime:  finishTime,
	}
}

func new_CANStream(name string, arrivalTime int, datasize float64, deadline int, finishTime int) *Stream {
	return &Stream{
		Name:        name,
		ArrivalTime: arrivalTime,
		DataSize:    datasize,
		Deadline:    deadline,
		FinishTime:  finishTime,
	}
}

type Flow struct {
	Period      int
	Deadline    int
	DataSize    float64
	HyperPeriod int
	Source      int
	Destination int
	Streams     []*Stream
}

func new_TTFlow(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	return &Flow{
		Period:      period,
		Deadline:    deadline,
		DataSize:    datasize,
		HyperPeriod: HyperPeriod,
	}
}

func new_CANFlow(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	return &Flow{
		Period:      period,
		Deadline:    deadline,
		DataSize:    datasize,
		HyperPeriod: HyperPeriod,
	}
}

type Flows struct {
	TSNFlows        []*Flow
	AVBFlows        []*Flow
	Encapsulate    []*Method
}

func new_Flows() *Flows {
	return &Flows{}
}

type Method struct{
	Method_Name		string		
	CAN2TSNFlows    []*Flow
	BytesSent	float64
	TSNFrameCount    int 
	CAN2TSN_O1_Drop int
	CAN2TSN_Delay   time.Duration
}