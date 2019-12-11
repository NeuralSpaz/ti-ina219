# TI-INA219
golang library for Texas Instruments INA219 High-Side, I2C Out Current/Power Monitor

Example
```go
package main

import (
	"fmt"
	"log"
	"time"

	ina219 "github.com/NeuralSpaz/ti-ina219"
)

func main() {
	// Create new sensor with an address of 0x41 and on i2c bus 1.
	Sensor1 := ina219.New(0x41, 1)
	for {
		select {
		case <-time.After(time.Second * 1):
			if err := ina219.Fetch(Sensor1); err != nil {
				log.Println(err)
			}
			fmt.Println(Sensor1)
		}
	}
}
```
