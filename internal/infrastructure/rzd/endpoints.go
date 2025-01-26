package rzd

import (
	"fmt"
	"net/url"
	"strings"
)

// Constants for layer IDs
const (
	RoutesLayer         = 5827
	CarriagesLayer      = 5764
	StationsStructureID = 704
)

// Endpoints содержит все пути эндпоинтов относительно базового пути
type Endpoints struct {
	TrainRoutes       string
	TrainRoutesReturn string
	TrainCarriages    string
	TrainStationList  string
	StationCode       string
}

// NewEndpoints создает все эндпоинты на основе базового пути и языка, используя url.URL
func NewEndpoints(basePath, language string) Endpoints {
	baseURL, err := url.Parse(basePath)
	if err != nil {
		panic(fmt.Sprintf("invalid base path: %v", err))
	}

	return Endpoints{
		TrainRoutes:       buildEndpointURL(baseURL, fmt.Sprintf("timetable/public/%s", language), RoutesLayer),
		TrainRoutesReturn: buildEndpointURL(baseURL, fmt.Sprintf("timetable/public/%s", language), RoutesLayer),
		TrainCarriages:    buildEndpointURL(baseURL, fmt.Sprintf("timetable/public/%s", language), CarriagesLayer),
		TrainStationList:  buildStaticEndpoint(baseURL, "ticket/services/route/basicRoute"),
		StationCode:       buildStaticEndpoint(baseURL, "suggester"),
	}
}

// buildEndpointURL конструирует URL эндпоинта с параметром layer_id
func buildEndpointURL(baseURL *url.URL, path string, layerID int) string {
	u := *baseURL // Копируем baseURL, чтобы избежать изменения оригинального
	u.Path = joinPaths(u.Path, path)
	query := u.Query()
	query.Set("layer_id", fmt.Sprintf("%d", layerID))
	u.RawQuery = query.Encode()
	return u.String()
}

// buildStaticEndpoint конструирует статический URL эндпоинта без дополнительных параметров
func buildStaticEndpoint(baseURL *url.URL, path string) string {
	u := *baseURL // Копируем baseURL, чтобы избежать изменения оригинального
	u.Path = joinPaths(u.Path, path)
	return u.String()
}

// joinPaths корректно объединяет пути
func joinPaths(basePath, addPath string) string {
	trimmedBase := strings.TrimSuffix(basePath, "/")
	trimmedAdd := strings.TrimPrefix(addPath, "/")
	return fmt.Sprintf("%s/%s", trimmedBase, trimmedAdd)
}
