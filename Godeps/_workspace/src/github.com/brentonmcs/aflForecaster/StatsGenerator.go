package AflForecaster

var pointLow = 0
var betCount = 0
var wonCount = 0
var won40plus = 0
var wonUnder40 = 0
var lose40plus = 0
var loseUnder40 = 0

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
		if i > 9 {
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

func reset() {
	betCount = 0
	wonCount = 0
	won40plus = 0
	wonUnder40 = 0
	lose40plus = 0
	loseUnder40 = 0
}

func addToArray(result []StatsModel, pointHigh int) []StatsModel {
	return append(result,
		StatsModel{PointHigh: pointHigh,
			PointLow: pointLow,
			BetCount: betCount,
			WonCount: wonCount,
			Won40Plus: won40plus,
			WonUnder40: wonUnder40,
			Lose40Plus: lose40plus,
			LoseUnder40:loseUnder40,
			WinPercentage: float32(wonCount) / float32(betCount),
			LosePercentage: float32(betCount-wonCount) / float32(betCount),
			WinOver40Percentage: float32(won40plus) / float32(betCount),
			WinUnder40Percentage: float32(wonUnder40) / float32(betCount),
			LoseOver40Percentage: float32(lose40plus) / float32(betCount),
			LoseUnder40Percentage: float32(loseUnder40) / float32(betCount),
		})
}

