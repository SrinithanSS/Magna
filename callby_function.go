package main

import (
	"fmt"
)

type Staff struct {
	ID     int
	FullName string
	Pay    float64
}

// Function demonstrating pass by value
func CallByValue(staffList [4]Staff) {
	staffList[0].ID = 1111
	staffList[1].FullName = "Updated Name"
	staffList[2].Pay = 9999.99
}

// Function demonstrating pass by reference
func CallByReference(staffSlice *[]Staff) {
	(*staffSlice)[0].ID = 2222
	(*staffSlice)[1].FullName = "Modified Person"
	(*staffSlice)[2].Pay = 8888.88
	(*staffSlice)[3].Pay = 7777.77
}

// Reusable print function
func DisplayStaff(title string, staffArr [4]Staff) {
	fmt.Println("\n" + title)
	for _, s := range staffArr {
		fmt.Printf("ID: %d, Name: %s, Salary: %.2f\n", s.ID, s.FullName, s.Pay)
	}
}

func main() {
	team := [4]Staff{
		{ID: 1001, FullName: "Alice", Pay: 5000.00},
		{ID: 1002, FullName: "Bob", Pay: 6000.00},
		{ID: 1003, FullName: "Charlie", Pay: 7000.00},
		{ID: 1004, FullName: "Diana", Pay: 8000.00},
	}

	DisplayStaff("Before CallByValue:", team)

	CallByValue(team)

	DisplayStaff("After CallByValue (No Change Expected):", team)

	// Slice for passing by reference
	teamSlice := team[:]

	CallByReference(&teamSlice)

	DisplayStaff("After CallByReference (Changes Reflected):", team)
}
