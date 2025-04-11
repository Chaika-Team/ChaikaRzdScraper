// pkg/rzd/schemas/train_carriages.go
package schemas

// TrainCarriagesResponse представляет полный ответ от API РЖД на запрос информации о вагонах.
type TrainCarriagesResponse struct {
	Result                string                 `json:"result"`                // Результат запроса, обычно "OK"
	Lst                   []TrainResult          `json:"lst"`                   // Список объектов, описывающих конкретные вагоны (поезда)
	Schemes               []Schemes              `json:"schemes"`               // Список схем вагонов
	InsuranceCompany      []InsuranceCompany     `json:"insuranceCompany"`      // Список страховых компаний
	InsuranceCompanyTypes []InsuranceCompanyType `json:"insuranceCompanyTypes"` // Типы страховых тарифов и программ
	Psaction              interface{}            `json:"psaction"`              // Дополнительное действие (может быть null)
	ChildrenAge           int                    `json:"childrenAge"`           // Возраст детей для тарифных расчетов
	MotherAndChildAge     int                    `json:"motherAndChildAge"`     // Возраст для тарифов "мать и ребёнок"
	PartialPayment        bool                   `json:"partialPayment"`        // Флаг частичной оплаты
	Timestamp             string                 `json:"timestamp"`             // Временная метка ответа
}

// TrainResult представляет данные по конкретному поезду, полученные от РЖД.
type TrainResult struct {
	Result         string `json:"result"`         // Результат для данного поезда, обычно "OK"
	TrainNumber    string `json:"number"`         // Номер поезда (например, "119А")
	TrainNumber2   string `json:"number2"`        // Дублирующий номер поезда (если имеется)
	DefShowTime    string `json:"defShowTime"`    // Тип отображения времени (например, "local")
	Date0          string `json:"date0"`          // Дата отправления (формат DD.MM.YYYY)
	Time0          string `json:"time0"`          // Время отправления (формат HH:MM)
	Date1          string `json:"date1"`          // Дата прибытия (формат DD.MM.YYYY)
	Time1          string `json:"time1"`          // Время прибытия (формат HH:MM)
	Type           string `json:"type"`           // Тип вагона (например, "СК")
	Virtual        bool   `json:"virtual"`        // Флаг виртуального поезда?
	Bus            bool   `json:"bus"`            // Флаг автобусного соединения
	Boat           bool   `json:"boat"`           // Флаг водного соединения
	Station0       string `json:"station0"`       // Название станции отправления поезда
	Code0          string `json:"code0"`          // Код станции отправления (в виде строки)
	Station1       string `json:"station1"`       // Название станции прибытия
	Code1          string `json:"code1"`          // Код станции прибытия (в виде строки)
	TimeSt0        string `json:"timeSt0"`        // Дополнительное время отправления (например, время стоянки)
	TimeSt1        string `json:"timeSt1"`        // Дополнительное время прибытия
	Route0         string `json:"route0"`         // Краткое название маршрута отправления
	Route1         string `json:"route1"`         // Краткое название маршрута прибытия
	Cars           []Car  `json:"cars"`           // Список конкретных вагонов (детали состава)
	AddCompLuggage bool   `json:"addCompLuggage"` // Флаг дополнительного багажа
	Timestamp      string `json:"timestamp"`      // Временная метка ответа формата "12.02.2025 16:03:33.515"
	// FunctionBlocks

}

// Car представляет подробную информацию о конкретном вагоне.
type Car struct {
	Cnumber            string      `json:"cnumber"`            // Номер вагона (внутренний номер, например "01")
	Type               string      `json:"type"`               // Тип вагона (например, "Купе", "Плац", "Люкс")
	CatLabelLoc        string      `json:"catLabelLoc"`        // Локальная метка категории вагона (например, "Купе")
	TypeLoc            string      `json:"typeLoc"`            // Локальное наименование типа вагона (например, "Купе")
	CatCode            string      `json:"catCode"`            // Код категории вагона
	Ctypei             int         `json:"ctypei"`             // Идентификатор типа вагона (числовой, например, 4)
	Ctype              int         `json:"ctype"`              // Повторный идентификатор типа вагона
	Letter             string      `json:"letter"`             // Буква вагона (например, "А")
	ClsType            string      `json:"clsType"`            // Тип класса вагона (например, "2Ш")
	SubType            string      `json:"subType"`            // Подтип вагона (например, "66К")
	ClsName            string      `json:"clsName"`            // Полное описание вагона с HTML-разметкой
	Services           []Service   `json:"services"`           // Список услуг, предоставляемых в вагоне
	Tariff             string      `json:"tariff"`             // Тариф за билет (в виде строки, например, "2533")
	Tariff2            string      `json:"tariff2"`            // Дополнительный тариф (например, "3463", может быть null)
	TariffServ         *string     `json:"tariffServ"`         // Тариф за услугу (если указан, может быть null)
	AddSigns           string      `json:"addSigns"`           // Дополнительные знаки/отметки (например, пустая строка)
	Carrier            string      `json:"carrier"`            // Перевозчик вагона (например, "ФПК")
	CarrierId          int         `json:"carrierId"`          // Идентификатор перевозчика
	InsuranceFlag      bool        `json:"insuranceFlag"`      // Флаг наличия страхования
	InsuranceTypeId    int         `json:"insuranceTypeId"`    // Тип страхования (числовое значение)
	Owner              string      `json:"owner"`              // Владелец вагона (например, "РЖД/МСК")
	ElReg              bool        `json:"elReg"`              // Флаг электронной регистрации
	Food               bool        `json:"food"`               // Флаг наличия питания
	SelFood            bool        `json:"selFood"`            // Флаг возможности выбора питания
	EquippedSIOP       bool        `json:"equippedSIOP"`       // Флаг оснащенности СИОП (информационно-обслуживающей системы)
	RegularFoodService bool        `json:"regularFoodService"` // Флаг регулярного обслуживания питанием
	NoSmok             bool        `json:"noSmok"`             // Флаг запрета курения
	InetSaleOff        bool        `json:"inetSaleOff"`        // Флаг отсутствия интернет-продаж
	BVip               bool        `json:"bVip"`               // Флаг VIP-условий
	ConferenceRoomFlag bool        `json:"conferenceRoomFlag"` // Флаг наличия переговорной комнаты
	BDeck2             bool        `json:"bDeck2"`             // Флаг второго палубного вагона (если применимо)
	IntServiceClass    interface{} `json:"intServiceClass"`    // Дополнительная информация о классе обслуживания (может быть null)
	SpecialSeatTypes   interface{} `json:"specialSeatTypes"`   // Особые типы мест (может быть null)
	DeferredPayment    bool        `json:"deferredPayment"`    // Флаг отложенной оплаты
	VarPrice           bool        `json:"varPrice"`           // Флаг вариативного тарифа
	Ferry              bool        `json:"ferry"`              // Флаг паромного сообщения
	SeniorTariff       int         `json:"seniorTariff"`       // Тариф для пожилых пассажиров
	Bedding            bool        `json:"bedding"`            // Флаг предоставления постельного белья
	NonRefundable      bool        `json:"nonRefundable"`      // Флаг безвозвратности билета
	AddTour            bool        `json:"addTour"`            // Флаг наличия тура/экскурсии
	CarNumeration      *string     `json:"carNumeration"`      // Информация о нумерации вагона (например, "FromHead"; может быть null)
	AddHandLuggage     bool        `json:"addHandLuggage"`     // Флаг добавления услуги провоза ручной клади
	Youth              bool        `json:"youth"`              // Флаг тарифов для молодежи
	Unior              bool        `json:"unior"`              // Флаг тарифов для студентов или иной категории
	Seats              []Seat      `json:"seats"`              // Список мест в данном вагоне
	Places             string      `json:"places"`             // Строка с перечнем мест (номера мест)
	SchemeID           int         `json:"schemeId"`           // Идентификатор схемы вагона
	SchemeInfo         SchemeInfo  `json:"schemeInfo"`         // Информация о схеме вагона (пути к изображениям, легенда)
	ForcedBedding      bool        `json:"forcedBedding"`      // Флаг принудительного выкупа спальных мест
	PolicyEnabled      bool        `json:"policyEnabled"`      // Флаг включения политики
	Msr                bool        `json:"msr"`                // Флаг MSR (специфическая информация РЖД)
	Medic              bool        `json:"medic"`              // Флаг наличия медицинского обслуживания
}

// Seat представляет информацию о конкретном месте в вагоне.
type Seat struct {
	Type         string  `json:"type"`         // Тип места: "dn" (нижнее) или "up" (верхнее)
	Code         string  `json:"code"`         // Код места (например, "Н" для нижнего, "В" для верхнего)
	Label        string  `json:"label"`        // Наименование места (например, "Нижнее")
	Tariff       string  `json:"tariff"`       // Тариф за место (в виде строки)
	Tariff2      string  `json:"tariff2"`      // Дополнительный тариф за место (может быть null)
	TariffServ   string  `json:"tariffServ"`   // Тариф за услугу для места (если указан)
	Free         int     `json:"free"`         // Количество свободных мест данного типа
	PlacesNonRef *string `json:"placesNonRef"` // Места, недоступные для возврата (может быть null)
	FreeRef      int     `json:"freeRef"`      // Количество мест, доступных для возврата
	PlacesRef    string  `json:"placesRef"`    // Номера мест, доступных для возврата (строка)
	Places       string  `json:"places"`       // Полный перечень номеров мест (строка)
}

// Service представляет услугу, предоставляемую в вагоне.
type Service struct {
	ID          int    `json:"id"`          // Идентификатор услуги
	Name        string `json:"name"`        // Название услуги (с иконкой)
	Description string `json:"description"` // Описание услуги
	HasImage    bool   `json:"hasImage"`    // Флаг: у услуги имеется изображение
}

// SchemeInfo представляет информацию о схеме вагона.
type SchemeInfo struct {
	Dir     string `json:"dir"`     // Путь к изображению схемы вагона
	DirVert string `json:"dirVert"` // Путь к вертикальной схеме вагона
	Legend  string `json:"legend"`  // Легенда (описание значков) для схемы вагона
}

// InsuranceCompany представляет данные о страховой компании.
type InsuranceCompany struct {
	ID               int    `json:"id"`               // Идентификатор страховой компании
	ShortName        string `json:"shortName"`        // Краткое название страховой компании
	OfferUrl         string `json:"offerUrl"`         // URL предложения страховой компании
	InsuranceCost    int    `json:"insuranceCost"`    // Стоимость страхования
	InsuranceBenefit int    `json:"insuranceBenefit"` // Страховая выплата или лимит
	SortOrder        int    `json:"sortOrder"`        // Порядок сортировки
}

// InsuranceCompanyType представляет тип страховой компании с тарифами.
type InsuranceCompanyType struct {
	TypeId           int               `json:"typeId"`           // Идентификатор типа страхования
	InsuranceTariffs []InsuranceTariff `json:"insuranceTariffs"` // Список тарифов для данного типа
}

// InsuranceTariff представляет тариф страхования.
type InsuranceTariff struct {
	ID                int                `json:"id"`                // Идентификатор тарифа
	Name              string             `json:"name"`              // Название тарифа (например, "Базовый")
	InsuranceCost     int                `json:"insuranceCost"`     // Стоимость страхования по тарифу
	InsuranceBenefit  int                `json:"insuranceBenefit"`  // Страховая выплата или лимит по тарифу
	Default           bool               `json:"default"`           // Флаг тарифа по умолчанию
	InsurancePrograms []InsuranceProgram `json:"InsurancePrograms"` // Список страховых программ для тарифа
}

// InsuranceProgram представляет страховую программу.
type InsuranceProgram struct {
	ID        int    `json:"id"`        // Идентификатор программы
	OfferUrl  string `json:"offerUrl"`  // URL предложения программы
	SortOrder int    `json:"sortOrder"` // Порядок сортировки
	ShortName string `json:"shortName"` // Краткое название программы
}

type Schemes struct {
	ID    int    `json:"id"`    // Идентификатор схемы
	HTML  string `json:"html"`  // HTML-разметка схемы
	Image string `json:"image"` // Путь к изображению схемы (если имеется)
}
