package AflForecaster

import (
	"aflForecasterRunner/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"aflForecasterRunner/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"fmt"
	"log"
	"os"
	"time"
)

func AddForecast(forecast *ForecastModel) {

	session := connect()
	defer session.Close()

	c := session.DB("aflForecaster").C("forecast")

	index := mgo.Index{
		Key:        []string{"Round", "WinTeam", "Year"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Insert(&forecast)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(forecast)
}

func addPriceRecord(matchRecord *MatchPrices) {
	session := connect()
	defer session.Close()

	c := session.DB("aflForecaster").C("prices")

	err := c.Insert(&matchRecord)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(matchRecord)
}

func GetCurrentRoundDetails() []ForecastModel {
	session := connect()
	defer session.Close()

	year, round := getCurrentRound()

	c := session.DB("aflForecaster").C("forecast")
	query := bson.M{"year": year, "round": round}

	var result []ForecastModel
	c.Find(query).All(&result)

	return result
}

func getCurrentRound() (int, int) {
	session := connect()
	defer session.Close()

	c := session.DB("aflForecaster").C("forecast")

	var years, rounds []int
	err := c.Find(nil).Sort("-year").Distinct("year", &years)
	if err != nil {
		log.Fatal(err)
	}

	curYear := 2012
	curRound := 1
	if len(years) > 0 {
		curYear = years[0]

		err = c.Find(bson.M{"year": curYear}).Sort("-round").Distinct("round", &rounds)
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range rounds {
			if r > curRound {
				curRound = r
			}
		}
	}

	return curYear, curRound
}

func removeYearRoundMatches(year, round int) {
	session := connect()
	defer session.Close()

	c := session.DB("aflForecaster").C("forecast")

	query := bson.M{"year": year, "round": round}
	info, err := c.RemoveAll(query)

	if info.Removed == 0 {
		log.Fatal("Something should have been deleted")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func getActiveRounds() []ActiveRound {
	session := connect()
	defer session.Close()

	c := session.DB("aflForecaster").C("forecast")

	query := bson.M{"resultmodel.winteam": "", "round": bson.M{"$ne": 23}}

	var results []ForecastModel
	var activeRound []ActiveRound
	err := c.Find(query).Sort("year, round").All(&results)
	if err != nil {
		log.Fatal(err)
	}

	curYear := 0
	curRound := 0
	for _, e := range results {
		if curYear != e.Year || curRound != e.Round {
			curYear = e.Year
			curRound = e.Round

			activeRound = append(activeRound, ActiveRound{round: curRound, year: curYear})

			removeYearRoundMatches(curYear, curRound)
		}
	}

	return activeRound
}

func groupMatchesByPoints() []AggregatePoints {
	session := connect()
	defer session.Close()

	c := session.DB("aflForecaster").C("forecast")

	eqWon := bson.M{"$eq": []interface{}{"$resultmodel.won", true}}
	eqLose := bson.M{"$eq": []interface{}{"$resultmodel.won", false}}
	eqOver40 := bson.M{"$gte": []interface{}{"$resultmodel.winpoints", 40}}
	eqUnder40 := bson.M{"$lt": []interface{}{"$resultmodel.winpoints", 40}}

	and := bson.M{"$and": []interface{}{eqWon, eqOver40}}
	andUnder := bson.M{"$and": []interface{}{eqWon, eqUnder40}}
	andLose := bson.M{"$and": []interface{}{eqLose, eqOver40}}
	andUnderLose := bson.M{"$and": []interface{}{eqLose, eqUnder40}}

	query := []interface{}{
		bson.M{"$match": bson.M{"resultmodel.winteam": bson.M{"$ne": ""}}},
		bson.M{"$group": bson.M{
			"_id":              "$winpoints",
			"betTotal":         bson.M{"$sum": 1},
			"wonCount":         bson.M{"$sum": bson.M{"$cond": []interface{}{eqWon, 1, 0}}},
			"wonOver40Count":   bson.M{"$sum": bson.M{"$cond": []interface{}{and, 1, 0}}},
			"wonUnder40Count":  bson.M{"$sum": bson.M{"$cond": []interface{}{andUnder, 1, 0}}},
			"loseOver40Count":  bson.M{"$sum": bson.M{"$cond": []interface{}{andLose, 1, 0}}},
			"loseUnder40Count": bson.M{"$sum": bson.M{"$cond": []interface{}{andUnderLose, 1, 0}}}}},
		bson.M{"$sort": bson.M{"_id": -1}}}

	var results []AggregatePoints
	c.Pipe(query).All(&results)
	return results
}

func getCurrentRoundPrices() []MatchPrices {

	daysTilSunday := int(time.Saturday-time.Now().Weekday()) + 2

	duration := time.Duration(time.Duration(daysTilSunday*24) * time.Hour)
	currentDate := time.Now().Add(duration)

	fmt.Println(currentDate)
	session := connect()
	defer session.Close()

	c := session.DB("aflForecaster").C("prices")

	var result []MatchPrices
	c.Find(bson.M{"matchdate": bson.M{"$lt": currentDate}}).All(&result)
	return result
}

func clearPrices() {
	session := connect()
	defer session.Close()

	c := session.DB("aflForecaster").C("prices")

	count, err := c.Count()

	if count == 0 {
		return
	}
	info, err := c.RemoveAll(bson.M{})

	if info.Removed == 0 {
		log.Fatal("Something should have been deleted")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func connect() *mgo.Session {

	mongoConnectionString := os.Getenv("MONGOLAB_URI")
	if mongoConnectionString == "" {
		log.Fatal("Mongo Server not found")
	}

	session, err := mgo.Dial(mongoConnectionString)
	if err != nil {
		panic(err)
	}

	if false {
		mgo.SetDebug(true)

		var aLogger *log.Logger
		aLogger = log.New(os.Stderr, "", log.LstdFlags)
		mgo.SetLogger(aLogger)
	}

	session.SetMode(mgo.Monotonic, true)
	return session
}
