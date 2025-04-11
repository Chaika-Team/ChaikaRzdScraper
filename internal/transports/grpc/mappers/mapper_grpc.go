// internal/transports/grpc/mappers/mappers.go
package mappers

import (
	"time"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/transports/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MapTrainRoutesToPb преобразует срез доменных TrainRoute в pb.GetTrainRoutesResponse.
// В данной реализации используются google.protobuf.Timestamp для полей времени.
func MapTrainRoutesToPb(routes []domain.TrainRoute) *pb.GetTrainRoutesResponse {
	var pbRoutes []*pb.TrainRoute
	for _, r := range routes {
		pbRoute := &pb.TrainRoute{
			TrainNumber: r.TrainNumber,
			TrainType:   int32(r.TrainType),
			Departure:   timestamppb.New(r.Departure),
			Arrival:     timestamppb.New(r.Arrival),
			From:        MapStationToPb(r.From),
			To:          MapStationToPb(r.To),
		}
		// Маппим агрегированные типы вагонов
		for _, ct := range r.CarTypes {
			pbCT := &pb.CarriageType{
				Type:           int32(ct.Type),
				TypeShortLabel: ct.TypeShortLabel,
				TypeLabel:      ct.TypeLabel,
				Class:          ct.Class,
				Tariff:         int32(ct.Tariff),
				TariffExtra:    int32(ct.TariffExtra),
				FreeSeats:      int32(ct.FreeSeats),
				Disabled:       ct.Disabled,
			}
			pbRoute.CarTypes = append(pbRoute.CarTypes, pbCT)
		}
		pbRoutes = append(pbRoutes, pbRoute)
	}
	return &pb.GetTrainRoutesResponse{
		Routes: pbRoutes,
	}
}

// MapTrainCarriagesToPb преобразует срез доменных Car в pb.GetTrainCarriagesResponse.
func MapTrainCarriagesToPb(cars []domain.Car) *pb.GetTrainCarriagesResponse {
	var pbCars []*pb.Car
	for _, c := range cars {
		pbCar := &pb.Car{
			CarNumber:          c.CarNumber,
			Type:               c.Type,
			CategoryLabelLocal: c.CategoryLabelLocal,
			TypeLabel:          c.TypeLabel,
			CategoryCode:       c.CategoryCode,
			CarTypeId:          int32(c.CarTypeID),
			CarType:            int32(c.CarType),
			Letter:             c.Letter,
			ClassType:          c.ClassType,
			Tariff:             int32(c.Tariff),
			TariffExtra:        int32(c.Tariff2),
			Carrier:            MapCarrierToPb(c.Carrier),
			CarNumeration:      int32(MapCarNumerationToPb(c.CarNumeration)),
		}
		// Маппим список услуг
		for _, s := range c.Services {
			pbCar.Services = append(pbCar.Services, &pb.Service{
				Id:          s.ID,
				Name:        s.Name,
				Description: s.Description,
			})
		}
		// Если потребуется маппить места, их можно добавить сюда.
		pbCars = append(pbCars, pbCar)
	}
	return &pb.GetTrainCarriagesResponse{
		Carriages: pbCars,
	}
}

// MapStationsToPb преобразует срез доменных Station в pb.SearchStationResponse.
func MapStationsToPb(stations []domain.Station) *pb.SearchStationResponse {
	var pbStations []*pb.Station
	for _, s := range stations {
		pbStations = append(pbStations, &pb.Station{
			Name:      s.Name,
			Code:      int32(s.Code),
			RouteName: s.RouteName,
			Level:     int32(s.Level),
			Score:     int32(s.Score),
		})
	}
	return &pb.SearchStationResponse{
		Stations: pbStations,
	}
}

// MapStationToPb преобразует доменную Station в pb.Station.
func MapStationToPb(s domain.Station) *pb.Station {
	return &pb.Station{
		Name:      s.Name,
		Code:      int32(s.Code),
		RouteName: s.RouteName,
		Level:     int32(s.Level),
		Score:     int32(s.Score),
	}
}

// MapCarrierToPb преобразует доменного Carrier в pb.Carrier.
func MapCarrierToPb(c domain.Carrier) *pb.Carrier {
	return &pb.Carrier{
		Id:   c.ID,
		Name: c.Name,
	}
}

// MapCarNumerationToPb преобразует доменное CarNumeration в int32 для pb.
// Используем определённые константы: Head=0, Tail=1, Unknown=2.
func MapCarNumerationToPb(cn domain.CarNumeration) int32 {
	// Если типы совпадают, можно вернуть значение напрямую.
	// В данном случае возвращаем cn, и в endpoint-е результат приведём к int32.
	switch cn {
	case domain.Head:
		return 0 // Head
	case domain.Tail:
		return 1 // Tail
	default:
		return 2 // Unknown
	}
}

// ParseTimestampToTime преобразует protobuf Timestamp в time.Time.
func ParseTimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// ParseTimeRequest – пример функции для преобразования поля времени из запроса.
// Если в запросе используется google.protobuf.Timestamp, достаточно вызвать ts.AsTime()
func ParseTimeRequest(ts *timestamppb.Timestamp) time.Time {
	return ParseTimestampToTime(ts)
}

// ParseDateRequest аналогично преобразует Timestamp в дату (time.Time).
func ParseDateRequest(ts *timestamppb.Timestamp) time.Time {
	return ParseTimestampToTime(ts)
}
