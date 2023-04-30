package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type Transaction struct {
	CCNum      string  `bson:"ccnum"`
	Date       string  `bson:"amount"`
	Amount     float32 `bson:"amount"`
	Cvv        string  `bson:"cvv"`
	Expiration string  `bson:"exp"`
}

func loot_mongo(wg *sync.WaitGroup) {
	defer wg.Done()
	maxWait := time.Duration(5 * time.Second)
	session, err := mgo.DialWithTimeout("127.0.0.1:27017", maxWait)
	if err != nil {
		log.Panicln("Error connecting to server", err)
	}

	results := make([]Transaction, 0)
	if err := session.DB("test").C("transations").Find(nil).All(&results); err != nil {
		log.Panicln(err)
	}

	for _, txn := range results {
		fmt.Println(txn.CCNum, txn.Date, txn.Amount, txn.Cvv, txn.Expiration)
	}
}

func loot_mysql(wg *sync.WaitGroup) {
	defer wg.Done()
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/store")
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	var (
		ccnum, date, cvv, exp string
		amount                float32
	)
	rows, err := db.Query("SELECT ccnum,date,amount,cvv,exp FROM transactions")
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ccnum, &date, &amount, &cvv, &exp)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(ccnum, date, amount, cvv, exp)
	}
	if rows.Err() != nil {
		log.Panicln(err)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)
	go loot_mongo(&wg)
	go loot_mysql(&wg)

	wg.Wait()
}
