// pkg/rzd/schemas/train_station_list.go
package schemas

// TrainStationListResponse // TODO не готово
type TrainStationListResponse struct {
	Data TrainStationListData `json:"data"`
}

// TrainStationListData // TODO не готово
type TrainStationListData struct {
	TrainInfo TrainInfo   `json:"trainInfo"`
	Routes    []RouteInfo `json:"routes"`
}

// TrainInfo // TODO не готово
type TrainInfo struct {
	Number string `json:"number"`
	// Добавьте другие необходимые поля
}

// RouteInfo // TODO не готово
type RouteInfo struct {
	Station     string `json:"station"`
	ArvTime     string `json:"ArvTime"`
	WaitingTime string `json:"WaitingTime"`
	DepTime     string `json:"DepTime"`
	Distance    int    `json:"Distance"`
}
