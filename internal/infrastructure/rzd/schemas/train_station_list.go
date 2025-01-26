// pkg/rzd/schemas/train_station_list.go
package schemas

// TrainStationListResponse represents the response from the RZD API
type TrainStationListResponse struct {
	Data TrainStationListData `json:"data"`
}

// TrainStationListData represents the data from the RZD API
type TrainStationListData struct {
	TrainInfo TrainInfo   `json:"trainInfo"`
	Routes    []RouteInfo `json:"routes"`
}

// TrainInfo represents a train info from the RZD API
type TrainInfo struct {
	Number string `json:"number"`
	// Добавьте другие необходимые поля
}

// RouteInfo represents a route info from the RZD API
type RouteInfo struct {
	Station     string `json:"station"`
	ArvTime     string `json:"ArvTime"`
	WaitingTime string `json:"WaitingTime"`
	DepTime     string `json:"DepTime"`
	Distance    int    `json:"Distance"`
}
