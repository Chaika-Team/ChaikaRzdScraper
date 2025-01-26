// pkg/rzd/schemas/train_carriages.go
package schemas

type TrainCarriagesResponse struct {
	Lst              []TrainCarriage `json:"lst"`
	Schemes          []string        `json:"schemes"`
	InsuranceCompany []string        `json:"insuranceCompany"`
}

type TrainCarriage struct {
	Cars           []Carriage `json:"cars"`
	FunctionBlocks []string   `json:"functionBlocks"`
}
