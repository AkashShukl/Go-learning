package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Area string
	City string
	Zip  int
}

type Employee struct {
	gorm.Model
	Name    string
	Address string
	Salary  int
	YOE     int
}

func main() {
	dsn := "host=localhost user=akash password=password dbname=golangdb1 port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("DB Connection Error %v \n", err)
	} else {
		db.AutoMigrate(&Employee{})
		// Set table options
		db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&Employee{})

		// Insert
		employee := Employee{
			Name:    "Sasssm",
			Address: "st 12 ghj sasjd",
			Salary:  111,
			YOE:     2,
		}
		result := db.Create(&employee)
		if result.Error != nil {
			fmt.Println("error occured while injecting data", result.Error)
		} else {
			fmt.Println(result.RowsAffected)
		}

		// Batch Insert
		// var employees = []Employee{employee1, employee2, employee3}
		// db.Create(&employees)

		searchedEmployee := Employee{}
		db.Find(&searchedEmployee, "id = ?", 1)
		fmt.Println("seaerch result -> ", searchedEmployee)

	}

}
