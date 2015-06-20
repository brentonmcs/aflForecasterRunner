package AflForecaster
import (
	"io"
	"net/http"
	"log"
	"encoding/json"
	"os"
	"fmt"
)

func AddHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
func StatsResult(w http.ResponseWriter, r *http.Request) {
	AddHeaders(w)
	result, _ := json.Marshal(GenerateStats())
	io.WriteString(w, string(result))
}

func currentRound(w http.ResponseWriter, r *http.Request) {
	AddHeaders(w)
	result, _ := json.Marshal(GetCurrentRoundDetails())
	io.WriteString(w, string(result))
}

func currentRoundStats(w http.ResponseWriter, r *http.Request) {
	AddHeaders(w)
	result, _ := json.Marshal(GenerateCurrentRoundStats())
	io.WriteString(w, string(result))

}

func currentPrices(w http.ResponseWriter, r *http.Request) {
	AddHeaders(w)
	result, _ := json.Marshal(getCurrentRoundPrices())
	io.WriteString(w, string(result))

}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func StartHttpServer() {
	http.HandleFunc("/stats", StatsResult)
	http.HandleFunc("/currentRound", currentRound)
	http.HandleFunc("/currentRoundStats", currentRoundStats)
	http.HandleFunc("/prices", currentPrices)

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
