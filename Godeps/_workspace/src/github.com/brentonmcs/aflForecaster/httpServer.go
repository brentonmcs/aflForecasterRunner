package AflForecaster
import (
	"io"
	"net/http"
	"log"
	"encoding/json"
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

func StartHttpServer() {
	http.HandleFunc("/stats", StatsResult)
	http.HandleFunc("/currentRound", currentRound)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
