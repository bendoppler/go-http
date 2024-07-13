package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Employee struct to hold individual employee data
type Employee struct {
	ID             int    `json:"id"`
	EmployeeName   string `json:"employee_name"`
	EmployeeSalary int    `json:"employee_salary"`
	EmployeeAge    int    `json:"employee_age"`
	ProfileImage   string `json:"profile_image"`
}

// Response struct to hold the API response
type Response struct {
	Status string     `json:"status"`
	Data   []Employee `json:"data"`
}

func main() {
	url := "https://dummy.restapiexample.com/api/v1/employees"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get a valid response: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	fmt.Printf("Status: %s\n", response.Status)
	for _, employee := range response.Data {
		fmt.Printf("ID: %d, Name: %s, Salary: %d, Age: %d, Profile Image: %s\n",
			employee.ID, employee.EmployeeName, employee.EmployeeSalary, employee.EmployeeAge, employee.ProfileImage)
	}
}
