package mappers

import (
	"strings"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd/schemas"
)

// MapStationCodeResponse преобразует ответ API (StationCodeResponse) в срез доменных моделей []domain.Station.
// Фильтрация производится по тому, что название станции (поле "n") содержит поисковый запрос (без учёта регистра).
// Остальные поля ответа (например, S, L) игнорируются.
func MapStationCodeResponse(resp schemas.StationCodeResponse, query string) ([]domain.Station, error) {
	var stations []domain.Station
	lowerQuery := strings.ToLower(query)

	for _, s := range resp {
		// Если имя станции содержит искомую подстроку (без учёта регистра)
		if strings.Contains(strings.ToLower(s.N), lowerQuery) {
			// Преобразуем код в строку, если доменная модель ожидает строку, либо сохраняем как число
			// Здесь в доменной модели Station.Code определён как int.
			stations = append(stations, domain.Station{
				Name: s.N,
				Code: s.C,
			})
		}
	}

	return stations, nil
}
