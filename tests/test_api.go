// tests/test_api.go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Chaika-Team/rzd-api/internal/adapters/api"
	"github.com/Chaika-Team/rzd-api/internal/config"
	"github.com/Chaika-Team/rzd-api/internal/domain"
	"github.com/Chaika-Team/rzd-api/internal/infrastructure"
	"github.com/go-kit/log"
)

func main() {
	logger := log.NewNopLogger() // или log.NewLogfmtLogger(os.Stdout)
	cfg := config.GetConfig(logger, "config.yaml")

	httpClient := infrastructure.NewGuzzleHttpClient(cfg, logger)
	rzdAPI := api.NewRzdAPI(httpClient, cfg, logger)

	ctx := context.Background()

	// 1. TrainRoutes
	fmt.Println("=== TrainRoutes ===")
	tr, err := rzdAPI.TrainRoutes(ctx, domain.TrainRoutesParams{
		Dir:        0,
		Tfl:        3,
		CheckSeats: 1,
		Code0:      "2004000", // СПб
		Code1:      "2000000", // Москва
		Dt0:        time.Now().AddDate(0, 0, 1),
		Md:         0,
	})
	if err != nil {
		fmt.Println("TrainRoutes error:", err)
	} else {
		fmt.Println("TrainRoutes result:", tr)
	}

	time.Sleep(1 * time.Second)

	// 2. TrainRoutesReturn
	fmt.Println("=== TrainRoutesReturn ===")
	trr, err := rzdAPI.TrainRoutesReturn(ctx, domain.TrainRoutesReturnParams{
		Dir:        1,
		Tfl:        3,
		CheckSeats: 1,
		Code0:      "2004000",
		Code1:      "2000000",
		Dt0:        time.Now().AddDate(0, 0, 1),
		Dt1:        time.Now().AddDate(0, 0, 5),
	})
	if err != nil {
		fmt.Println("TrainRoutesReturn error:", err)
	} else {
		fmt.Println("TrainRoutesReturn result:", trr)
	}

	time.Sleep(1 * time.Second)
	// 3. TrainCarriages
	fmt.Println("=== TrainCarriages ===")
	// чтобы carriages было на что смотреть, нужно сначала получить routes
	// Но для теста можно просто проверить сам запрос:
	tc, err := rzdAPI.TrainCarriages(ctx, domain.TrainCarriagesParams{
		Dir:   0,
		Code0: "2004000",
		Code1: "2000000",
		Dt0:   time.Now().AddDate(0, 0, 1),
		Time0: time.Date(2025, 1, 1, 15, 30, 0, 0, time.UTC),
		Tnum0: "072E",
	})
	if err != nil {
		fmt.Println("TrainCarriages error:", err)
	} else {
		fmt.Println("TrainCarriages result:", tc)
	}

	time.Sleep(1 * time.Second)
	// 4. TrainStationList
	fmt.Println("=== TrainStationList ===")
	ts, err := rzdAPI.TrainStationList(ctx, domain.TrainStationListParams{
		TrainNumber: "054Г",
		DepDate:     time.Now().AddDate(0, 0, 1),
	})
	if err != nil {
		fmt.Println("TrainStationList error:", err)
	} else {
		fmt.Println("TrainStationList result:", ts)
	}

	time.Sleep(1 * time.Second)
	// 5. StationCode
	fmt.Println("=== StationCode ===")
	sc, err := rzdAPI.StationCode(ctx, domain.StationCodeParams{
		StationNamePart: "ЧЕБ",
		CompactMode:     "y",
	})
	if err != nil {
		fmt.Println("StationCode error:", err)
	} else {
		fmt.Println("StationCode result:", sc)
	}
}
