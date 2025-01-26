// pkg/rzd/schemas/station_code.go
package schemas

type StationCodeResponse []StationCode

type StationCode struct {
	N string `json:"n"`
	C string `json:"c"`
}
