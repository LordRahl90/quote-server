package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/LordRahl90/quote-server/proto"
)

type stoicResponse struct {
	Data []Quote `json:"data"`
}

// Quote quote object
type Quote struct {
	ID       uint32 `json:"id"`
	AuthorID uint32 `json:"author_id"`
	Body     string `json:"body"`
	Author   string `json:"author"`
}

// QuoteServer servers for quote management
type QuoteServer struct {
	proto.QuoteServer
	// pass every other connection
}

// GetQuotes return quotes
// no need for pointers
func (q *QuoteServer) GetQuotes(req *proto.QuoteRequest, stream proto.Quote_GetQuotesServer) error {
	// var c chan Quote
	c := make(chan Quote)
	go spitOut(stream, c)
	if req.Limit == 0 {
		req.Limit = 1
	}
	var wg sync.WaitGroup
	for i := 0; i < int(req.Limit); i++ {
		wg.Add(1)
		go func() {
			i := 1
			defer wg.Done()
			res, err := getQuotes(i)
			if err != nil {
				fmt.Printf("Err: %v", err)
			}
			for _, v := range res {
				fmt.Printf("\nV: %+v\n\n", v)
				c <- v
			}
		}()
	}
	fmt.Printf("Waiting\n")
	wg.Wait()
	return nil
}

func spitOut(stream proto.Quote_GetQuotesServer, c chan Quote) error {
	fmt.Printf("Spitting started\n")
	for msg := range c {
		fmt.Printf("\nReceived: %v\n\n", msg)
		res := &proto.QuoteResponse{
			ID:       msg.ID,
			Body:     msg.Body,
			Author:   msg.Author,
			AuthorID: msg.AuthorID,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		fmt.Printf("Message Sent\n")
	}
	return nil
}

// GetFilteredQuotes return filtered version of quotes
func (q *QuoteServer) GetFilteredQuotes(req *proto.QuoteRequest, stream proto.Quote_GetFilteredQuotesServer) error {
	return nil
}

func getQuotes(limit int) (result []Quote, err error) {
	fmt.Printf("Getting Quotes\n")
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	path := "https://stoicquotesapi.com/v1/api/quotes"

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return
	}
	response, err := client.Do(req)
	if err != nil {
		return
	}
	var res stoicResponse
	if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
		return
	}
	result = res.Data
	return
}
