package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	con = flag.Int("con", 5, "number of concurrent clients repeating request")
	hld = flag.Duration("hold-for", time.Second, "hold for time")
	out = flag.String("out", "", "output file name")
	url = flag.String("url", "", "")
)

var (
	rch = make(chan record)
	w   = os.Stdout
)

func main() {
	flag.Parse()
	req, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		log.Fatal(err)
	}
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w = f
	}
	for i := 0; i < *con; i++ {
		go func() {
			for {
				do(req)
			}
		}()
	}
	done := time.After(*hld)
	for {
		select {
		case r := <-rch:
			fmt.Fprintln(w, r.sc, r.dur)
		case <-done:
			return
		}
	}
}

func do(req *http.Request) {
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
