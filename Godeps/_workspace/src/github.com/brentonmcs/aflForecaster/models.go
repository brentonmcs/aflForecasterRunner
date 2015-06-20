package AflForecaster

import "time"

type AggregatePoints struct {
	Point            int `bson:"_id" json:"_id"`
	WonCount         int `bson:"wonCount" json:"wonCount"`
	BetTotal         int `bson:"betTotal" json:"betTotal"`
	WonOver40Count   int `bson:"wonOver40Count" json:"wonOver40Count"`
	WonUnder40Count  int `bson:"wonUnder40Count" json:"wonUnder40Count"`
	LoseOver40Count  int `bson:"loseOver40Count" json:"loseOver40Count"`
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
	PointLow  int
	PointHigh int
	BetCount  int
	WinStat   Stat
	LoseStat  Stat
}

type Stat struct {
	WonCount             int
	Under40              int
	Plus40               int
	WinPercentage        float32
	WinOver40Percentage  float32
	WinUnder40Percentage float32
}

type PriceModel struct {
	Name       string
	HeadToHead float32
	Under39    float32
	Over40     float32
}

type FavPrices struct {
	Favourite PriceModel
	OtherTeam PriceModel
}

type MatchPrices struct {
	HomeTeam  PriceModel
	AwayTeam  PriceModel
	MatchDate time.Time
}

type StakeModel struct {
	HeadToHead float32
	Under40    float32
	Over40     float32
}

type MatchStatsPriceModel struct {

	Stats           StatsModel
	FavPrices       FavPrices
	FavouriteStakes StakeModel
	OtherTeamStakes StakeModel
	Forecast        ForecastModel
}
