// Implement a part of a web-crowler program that concurrently sends HTTP GET requests to different endpoints and returns the results and errors into the main function.
// Apply a limitation on a number of simultaneously running HTTP requests.
// The order of results doesn't matter.
// The goal is to build a program that is capable of efficiently managing the execution of external requests in parallel, leveraging Go's concurrency features.


package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

// concurrentRequests function takes a slice of URLs,
// fetches data from each URL concurrently, and prints the results.
// The order of results is not important.

func concurrentRequests(urls []string, maxGR int) []Result {
	if maxGR <= 0 {
		maxGR = 10
	}

	var wg sync.WaitGroup
	ch := make(chan Result, len(urls))
	sem := make(chan struct{}, maxGR) 
	
	
	for _, u := range urls {
		wg.Add(1)
		sem <- struct{}{}
		
		go func(u string, wg *sync.WaitGroup){
			defer wg.Done()
			defer func(){<-sem}()
			
			resp, err := http.Get(u)
			if err != nil {
				ch <- Result{
					URL: u,
					Err: err,
				}
				return
			}
			
			defer resp.Body.Close()
			
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				ch <- Result{
					URL: u,
					Err: err,
				}
				return
			}

			ch <- Result{
				URL: u,
				Body: string(b),
				Err: nil,
			}
 			
		}(u, &wg)
	}

	wg.Wait()
	close(ch)
	results := make([]Result, 0, len(urls))
	for result := range ch {
		results = append(results, result)
		fmt.Println(result)
	}
	
	return results
 
}

func main() {
    
	urls := []string{
		"https://example.com",
		"https://httpbin.org/get",
		"https://jsonplaceholder.typicode.com/posts/1",
	}

    results := concurrentRequests(urls, 10)

	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("%s -> error: %v\n", r.URL, r.Err)
			continue
		}
		fmt.Printf("%s -> %d bytes\n", r.URL, len(r.Body))
	}
}


type Result struct {
	URL  string
	Body string
	Err  error
}
