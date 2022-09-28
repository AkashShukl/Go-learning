package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"directory/database"
	pb "directory/directoryProto"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

const (
	port = ":50051"
)

type EmployeeDirectoryServer struct {
	db *gorm.DB
	pb.UnimplementedEmployeeDirectoryServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	db, err := database.EstablishDBConn()
	if err != nil {
		fmt.Println("check Database connection!!")
	} else {
		fmt.Println("Connection to Database Established!")
	}

	pb.RegisterEmployeeDirectoryServer(s, &EmployeeDirectoryServer{db: db})
	log.Printf("Server Listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *EmployeeDirectoryServer) DeleteEmployee(ctx context.Context, in *pb.ID) (*pb.ID, error) {
	db := s.db
	res, err := database.DeleteById(db, int(in.GetId()))
	fmt.Println("DEL", res.RowsAffected, err)
	return &pb.ID{Id: in.GetId()}, err
}

func (s *EmployeeDirectoryServer) GetEmployeeById(ctx context.Context, in *pb.ID) (*pb.Employee, error) {
	db := s.db
	res, err := database.ReadById(db, int(in.GetId()))
	return &pb.Employee{Name: res.Name,
		Salary: int32(res.Salary),
		Age:    int32(res.Age),
		Yoe:    int32(res.YOE)}, err
}

func (s *EmployeeDirectoryServer) CreateNewEmployee(ctx context.Context, in *pb.Employee) (*pb.Employee, error) {

	address := []database.Address{}
	for _, y := range in.Address {
		fmt.Println(y.GetEmployeeID())
		a := database.Address{
			Area:       y.GetArea(),
			City:       y.GetCity(),
			Zip:        int(y.GetZip()),
			EmployeeID: int(in.GetEmployeeID()),
		}
		address = append(address, a)
	}
	e := database.Employee{
		Name:       in.Name,
		Age:        int(in.Age),
		Salary:     int(in.Salary),
		YOE:        int(in.Yoe),
		Address:    address,
		EmployeeID: int(in.EmployeeID),
	}
	db := s.db
	result, createError := database.Create(db, &e)
	if createError != nil {
		return &pb.Employee{}, createError
	} else {
		fmt.Println("\nrows affected", result.RowsAffected)
	}
	return in, nil
}

func (s *EmployeeDirectoryServer) UpdateEmployee(ctx context.Context, in *pb.Employee) (*pb.Employee, error) {
	db := s.db
	address := []database.Address{}
	for _, y := range in.Address {
		fmt.Println(y.GetEmployeeID())
		a := database.Address{
			Area:       y.GetArea(),
			City:       y.GetCity(),
			Zip:        int(y.GetZip()),
			EmployeeID: int(in.GetEmployeeID()),
		}
		address = append(address, a)
	}
	e := database.Employee{
		Name:       in.Name,
		Age:        int(in.Age),
		Salary:     int(in.Salary),
		YOE:        int(in.Yoe),
		Address:    address,
		EmployeeID: int(in.EmployeeID),
	}
	_, err := database.Update(db, &e)
	return in, err
}
