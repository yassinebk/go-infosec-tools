package main

import (
	dbminer "go-db-miner/dbminer"
	mongominer "go-db-miner/mongominer"
	mysqlminer "go-db-miner/mysqlminer"
	"os"
)

func main() {

	mongo_miner, err := mongominer.New(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = dbminer.Search(mongo_miner)

	if err != nil {
		panic(err)
	}

	sql_miner, err := mysqlminer.New(os.Args[1])

	if err != nil {
		panic(err)
	}

	err = dbminer.Search(sql_miner)

	if err != nil {
		panic(err)
	}

}
