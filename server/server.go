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
	// pass every other connection resource
}

// GetQuotes return quotes
// no need for pointers
func (q *QuoteServer) GetQuotes(req *proto.QuoteRequest, stream proto.Quote_GetQuotesServer) error {
	c := make(chan Quote)
	go spit(stream, c)
	if req.Limit == 0 {
		req.Limit = 1
	}
	var wg sync.WaitGroup
	for i := 1; i <= int(req.Limit); i++ {
		i := i
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			res, err := getBasicQuotes(i)
			if err != nil {
				fmt.Printf("Err: %v", err)
				return
			}
			for _, v := range res {
				c <- v
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func spit(stream proto.Quote_GetQuotesServer, c chan Quote) error {
	for msg := range c {
		res := &proto.QuoteResponse{
			ID:       msg.ID,
			Body:     msg.Body,
			Author:   msg.Author,
			AuthorID: msg.AuthorID,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil
}

// GetFilteredQuotes return filtered version of quotes
func (q *QuoteServer) GetFilteredQuotes(req *proto.QuoteRequest, stream proto.Quote_GetFilteredQuotesServer) error {
	c := make(chan Quote)
	go spit(stream, c)
	if req.Limit == 0 {
		req.Limit = 1
	}

	res, err := getFilteredPaths(req.Author)
	if err != nil {
		fmt.Printf("Err: %v", err)
		return err
	}
	for _, v := range res {
		c <- v
	}

	return nil
}

func handleRequest(endpoint string) (result []Quote, err error) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
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

func getFilteredPaths(author string) (result []Quote, err error) {
	path := "https://stoicquotesapi.com/v1/api/quotes/" + author
	return handleRequest(path)
}

func getBasicQuotes(page int) (result []Quote, err error) {
	path := fmt.Sprintf("https://stoicquotesapi.com/v1/api/quotes?page=%d", page)
	return handleRequest(path)
}
