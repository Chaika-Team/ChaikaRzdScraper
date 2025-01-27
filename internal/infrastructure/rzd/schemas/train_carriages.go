// pkg/rzd/schemas/train_carriages.go
package schemas

// TrainCarriagesResponse represents the response from the RZD API
type TrainCarriagesResponse struct {
	Lst              []TrainCarriage `json:"list"`
	Schemes          []string        `json:"schemes"`
	InsuranceCompany []string        `json:"insuranceCompany"`
}

// TrainCarriage represents a train carriage from the RZD API
type TrainCarriage struct {
	Cars           []CarriageType `json:"cars"`
	FunctionBlocks []string       `json:"functionBlocks"`
}
