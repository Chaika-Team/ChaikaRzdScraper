// pkg/rzd/schemas/train_station_list.go
package schemas

type TrainStationListResponse struct {
	Data TrainStationListData `json:"data"`
}

type TrainStationListData struct {
	TrainInfo TrainInfo   `json:"trainInfo"`
	Routes    []RouteInfo `json:"routes"`
}

type TrainInfo struct {
	Number string `json:"number"`
	// Добавьте другие необходимые поля
}

type RouteInfo struct {
	Station     string `json:"station"`
	ArvTime     string `json:"ArvTime"`
	WaitingTime string `json:"WaitingTime"`
	DepTime     string `json:"DepTime"`
	Distance    int    `json:"Distance"`
}
