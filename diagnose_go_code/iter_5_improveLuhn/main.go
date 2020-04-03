// Read the txt file
// Extract the ID (first column) and the JSON (the second column) that contain countries and continents.
// Display how many ID's are qualify as Luhn numbers
// Display the country with the highest number of hits
// Display the continent that has largest number of countries with hits over 1000

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/Danr17/dev-state_blog_code/tree/master/diagnose_go_code/iter_5_improveLuhn/luhn"
)

type result struct {
	validLuhn     int
	nr            int
	hitsCountry   string
	hits          int
	hitsContinent string
	hitsR         int
}

func main() {

	/*
		tracefile, err := os.OpenFile("m.trace", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer tracefile.Close()

		trace.Start(tracefile)
		defer trace.Stop()
	*/

	file, err := os.Open("../csv_files/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	results := result{}
	routines := runtime.NumCPU() * 2
	results.getStatistics(file, routines)

	fmt.Printf("There are %d, out of %d, valid Luhn numbers. \n", results.validLuhn, results.nr)
	fmt.Printf("%s has the biggest # of visitors, with %d of hits. \n", results.hitsCountry, results.hits)
	fmt.Printf("%s is the continent with most unique countries that accessed the site more than 1000 times. It has %d unique countries. \n", results.hitsContinent, results.hitsR)

}

func (r *result) getStatistics(stream io.Reader, routines int) {

	countries := map[string]int{}
	continents := map[string]string{}

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	const CacheSize = 64 * 1024

	linesPool := sync.Pool{
		New: func() interface{} {
			return new([CacheSize]string)
		},
	}

	lines := make(chan *[CacheSize]string, routines)

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for cache := range lines {

				validLuhn := 0
				continent := ""
				country := ""

				for _, text := range cache {

					split := strings.Split(text, "#")
					if len(split) < 2 {
						continue
					}
					number := strings.TrimSpace(split[0])
					description := strings.TrimSpace(split[1])

					dep := strings.SplitN(description, "\"", 13)
					continent = dep[3]
					country = dep[11]

					if luhn.Valid(number) {
						validLuhn++
					}

					mutex.Lock()
					countries[country]++

					if _, ok := continents[country]; !ok {
						continents[country] = continent
					}
					mutex.Unlock()
				}

				mutex.Lock()
				r.validLuhn += validLuhn
				mutex.Unlock()

				linesPool.Put(cache)
			}
		}()
	}

	scanner := bufio.NewScanner(stream)
	iter := 0
	pool := linesPool.Get().(*[CacheSize]string)

	for {
		valid := scanner.Scan()
		if iter == CacheSize || !valid {
			lines <- pool
			iter = 0
			pool = linesPool.Get().(*[CacheSize]string)

		}
		if !valid {
			break
		}
		r.nr++
		pool[iter] = scanner.Text()
		iter++
	}

	close(lines)
	wg.Wait()

	for k, v := range countries {
		if v > r.hits {
			r.hits = v
			r.hitsCountry = k
		}
	}

	regions := map[string]int{}
	for k, v := range continents {
		if countries[k] > 1000 {
			regions[v]++
		}
	}

	r.hitsR = 1
	for k, v := range regions {
		if v > r.hitsR {
			r.hitsContinent = k
			r.hitsR = v
		}
	}
}
