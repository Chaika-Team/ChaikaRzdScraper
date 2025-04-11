// pkg/rzd/schemas/station_code.go
package schemas

// StationCodeResponse represents the response from the RZD API
type StationCodeResponse []StationCode

// StationCode represents a station code from the RZD API
type StationCode struct {
	N string `json:"n"` // Название станции, как в API
	C int    `json:"c"` // Код станции, как в API
	L int    `json:"L"` // Уровень станции (0-5, 5 - самый высокий)
	S int    `json:"S"` // Значение сортировки (0-5, 5 - самое высокое)
}
