// internal/adapters/http/handlers.go

package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/usecases"

	"github.com/gorilla/mux"
)

type Handlers struct {
	service *usecases.RzdService
}

func NewHandlers(service *usecases.RzdService) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/train_routes", h.GetTrainRoutes).Methods("GET")
	router.HandleFunc("/train_routes_return", h.GetTrainRoutesReturn).Methods("GET")
	router.HandleFunc("/train_carriages", h.GetTrainCarriages).Methods("GET")
	router.HandleFunc("/train_station_list", h.GetTrainStationList).Methods("GET")
	router.HandleFunc("/station_code", h.GetStationCode).Methods("GET")
}

func (h *Handlers) GetTrainRoutes(w http.ResponseWriter, r *http.Request) {
	var params domain.TrainRoutesParams
	// Парсинг параметров из запроса
	params.Dir = parseInt(r, "dir", 0)
	params.Tfl = parseInt(r, "tfl", 3)
	params.CheckSeats = parseInt(r, "checkSeats", 1)
	params.Code0 = r.URL.Query().Get("code0")
	params.Code1 = r.URL.Query().Get("code1")
	params.Md = parseInt(r, "md", 0)

	dt0Str := r.URL.Query().Get("dt0")
	if dt0Str == "" {
		params.Dt0 = time.Now().AddDate(0, 0, 1)
	} else {
		params.Dt0, _ = time.Parse("02.01.2006", dt0Str)
	}

	trips, err := h.service.GetTrainRoutes(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, trips)
}

func (h *Handlers) GetTrainRoutesReturn(w http.ResponseWriter, r *http.Request) {
	var params domain.TrainRoutesReturnParams
	// Парсинг параметров из запроса
	params.Dir = parseInt(r, "dir", 1)
	params.Tfl = parseInt(r, "tfl", 3)
	params.CheckSeats = parseInt(r, "checkSeats", 1)
	params.Code0 = r.URL.Query().Get("code0")
	params.Code1 = r.URL.Query().Get("code1")

	dt0Str := r.URL.Query().Get("dt0")
	if dt0Str == "" {
		params.Dt0 = time.Now().AddDate(0, 0, 1)
	} else {
		params.Dt0, _ = time.Parse("02.01.2006", dt0Str)
	}

	dt1Str := r.URL.Query().Get("dt1")
	if dt1Str == "" {
		params.Dt1 = time.Now().AddDate(0, 0, 5)
	} else {
		params.Dt1, _ = time.Parse("02.01.2006", dt1Str)
	}

	tripsReturn, err := h.service.GetTrainRoutesReturn(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, tripsReturn)
}

func (h *Handlers) GetTrainCarriages(w http.ResponseWriter, r *http.Request) {
	var params domain.TrainCarriagesParams
	// Парсинг параметров из запроса
	params.Dir = parseInt(r, "dir", 0)
	params.Code0 = r.URL.Query().Get("code0")
	params.Code1 = r.URL.Query().Get("code1")

	dt0Str := r.URL.Query().Get("dt0")
	if dt0Str == "" {
		params.Dt0 = time.Now().AddDate(0, 0, 1)
	} else {
		params.Dt0, _ = time.Parse("02.01.2006", dt0Str)
	}

	time0Str := r.URL.Query().Get("time0")
	if time0Str == "" {
		params.Time0 = time.Now()
	} else {
		params.Time0, _ = time.Parse("15:04", time0Str)
	}

	params.Tnum0 = r.URL.Query().Get("tnum0")

	carriages, err := h.service.GetTrainCarriages(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, carriages)
}

func (h *Handlers) GetTrainStationList(w http.ResponseWriter, r *http.Request) {
	var params domain.TrainStationListParams
	// Парсинг параметров из запроса
	params.TrainNumber = r.URL.Query().Get("trainNumber")

	depDateStr := r.URL.Query().Get("depDate")
	if depDateStr == "" {
		params.DepDate = time.Now().AddDate(0, 0, 1)
	} else {
		params.DepDate, _ = time.Parse("02.01.2006", depDateStr)
	}

	stations, err := h.service.GetTrainStationList(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, stations)
}

func (h *Handlers) GetStationCode(w http.ResponseWriter, r *http.Request) {
	var params domain.StationCodeParams
	// Парсинг параметров из запроса
	params.StationNamePart = r.URL.Query().Get("stationNamePart")
	params.CompactMode = r.URL.Query().Get("compactMode")

	stationCodes, err := h.service.GetStationCodes(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, stationCodes)
}

// Вспомогательные функции

func parseInt(r *http.Request, key string, defaultVal int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Ошибка сериализации JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
