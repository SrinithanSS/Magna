package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//// ===== Structs =====

// Employee
type Employee struct {
	ID     int     `bson:"id"`
	Name   string  `bson:"name"`
	Salary float64 `bson:"salary"`
}

// Department
type Department struct {
	Name  string `bson:"name"`
	EmpID int    `bson:"emp_id"`
}

// Developer
type Developer struct {
	Language string `bson:"language"`
	EmpID    int    `bson:"emp_id"`
}

// Tester
type Tester struct {
	Language string `bson:"language"`
	EmpID    int    `bson:"emp_id"`
}

// Aggregated result with JOIN
type EmployeeFull struct {
	ID          int          `bson:"id"`
	Name        string       `bson:"name"`
	Salary      float64      `bson:"salary"`
	Departments []Department `bson:"department_info"`
	Developers  []Developer  `bson:"developer_info"`
	Testers     []Tester     `bson:"tester_info"`
}

// // ===== Mongo Connection =====
func connectMongoDB(username, password string) (*mongo.Client, context.Context, context.CancelFunc) {
	encodedUser := url.QueryEscape(username)
	encodedPass := url.QueryEscape(password)

	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.im3q7gc.mongodb.net/", encodedUser, encodedPass)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Ping failed:", err)
	}
	fmt.Println("✅ Connected to MongoDB")
	return client, ctx, cancel
}

//// ===== CRUD Operations =====

// Insert new employee (with dept/dev/tester)
func insertEmployee(db *mongo.Database, ctx context.Context) {
	var emp Employee
	var dept Department
	var dev Developer
	var tester Tester

	fmt.Print("Enter Employee ID: ")
	fmt.Scan(&emp.ID)
	fmt.Print("Enter Name: ")
	fmt.Scan(&emp.Name)
	fmt.Print("Enter Salary: ")
	fmt.Scan(&emp.Salary)

	fmt.Print("Enter Department: ")
	fmt.Scan(&dept.Name)
	dept.EmpID = emp.ID

	fmt.Print("Enter Developer Language: ")
	fmt.Scan(&dev.Language)
	dev.EmpID = emp.ID

	fmt.Print("Enter Tester Language: ")
	fmt.Scan(&tester.Language)
	tester.EmpID = emp.ID

	_, _ = db.Collection("Employee").InsertOne(ctx, emp)
	_, _ = db.Collection("Department").InsertOne(ctx, dept)
	_, _ = db.Collection("Developer").InsertOne(ctx, dev)
	_, _ = db.Collection("Tester").InsertOne(ctx, tester)

	fmt.Println("✅ Employee inserted successfully")
}

// Update employee (only name & salary here)
func updateEmployee(db *mongo.Database, ctx context.Context) {
	var empID int
	var newName string
	var newSalary float64

	fmt.Print("Enter Employee ID to update: ")
	fmt.Scan(&empID)
	fmt.Print("Enter new Name: ")
	fmt.Scan(&newName)
	fmt.Print("Enter new Salary: ")
	fmt.Scan(&newSalary)

	filter := bson.M{"id": empID}
	update := bson.M{"$set": bson.M{"name": newName, "salary": newSalary}}

	_, err := db.Collection("Employee").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("❌ Update failed:", err)
	} else {
		fmt.Println("✅ Employee updated successfully")
	}
}

// Delete employee (cascade delete across collections)
func deleteEmployee(db *mongo.Database, ctx context.Context) {
	var empID int
	fmt.Print("Enter Employee ID to delete: ")
	fmt.Scan(&empID)

	collections := []string{"Employee", "Department", "Developer", "Tester"}
	for _, coll := range collections {
		_, _ = db.Collection(coll).DeleteOne(ctx, bson.M{"emp_id": empID})
		_, _ = db.Collection(coll).DeleteOne(ctx, bson.M{"id": empID}) // for Employee
	}
	fmt.Println("✅ Employee deleted successfully")
}

// Read with $lookup JOIN
func readEmployeesWithJoin(db *mongo.Database, ctx context.Context) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Department"},
			{Key: "localField", Value: "id"},
			{Key: "foreignField", Value: "emp_id"},
			{Key: "as", Value: "department_info"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Developer"},
			{Key: "localField", Value: "id"},
			{Key: "foreignField", Value: "emp_id"},
			{Key: "as", Value: "developer_info"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Tester"},
			{Key: "localField", Value: "id"},
			{Key: "foreignField", Value: "emp_id"},
			{Key: "as", Value: "tester_info"},
		}}},
	}

	cursor, err := db.Collection("Employee").Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal("Aggregation failed:", err)
	}
	defer cursor.Close(ctx)

	fmt.Println("\n=== Employee Records with Full Details ===")
	for cursor.Next(ctx) {
		var emp EmployeeFull
		if err := cursor.Decode(&emp); err != nil {
			log.Println("Decode failed:", err)
			continue
		}
		jsonData, _ := json.MarshalIndent(emp, "", "  ")
		fmt.Println(string(jsonData))
	}
}

// // ===== Sample Data Loader =====
func loadSampleData(db *mongo.Database, ctx context.Context) {
	// Drop old data
	db.Collection("Employee").Drop(ctx)
	db.Collection("Department").Drop(ctx)
	db.Collection("Developer").Drop(ctx)
	db.Collection("Tester").Drop(ctx)

	// Insert sample employees
	emps := []interface{}{
		Employee{ID: 1, Name: "Alice", Salary: 50000},
		Employee{ID: 2, Name: "Bob", Salary: 60000},
		Employee{ID: 3, Name: "Charlie", Salary: 55000},
	}
	depts := []interface{}{
		Department{Name: "IT", EmpID: 1},
		Department{Name: "HR", EmpID: 2},
		Department{Name: "Finance", EmpID: 3},
	}
	devs := []interface{}{
		Developer{Language: "Go", EmpID: 1},
		Developer{Language: "Python", EmpID: 2},
		Developer{Language: "Java", EmpID: 3},
	}
	tests := []interface{}{
		Tester{Language: "JavaScript", EmpID: 1},
		Tester{Language: "Ruby", EmpID: 2},
		Tester{Language: "C#", EmpID: 3},
	}

	db.Collection("Employee").InsertMany(ctx, emps)
	db.Collection("Department").InsertMany(ctx, depts)
	db.Collection("Developer").InsertMany(ctx, devs)
	db.Collection("Tester").InsertMany(ctx, tests)

	fmt.Println("✅ Sample data loaded")
}

// // ===== Main Menu =====
func main() {
	var user, pass string
	fmt.Print("Enter MongoDB Username: ")
	fmt.Scan(&user)
	fmt.Print("Enter MongoDB Password: ")
	fmt.Scan(&pass)

	client, ctx, cancel := connectMongoDB(user, pass)
	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database("unified_demo")

	// Load sample data once
	loadSampleData(db, ctx)

	for {
		fmt.Println("\n===== MENU =====")
		fmt.Println("1. Insert Employee")
		fmt.Println("2. Update Employee")
		fmt.Println("3. Delete Employee")
		fmt.Println("4. Read Employees with Join")
		fmt.Println("5. Exit")
		fmt.Print("Choose option: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			insertEmployee(db, ctx)
		case 2:
			updateEmployee(db, ctx)
		case 3:
			deleteEmployee(db, ctx)
		case 4:
			readEmployeesWithJoin(db, ctx)
		case 5:
			fmt.Println("👋 Exiting...")
			return
		default:
			fmt.Println("❌ Invalid choice, try again.")
		}
	}
}
