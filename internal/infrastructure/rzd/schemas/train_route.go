// pkg/rzd/schemas/train_route.go
package schemas

type TrainRouteResponse struct {
	Result string `json:"result"`
	TP     []TP   `json:"tp"`
}

type TP struct {
	From        string      `json:"from"`
	FromCode    int         `json:"fromCode"`
	Where       string      `json:"where"`
	WhereCode   int         `json:"whereCode"`
	Date        string      `json:"date"`
	NoSeats     bool        `json:"noSeats"`
	DefShowTime string      `json:"defShowTime"`
	State       string      `json:"state"`
	List        []TrainList `json:"list"`
	Cur         []int       `json:"cur"`
}

type TrainList struct {
	MidCode            int           `json:"midCode"`
	MidPt              string        `json:"midPt"`
	TransTime          int           `json:"transTime"`
	TimeBetweenCurCase int           `json:"timeBetweenCurCase"`
	TotalTravelTime    int           `json:"totalTravelTime"`
	Cases              [][]TrainCase `json:"cases"`
	Cur                []int         `json:"cur"`
}

type TrainCase struct {
	Number            string     `json:"number"`
	Number2           string     `json:"number2"`
	Type              int        `json:"type"`
	TypeEx            int        `json:"typeEx"`
	Depth             int        `json:"depth"`
	New               bool       `json:"new"`
	ElReg             bool       `json:"elReg"`
	DeferredPayment   bool       `json:"deferredPayment"`
	VarPrice          bool       `json:"varPrice"`
	Code0             int        `json:"code0"`
	Code1             int        `json:"code1"`
	BEntire           bool       `json:"bEntire"`
	TrainName         string     `json:"trainName"`
	BFirm             bool       `json:"bFirm"`
	Brand             string     `json:"brand"`
	Carrier           string     `json:"carrier"`
	Route0            string     `json:"route0"`
	Route1            string     `json:"route1"`
	TrDate0           string     `json:"trDate0"`
	TrTime0           string     `json:"trTime0"`
	Station0          string     `json:"station0"`
	Station1          string     `json:"station1"`
	Date0             string     `json:"date0"`
	Time0             string     `json:"time0"`
	Date1             string     `json:"date1"`
	Time1             string     `json:"time1"`
	TimeInWay         string     `json:"timeInWay"`
	FlMsk             int        `json:"flMsk"`
	TrainID           int        `json:"train_id"`
	Cars              []Carriage `json:"cars"`
	CarNumeration     string     `json:"carNumeration"`
	DisabledType      bool       `json:"disabledType"`
	AddCompLuggageNum int        `json:"addCompLuggageNum"`
	AddCompLuggage    bool       `json:"addCompLuggage"`
	AddHandLuggage    bool       `json:"addHandLuggage"`
}

type Carriage struct {
	CarDataType    int    `json:"carDataType"`
	Itype          int    `json:"itype"`
	Type           string `json:"type"`
	TypeLoc        string `json:"typeLoc"`
	FreeSeats      int    `json:"freeSeats"`
	Pt             int    `json:"pt"`
	Tariff         int    `json:"tariff"`
	ServCls        string `json:"servCls"`
	DisabledPerson bool   `json:"disabledPerson,omitempty"`
}
