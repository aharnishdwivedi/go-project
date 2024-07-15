package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

// our struct as per discussion
type AutoGenerated struct {
	ID        int    `json:"id"`
	TempValue int    `json:"tempValue"`
	TimeStamp string `json:"timeStamp"`
}

func main() {

	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every("1s").Do(sendApiCalls)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	// Start the scheduler asynchronously
	fmt.Println("Scheduler started")
	s.StartAsync()
	select {} //to keep main alive indefinately
}

func sendApiCalls() {
	var wg sync.WaitGroup
	fmt.Println("sendApiCalls triggered")

	// we are looping to send api calls
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			res, err := getTempReadings(id)
			if err != nil {
				fmt.Println("Error getting temp readings for ID:", id, err)
				return
			}
			// Print our respnose
			fmt.Printf("Sensor ID: %d, TempValue: %d, TimeStamp: %s\n", res.ID, res.TempValue, res.TimeStamp)
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func getTempReadings(id int) (AutoGenerated, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.google.com/v1/sensor/data/%d", id))
	if err != nil {
		// here if req fails we are sending hardcoed response
		fmt.Println("HTTP request failed for ID:", id, err)
		return AutoGenerated{
			ID:        id,
			TempValue: 25,
			TimeStamp: time.Now().Format(time.RFC3339),
		}, nil
	}
	defer resp.Body.Close()

	// Reading the response body using IO
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		// If reading the body fails, print an error and return hardcoed values here also
		fmt.Println("Reading response body failed for ID:", id, err)
		return AutoGenerated{
			ID:        id,
			TempValue: 25,
			TimeStamp: time.Now().Format(time.RFC3339),
		}, nil
	}

	// Return hardcoded values regardless of the actual response
	return AutoGenerated{
		ID:        id,
		TempValue: 25,
		TimeStamp: time.Now().Format(time.RFC3339),
	}, nil
}
