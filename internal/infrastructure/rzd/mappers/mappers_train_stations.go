package mappers

import (
	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure/rzd/schemas"
)

// MapStationCodeResponse преобразует ответ API (StationCodeResponse) в срез доменных моделей []domain.Station.
// Остальные поля ответа (например, S, L) игнорируются.
func MapStationCodeResponse(resp schemas.StationCodeResponse) ([]domain.Station, error) {
	var stations []domain.Station

	for _, s := range resp {
		stations = append(stations, domain.Station{
			Name:  s.N,
			Code:  s.C,
			Level: s.L,
			Score: s.S,
		})
	}

	return stations, nil
}
