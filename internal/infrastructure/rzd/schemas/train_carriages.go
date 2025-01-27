// pkg/rzd/schemas/train_carriages.go
package schemas

// TODO не готово
type TrainCarriagesResponse struct {
	Lst              []TrainCarriage `json:"list"`
	Schemes          []string        `json:"schemes"`
	InsuranceCompany []string        `json:"insuranceCompany"`
}

// TODO не готово
type TrainCarriage struct {
	Cars           []CarriageType `json:"cars"`
	FunctionBlocks []string       `json:"functionBlocks"`
}
