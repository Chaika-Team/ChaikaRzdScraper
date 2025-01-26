// internal/domain/enums.go
package domain

// Direction представляет направление поездки
type Direction int32

const (
	OneWay Direction = iota // Только в одну сторону
	Return                  // Туда и обратно
)

// TrainType представляет тип поезда
type TrainType int32

const (
	AllTrains TrainType = iota + 1 // Поезда и электрички
	Trains                         // Только поезда
	Electrics                      // Только электрички
)
