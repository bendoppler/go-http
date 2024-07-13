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

// Result struct to hold the computed result
type Result struct {
	ID     int
	Result float64
}

func worker(id int, jobs <-chan Employee, results chan<- Result, done chan<- bool) {
	for employee := range jobs {
		if employee.EmployeeAge != 0 {
			result := float64(employee.EmployeeSalary) / float64(employee.EmployeeAge)
			fmt.Printf("Worker %d processed employee %d\n", id, employee.ID)
			results <- Result{ID: employee.ID, Result: result}
		} else {
			fmt.Printf("Worker %d skipped employee %d due to zero age\n", id, employee.ID)
		}
	}
	done <- true
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

	const numWorkers = 3
	jobs := make(chan Employee, len(response.Data))
	results := make(chan Result, len(response.Data))
	done := make(chan bool, numWorkers)

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results, done)
	}

	// Send jobs
	for _, employee := range response.Data {
		jobs <- employee
	}
	close(jobs)

	// Wait for all workers to finish
	for i := 0; i < numWorkers; i++ {
		<-done
	}
	close(results)

	// Collect results
	for result := range results {
		fmt.Printf("Employee ID: %d, Salary/Age: %.2f\n", result.ID, result.Result)
	}
}
