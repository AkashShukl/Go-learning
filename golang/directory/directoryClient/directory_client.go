package main

import (
	"bufio"
	pb "directory/directoryProto"
	"fmt"
	"log"
	"os"
	"time"

	"context"

	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect server: %v", err)
	}
	defer conn.Close()
	c := pb.NewEmployeeDirectoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var eid int32 = 222

	fmt.Println(c)

	fmt.Print("1 to create an item demo; 2 to read item; 3 to update and 4 to remove: ")
	choice := bufio.NewReader(os.Stdin)
	text, _ := choice.ReadString('\n')
	switch text {
	case "1\n":
		item, err1 := c.CreateNewEmployee(ctx, &pb.Employee{
			Name: "jhonn Allen",
			Age:  20, Salary: 1000, Yoe: 2,
			EmployeeID: eid, Address: []*pb.Address{
				{
					City:       "donce",
					Zip:        34512,
					Area:       "code",
					EmployeeID: eid,
				}}})
		if err1 != nil {
			log.Fatalf("\nCould not create a new item: %v", err)
		} else {
			fmt.Println("\nInsertion Successful->", item.Name)
		}

	case "2\n":
		res, err2 := c.GetEmployeeById(ctx, &pb.ID{Id: eid})
		if err2 != nil {
			log.Fatalf("\nCould not search item: %v", err2)
		} else {
			fmt.Println("\nSearch Successful->", res)
		}

	case "3\n":
		item, err0 := c.UpdateEmployee(ctx, &pb.Employee{
			EmployeeID: eid,
			Name:       "updated name",
		})
		if err0 != nil {
			log.Fatalf("\nCould not update item: %v", err)
		} else {
			fmt.Println("\nUpdation Successful->", item)
		}

	case "4\n":
		res3, err3 := c.DeleteEmployee(ctx, &pb.ID{Id: eid})
		if err3 != nil {
			log.Fatalf("\nCould not delete item: %v", err3)
		} else {
			fmt.Println("\nDelete Successful->", res3)
		}
	default:
		fmt.Println("\nWrong option!")
	}
}
