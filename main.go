package main

import (
	"fmt"

	messagebird "github.com/messagebird/go-rest-api"
)

const apiKey = "0E8kldYVc5JwvaeXF2j0ew0ty"

func main() {
	client := messagebird.New(apiKey)
	// Request the balance information, returned as a Balance object.
	balance, err := client.Balance()
	if err != nil {
		// messagebird.ErrResponse means custom JSON errors.
		if err == messagebird.ErrResponse {
			for _, mbError := range balance.Errors {
				fmt.Printf("Error: %#v\n", mbError)
			}
		}

		return
	}

	fmt.Println("  payment :", balance.Payment)
	fmt.Println("  type    :", balance.Type)
	fmt.Println("  amount  :", balance.Amount)
}
