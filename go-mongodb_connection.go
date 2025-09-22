package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ====================
// Employee Struct
// ====================
type Employee struct {
	EmpID  int     `bson:"emp_id"`
	Name   string  `bson:"name"`
	Salary float64 `bson:"salary"`
}

// ====================
// Load .env File
// ====================
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables.")
	}
}

// ====================
// Connect to MongoDB
// ====================
func connectMongoDB() (*mongo.Client, context.Context, context.CancelFunc, error) {
	loadEnv()
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, nil, nil, fmt.Errorf("❌ MONGO_URI is not set in .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		cancel()
		_ = client.Disconnect(ctx)
		return nil, nil, nil, err
	}

	fmt.Println("✅ Connected to MongoDB")
	return client, ctx, cancel, nil
}

// ====================
// Insert Employees
// ====================
func insertEmployees(coll *mongo.Collection, employees []Employee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	docs := make([]interface{}, len(employees))
	for i, e := range employees {
		docs[i] = e
	}

	result, err := coll.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	fmt.Println("\n🟢 Inserted Employees with MongoDB IDs:")
	for _, id := range result.InsertedIDs {
		fmt.Println(" -", id)
	}
	return nil
}

// ====================
// Fetch Employees
// ====================
func fetchEmployees(coll *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	fmt.Println("\n📋 Employee Records:")
	for cursor.Next(ctx) {
		var emp Employee
		if err := cursor.Decode(&emp); err != nil {
			return err
		}
		fmt.Printf("ID: %d | Name: %s | Salary: %.2f\n", emp.EmpID, emp.Name, emp.Salary)
	}

	if err := cursor.Err(); err != nil {
		return err
	}
	return nil
}

// ====================
// Main Function
// ====================
func main() {
	client, ctx, cancel, err := connectMongoDB()
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	defer cancel()
	defer client.Disconnect(ctx)

	collection := client.Database("company_db").Collection("employees")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose action: 1) Insert Employees  2) Fetch Employees")
	fmt.Print("Enter your choice: ")
	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	switch choice {
	case "1":
		fmt.Print("Enter number of employees to insert: ")
		numStr, _ := reader.ReadString('\n')
		count, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil || count <= 0 {
			log.Fatal("❌ Invalid number of employees.")
		}

		employees := make([]Employee, count)
		for i := 0; i < count; i++ {
			fmt.Printf("\nEmployee %d:\n", i+1)

			fmt.Print("EmpID: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))

			fmt.Print("Name: ")
			name, _ := reader.ReadString('\n')

			fmt.Print("Salary: ")
			salStr, _ := reader.ReadString('\n')
			salary, _ := strconv.ParseFloat(strings.TrimSpace(salStr), 64)

			employees[i] = Employee{
				EmpID:  id,
				Name:   strings.TrimSpace(name),
				Salary: salary,
			}
		}

		if err := insertEmployees(collection, employees); err != nil {
			log.Fatal("❌ Insert failed:", err)
		}

	case "2":
		if err := fetchEmployees(collection); err != nil {
			log.Fatal("❌ Fetch failed:", err)
		}

	default:
		fmt.Println("❌ Invalid choice. Please enter 1 or 2.")
	}
}
