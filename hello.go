package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	con = flag.Int("con", 5, "number of concurrent clients repeating request")
	hld = flag.Duration("hold-for", time.Second, "hold for time")
	url = flag.String("url", "", "")
)

func main() {
	flag.Parse()
	req, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		log.Fatal(err)
	}
	rch := make(chan record)
	for i := 0; i < *con; i++ {
		go func() {
			for {
				do(req, rch)
			}
		}()
	}
	done := time.After(*hld)
	for {
		select {
		case r := <-rch:
			fmt.Println(r.sc, r.dur)
		case <-done:
			return
		}
	}
}

func do(req *http.Request, rch chan<- record) {
	start := time.Now()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	rch <- record{
		dur: time.Since(start),
		sc:  res.StatusCode,
	}
}

type record struct {
	dur time.Duration
	sc  int
}
