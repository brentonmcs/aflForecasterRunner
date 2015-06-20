package AflForecaster

type AggregatePoints struct {
	Point      int `bson:"_id" json:"_id"`
	WonCount int `bson:"wonCount" json:"wonCount"`
	BetTotal int `bson:"betTotal" json:"betTotal"`
	WonOver40Count int `bson:"wonOver40Count" json:"wonOver40Count"`
	WonUnder40Count int `bson:"wonUnder40Count" json:"wonUnder40Count"`
	LoseOver40Count int `bson:"loseOver40Count" json:"loseOver40Count"`
	LoseUnder40Count int `bson:"loseUnder40Count" json:"loseUnder40Count"`
}

type ActiveRound struct {
	round int
	year  int
}

type ForecastModel struct {
	WinTeam       string
	WinPercentage float32
	WinPoints     int
	Round         int
	ResultModel   ResultModel
	Year          int
}

type ResultModel struct {
	Won       bool
	WinTeam   string
	WinPoints int
}

type StatsModel struct {
	PointLow              int
	PointHigh             int
	BetCount              int
	WonCount              int
	Won40Plus             int
	WonUnder40            int
	Lose40Plus            int
	LoseUnder40           int
	WinPercentage         float32
	LosePercentage        float32
	WinOver40Percentage   float32
	WinUnder40Percentage  float32
	LoseOver40Percentage  float32
	LoseUnder40Percentage float32
}
