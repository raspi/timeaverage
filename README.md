# timeaverage

## Install
    go get -u github.com/raspi/timeaverage
    
## Example

```go
package main

func exampleSampler() (float64, error) {
	return 1, nil
}

func main() {
	avg := New(time.Second*10, time.Millisecond*500, 0.0, exampleSampler)
	avg.Start()

	for {
		v := samplah.Average()
		log.Printf(`%f`, v)
		time.Sleep(time.Second * 1)
	}

}
```
