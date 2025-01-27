// internal/domain/enums.go
package domain

// Direction представляет направление поездки
type Direction int32

const (
	OneWay Direction = iota // Только в одну сторону
	Return                  // Туда и обратно
)

// TrainSearchType представляет тип поезда
type TrainSearchType int32

const (
	AllTrains TrainSearchType = iota + 1 // Поезда и электрички
	Trains                               // Только поезда
	Electrics                            // Только электрички
)

type SeatType int32

const (
	Platz   SeatType = iota + 1 // Плац
	General                     // Общий
	Side                        // Cид
	Coupe                       // Купе
	Soft                        // Мягкий
	Lux                         // Люкс
)

type TrainType int32

const (
	Train    TrainType = iota // Поезд
	Suburban                  // Электричка TODO удостовериться
)
