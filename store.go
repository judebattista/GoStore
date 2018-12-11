//Coke and Pepsi

package main

import (
	"math/rand"
	"time"
)

type Bill struct {
	cokeDelivered  int
	pepsiDelivered int
}

type Brand int

const (
	Coke  Brand = 0
	Pepsi Brand = 1
)

type Shelf struct {
	Coke         int
	Pepsi        int
	CokeChannel  chan int
	PepsiChannel chan int
}

var cokeDelivery chan int
var pepsiDelivery chan int
var bill Bill
var shelf Shelf
var closed bool

//Go routine or for loop to spawn customers
//Each customer is a go routine
//Customers take soda in sets of six off the shelf and go to checkout
//Arrive every 100 ms, stay until they get what they want or the store closes
func spawnCustomers() {
	for closed == false {

	}

}

//Checkout is a go routine
//At checkout, record the amount of money

//Go routine
//Coke delivery - send to stocker
//Record delivery, add to bill
//Can deliver 24 cans every 500 ms
func deliverCoke() {
	for j := 0; j < 20; j++ {
		for i := 0; i < 24; i++ {
			cokeDelivery <- 1
		}
		bill.cokeDelivered += 1
		time.Sleep(500 * time.Millisecond)
	}
	close(cokeDelivery)
}

//Go routine
//Pepsi delivery - send to stocker
//Record delivery, add to bill
//Can deliver 24 cans every 500 ms
func deliverPepsi() {
	for j := 0; j < 20; j++ {
		for i := 0; i < 24; i++ {
			pepsiDelivery <- 1
		}
		bill.pepsiDelivered += 1
		time.Sleep(500 * time.Millisecond)
	}
	close(pepsiDelivery)
}

//Customers take sodas off the shelf in sets of 6, 12, 18, or 24
func customer() {
	targetNumber := (rand.Intn(4) * 6) + 6
	targetBrand := Brand(rand.Intn(2))
	if targetBrand == Coke {
		for sodas := 0; sodas < targetNumber; sodas++ {
			//If we've closed
			if shelf.CokeChannel == nil {
				return
			}
			//Take a coke off the shelf
			<-shelf.CokeChannel
			//Decrement the coke counter
			shelf.Coke--
		}
	} else {
		for sodas := 0; sodas < targetNumber; sodas++ {
			//If we've closed
			if shelf.PepsiChannel == nil {
				return
			}
			//take a Pepsi off the shelf
			<-shelf.PepsiChannel
			//Decrement the Pepsi counter
			shelf.Pepsi--
		}
	}

}

//Go routine
//Stocker - adds soda to shelf once per 10ms
//What if we make the shelf two channels?
func stocker() {
	//https://stackoverflow.com/questions/13666253/breaking-out-of-a-select-statement-when-all-channels-are-closed
	for {
		select {
		case _, ok := <-cokeDelivery:
			if ok {
				//Put the Coke on the shelf
				shelf.CokeChannel <- 1
				//Increment the Coke counter
				shelf.Coke++
			} else {
				cokeDelivery = nil
			}
		case _, ok := <-pepsiDelivery:
			if ok {
				//Put the Pepsi on the shelf
				shelf.PepsiChannel <- 1
				//Increase the Pepsi counter
				shelf.Pepsi++
			} else {
				pepsiDelivery = nil
			}
		}
		if pepsiDelivery == nil && cokeDelivery == nil {
			close(shelf.CokeChannel)
			close(shelf.PepsiChannel)
			return
		}
	}
	time.Sleep(10 * time.Millisecond)
}

func main() {
	var kind Brand
	shelf := Shelf{0, 0, make(chan int), make(chan int)}
	bill := Bill{0, 0}
	cokeDelivery := make(chan int)
	pepsiDelivery := make(chan int)
	kind = Coke
	time.Sleep(100 * time.Millisecond) //Sleep for 100 ms

}
