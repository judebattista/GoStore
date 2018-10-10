//Coke and Pepsi

package main

import (
	"time"
)

type Brand int

const (
	Coke  Brand = 0
	Pepsi Brand = 1
)

type Shelf struct {
	Coke  int
	Pepsi int
}

//Go routine or four loop to spawn customers
//Each customer is a go routine
//Customers take soda in sets of six off the shelf and go to checkout
//Arrive every 100 ms, stay until they get what they want or the store closes

//Checkout is a go routine
//At checkout, record the amount of money

//Go routine
//Coke delivery - send to stocker
//Record delivery, add to bill
//Can deliver 24 cans every 500 ms

//Go routine
//Pepsi delivery - send to stocker
//Record delivery, add to bill
//Can deliver 24 cans every 500 ms

//Go routine
//Stocker - adds soda to shelf once per 10ms

func main() {
	var kind Brand

	kind = Coke
	time.Sleep(100, time.Millisecond) //Sleep for 100 ms

}
