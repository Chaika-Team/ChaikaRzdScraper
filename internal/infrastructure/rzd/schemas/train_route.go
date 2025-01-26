// pkg/rzd/schemas/train_route.go
package schemas

// TrainRouteResponse represents the response from the RZD API
type TrainRouteResponse struct {
	Result string `json:"result"`
	TP     []TP   `json:"tp"`
}

// TP represents a train route from the RZD API
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

// TrainList represents a train list from the RZD API
type TrainList struct {
	Number            string         `json:"number"`
	Number2           string         `json:"number2"`
	Type              int            `json:"type"`
	TypeEx            int            `json:"typeEx"`
	Depth             int            `json:"depth"`
	New               bool           `json:"new"`
	ElReg             bool           `json:"elReg"`
	DeferredPayment   bool           `json:"deferredPayment"`
	VarPrice          bool           `json:"varPrice"`
	Code0             int            `json:"code0"`
	Code1             int            `json:"code1"`
	BEntire           bool           `json:"bEntire"`
	TrainName         string         `json:"trainName"`
	BFirm             bool           `json:"bFirm"`
	Brand             string         `json:"brand"`
	Carrier           string         `json:"carrier"`
	Route0            string         `json:"route0"`
	Route1            string         `json:"route1"`
	TrDate0           string         `json:"trDate0"`
	TrTime0           string         `json:"trTime0"`
	Station0          string         `json:"station0"`
	Station1          string         `json:"station1"`
	Date0             string         `json:"date0"`
	Time0             string         `json:"time0"`
	Date1             string         `json:"date1"`
	Time1             string         `json:"time1"`
	TimeInWay         string         `json:"timeInWay"`
	FlMsk             int            `json:"flMsk"`
	TrainID           int            `json:"train_id"`
	Cars              []Carriage     `json:"cars"`
	SeatCars          []SeatCarriage `json:"seatCars,omitempty"`
	CarNumeration     string         `json:"carNumeration"`
	DisabledType      bool           `json:"disabledType"`
	AddCompLuggageNum int            `json:"addCompLuggageNum"`
	AddCompLuggage    bool           `json:"addCompLuggage"`
	AddHandLuggage    bool           `json:"addHandLuggage"`
}

// Carriage represents a carriage from the RZD API
type Carriage struct {
	CarDataType    int    `json:"carDataType"`
	Itype          int    `json:"itype"`
	Type           string `json:"type"`
	TypeLoc        string `json:"typeLoc"`
	FreeSeats      int    `json:"freeSeats"`
	Pt             int    `json:"pt"`
	Tariff         int    `json:"tariff"` // Да-да, это int, а не string. А в seatcarriage - string
	ServCls        string `json:"servCls"`
	DisabledPerson bool   `json:"disabledPerson,omitempty"`
}

// SeatCarriage represents a seat carriage from the RZD API, whatever that means
type SeatCarriage struct { // Да, они разные
	CarDataType int    `json:"carDataType"`
	Itype       int    `json:"itype"`
	Type        string `json:"type"`
	TypeLoc     string `json:"typeLoc"`
	FreeSeats   int    `json:"freeSeats"`
	Pt          int    `json:"pt"`
	Tariff      string `json:"tariff"`            // Да-да, это string, а не int. А в carriage - int
	Tariff2     string `json:"tariff2,omitempty"` // И тут тоже string
	ServCls     string `json:"servCls"`
	// И тут нет disabledPerson, да
}
