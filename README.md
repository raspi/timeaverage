# timeaverage

Call a function which returns a float every X duration and keep history of N values and return average value of that window.  

## Install
    go get -u github.com/raspi/timeaverage
    
## Example

Take measurement every 500 ms and keep history of 10 seconds.

```go
package main

import (
	"github.com/raspi/timeaverage"
    "log"
	"time"
)

func exampleSampler() (float64, error) {
	return 1, nil
}

func main() {
	avg := timeaverage.New(time.Second*10, time.Millisecond*500, 0.0, exampleSampler)
	avg.Start()

	for {
		v := avg.Average()
		log.Printf(`%f`, v)
		time.Sleep(time.Second * 1)
	}

}
```

See [_examples](_examples/) directory for more examples.