package AflForecaster
import (
	"log"
)

var pointLow, betCount, wonCount, won40plus, wonUnder40, lose40plus, loseUnder40  int

func GenerateStats() []StatsModel {
	i := 0
	var result []StatsModel

	groupByPoint := groupMatchesByPoints()

	for _, p := range groupByPoint {
		if i == 0 {
			pointLow = p.Point
		}

		betCount += p.BetTotal
		wonCount += p.WonCount
		won40plus += p.WonOver40Count
		wonUnder40 += p.WonUnder40Count
		lose40plus += p.LoseOver40Count
		loseUnder40 += p.LoseUnder40Count

		i++;
		if i > 7 {
			result = addToArray(result, p.Point)
			reset();
			i = 0
		}
	}

	// Add Leftover results
	if (i > 1) {
		result = addToArray(result, groupByPoint[len(groupByPoint)-1].Point)
	}
	return result
}


func GenerateCurrentRoundStats() []MatchStatsPriceModel {
	currentRoundPrices := getCurrentRoundPrices()
	currentRound := GetCurrentRoundDetails()
	stats := GenerateStats()

	var result []MatchStatsPriceModel

	for _, cur := range currentRound {
		matchPrices := getRoundPrices(cur.WinTeam, currentRoundPrices)

		if (matchPrices == FavPrices{}) {
			continue
		}
		result = append(result, getStakingInformation(matchPrices, findStat(stats, cur.WinPoints), cur))
	}
	return result
}

func getRoundPrices(winTeam string, prices []MatchPrices) FavPrices {

	winTeam = convertName(winTeam)
	for _, r := range prices {

		if (r.HomeTeam.Name == winTeam) {
			return FavPrices{Favourite : r.HomeTeam, OtherTeam: r.AwayTeam}
		}
		if (r.AwayTeam.Name == winTeam) {
			return FavPrices{Favourite : r.AwayTeam, OtherTeam: r.HomeTeam}
		}
	}

	log.Println("Teams " + winTeam + " Prices Not Found")
	return FavPrices{}
}

func getStakingInformation(matchPrices FavPrices, stats StatsModel, forecast ForecastModel) MatchStatsPriceModel {

	return MatchStatsPriceModel{FavPrices:matchPrices,
		Stats:stats,
		FavouriteStakes: getTeamStakes(matchPrices.Favourite, stats.WinStat),
		OtherTeamStakes: getTeamStakes(matchPrices.OtherTeam, stats.LoseStat),
		Forecast: forecast}
}

func getTeamStakes(prices PriceModel, stat Stat) StakeModel {

	return StakeModel{HeadToHead: kellyCriterion(prices.HeadToHead, stat.WinPercentage),
		Under40: kellyCriterion(prices.HeadToHead, stat.WinUnder40Percentage),
		Over40: kellyCriterion(prices.HeadToHead, stat.WinOver40Percentage) }
}

func kellyCriterion(price float32, percentage float32) float32 {
	return ((percentage * price - 1) / (price - 1));
}

func findStat(stats []StatsModel, winPoints int) StatsModel {
	for _, s := range stats {
		if (s.PointHigh >= winPoints && s.PointLow <= winPoints) {
			return s;
		}
	}
	log.Fatal("Teams Stats Not Found")
	return StatsModel{}
}

func reset() {
	betCount = 0
	wonCount = 0
	won40plus = 0
	wonUnder40 = 0
	lose40plus = 0
	loseUnder40 = 0
}

func convertName(winTeam string) string {
	if (winTeam == "GW Sydney") {
		return "Greater Western Sydney";
	}

	if (winTeam == "Wstn Bulldogs") {
		return "Western Bulldogs";
	}
	return winTeam;
}

func addToArray(result []StatsModel, pointHigh int) []StatsModel {
	return append(result,
		StatsModel{PointHigh: pointLow,
			PointLow: pointHigh,
			BetCount: betCount,

			WinStat: Stat{
				WonCount: wonCount,
				Under40: wonUnder40,
				Plus40 : won40plus,
				WinPercentage: float32(wonCount) / float32(betCount),
				WinOver40Percentage: float32(won40plus) / float32(betCount),
				WinUnder40Percentage: float32(wonUnder40) / float32(betCount),
			},
			LoseStat: Stat{
				WonCount: betCount-wonCount,
				Under40: loseUnder40,
				Plus40 : lose40plus,
				WinPercentage: float32(betCount-wonCount) / float32(betCount),
				WinOver40Percentage: float32(lose40plus) / float32(betCount),
				WinUnder40Percentage: float32(loseUnder40) / float32(betCount) }})
}

