package AflForecaster
import (
	"strings"
	"strconv"
	"log"
)

func parseForecast(forecast, percentageStr, resultStr string, i, year int) ForecastModel {
	forecastSplit := strings.Split(forecast, "by")
	winTeam := strings.Trim(forecastSplit[0], " ")

	var resultModel ResultModel
	if (len(resultStr) > 10 ) {
		resultModel = getResult(resultStr, winTeam)
	}
	return ForecastModel{
		WinTeam: winTeam,
		WinPercentage: getPercentage(percentageStr, winTeam),
		WinPoints: int(getPoints(forecastSplit[1])),
		Round: i,
		Year: year,
		ResultModel: resultModel}
}

func getPercentage(percentageStr, winTeam string) float32 {
	percentageSplit := strings.SplitAfter(percentageStr, "%")

	var teamPercentageIndex int

	if (strings.Contains(percentageSplit[0], winTeam)) {
		teamPercentageIndex = 0
	} else {
		teamPercentageIndex = 1
	}

	extractedPercentage := strings.Split(percentageSplit[teamPercentageIndex], " ")

	percentage, err := strconv.ParseFloat(strings.Trim(extractedPercentage[len(extractedPercentage) -1], "%"), 32)
	if err != nil {
		log.Fatal(err)
	}
	return float32(percentage)
}

func getPoints(pointsStr string) int {
	pointsSplit := strings.Split(strings.Trim(pointsStr, " "), " ")

	points, err := strconv.ParseInt(pointsSplit[0], 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	return int(points)
}

func getResult(resultStr string, winTeam string) ResultModel {
	bySplit := strings.Split(resultStr, "by")
	won := strings.Trim(bySplit[0], " ")
	return ResultModel{Won : won == winTeam, WinPoints: int(getPoints(bySplit[1])), WinTeam: won}
}
