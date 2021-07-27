package main

import (
	"os"

	_ "gorm.io/driver/sqlite"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_TYPE"),
		os.Getenv("APP_DB_URI"),
	)
	a.DB.AutoMigrate(&product{})

	a.Run(":8010")
}
