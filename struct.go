package main

import "fmt"

type Employee struct {
	ID     int
	Name   string
	Salary float64
}

func main() {
	var count int
	fmt.Print("Enter number of employees: ")
	fmt.Scan(&count)

	var employees []Employee

	for i := 0; i < count; i++ {
		var emp Employee
		fmt.Printf("\nEnter details for Employee %d\n", i+1)

		fmt.Print("Enter ID: ")
		fmt.Scan(&emp.ID)

		fmt.Print("Enter Name: ")
		fmt.Scan(&emp.Name)

		fmt.Print("Enter Salary: ")
		fmt.Scan(&emp.Salary)

		employees = append(employees, emp)
	}

	fmt.Println("\n--- Employee Details ---")
	for _, emp := range employees {
		fmt.Printf("ID: %d\n", emp.ID)
		fmt.Printf("Name: %s\n", emp.Name)
		fmt.Printf("Salary: %.2f\n", emp.Salary)
		fmt.Println("-------------------------------")
	}
}
