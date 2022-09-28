package database

import (
	"directory/alerting"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Address struct {
	//gorm.Model
	Id         int `gorm:"PRIMARY_KEY"`
	Area       string
	City       string
	Zip        int
	EmployeeID int `gorm:"column:employee_id"`
}

type Employee struct {
	// gorm.Model
	Name       string
	Age        int
	Salary     int
	YOE        int
	EmployeeID int       `gorm:"PRIMARY_KEY"`
	Address    []Address `gorm:"foreignkey:employee_id" constraint:"OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

const (
	dsn = `host=localhost 
	user=akash 
	password=password 
	dbname=golangdb1 
	port=5432 
	sslmode=disable 
	TimeZone=Asia/Shanghai`
)

func (e *Employee) AfterCreate(tx *gorm.DB) (err error) {
	alerting.SendWelcomeEmail(e.Name)
	return
}

func EstablishDBConn() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("DB Connection Error %v \n", err)
	} else {
		db.AutoMigrate(&Employee{})
		db.AutoMigrate(&Address{})
		db.Set("gorm:table_options",
			"ENGINE=Distributed(cluster, default, hits)").
			AutoMigrate(&Employee{})
	}
	return db, err
}

func Create(db *gorm.DB, employee *Employee) (*gorm.DB, error) {
	result := db.Create(&employee)
	return result, result.Error
}

func Update(db *gorm.DB, employee *Employee) (Employee, error) {
	res := db.Model(&Employee{}).
		Where("employee_id = ?", employee.EmployeeID).
		Updates(employee)
	return *employee, res.Error
}

func DeleteById(db *gorm.DB, employeeId int) (*gorm.DB, error) {

	fmt.Println("Deleting ===>", employeeId)
	y := db.Where("employee_id=?", employeeId).
		Delete(&Address{})

	db.Model(&Employee{}).
		Association("Address").Clear()

	x := db.Where("employee_id = ?", employeeId).
		Delete(&Employee{})

	fmt.Println("rows affected ", x.RowsAffected, y.RowsAffected)
	return x, x.Error
}

func ReadAll(db *gorm.DB) ([]Employee, error) {
	var employees []Employee
	err := db.Model(&Employee{}).
		Preload("Address").
		Find(&employees).
		Error
	return employees, err
}

func ReadById(db *gorm.DB, employeeId int) (Employee, error) {
	var employee Employee
	res := db.First(&employee, "employee_id = ?", employeeId)
	return employee, res.Error
}
