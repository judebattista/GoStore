//Coke and Pepsi
//Currently the main bottleneck is the Stocker routine, as demonstrated by the time gap between the deliveries finishing and the checkout line closing
//We can work around this by adding more Stocker routines, or having the Stocker stock cases instead of cans

//In order to move the bottleneck to the delivery person, we can improve the Stockers performance and have the delivery people deliver individual cans instead of cases.

//To move the bottleneck to the the check-out clerk, we can just have the clerk. Slow. Way. Down. If they only check out one customer a minute, we will have
//a rapidly filling queue of unhappy customers.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const cokePurchasePrice float32 = .25
const pepsiPurchasePrice float32 = .20
const cokeSellPrice float32 = .55
const pepsiSellPrice float32 = .50

var waitGroup sync.WaitGroup

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

type customer struct {
	brand Brand
	sodas int
}

type ledger struct {
	cokeSold       int
	pepsiSold      int
	cokeDelivered  int
	pepsiDelivered int
	income         float32
	outlay         float32
}

var shelf Shelf
var closed bool
var books ledger

//Go routine
//Coke delivery - send to stocker
//Record delivery, add to bill
//Can deliver 24 cans every 500 ms
func deliverCoke(cokeDelivery chan int) {
	waitGroup.Add(1)
	for j := 0; j < 20; j++ {
		for i := 0; i < 24; i++ {
			cokeDelivery <- 1
			books.cokeDelivered++
		}
		fmt.Printf("	Case #%v of Coke delivered!\n", j)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Done with Coke deliveries!")
	close(cokeDelivery)
	waitGroup.Done()
}

//Go routine
//Pepsi delivery - send to stocker
//Record delivery, add to bill
//Can deliver 24 cans every 500 ms
func deliverPepsi(pepsiDelivery chan int) {
	waitGroup.Add(1)
	for j := 0; j < 20; j++ {
		for i := 0; i < 24; i++ {
			pepsiDelivery <- 1
			books.pepsiDelivered++
		}
		fmt.Printf("	Case #%v of Pepsi delivered!\n", j)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Done with Pepsi deliveries!")
	close(pepsiDelivery)
	waitGroup.Done()
}

func reconcileBill(books ledger) (balance float32) {
	balance = float32(books.cokeSold) * cokeSellPrice
	balance += float32(books.pepsiSold) * pepsiSellPrice
	balance -= float32(books.cokeDelivered) * cokePurchasePrice
	balance -= float32(books.pepsiDelivered) * pepsiPurchasePrice
	return
}

//Go routine
//Stocker - adds soda to shelf once per 10ms
//What if we make the shelf two channels?
func stocker(cokeDelivery chan int, pepsiDelivery chan int) {
	//https://stackoverflow.com/questions/13666253/breaking-out-of-a-select-statement-when-all-channels-are-closed
	waitGroup.Add(1)
	for {
		select {
		case _, ok := <-cokeDelivery:
			if ok {
				fmt.Println("			Stocking coke.")
				//Put the Coke on the shelf
				shelf.CokeChannel <- 1
				//Increment the Coke counter
				shelf.Coke++
			} else {
				cokeDelivery = nil
			}
		case _, ok := <-pepsiDelivery:
			if ok {
				fmt.Println("			Stocking pepsi.")
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
			closed = true
			fmt.Println("Finished stocking all sodas!")
			waitGroup.Done()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

}

//Go routine to spawn customers
//Each customer is a go routine
//Customers take soda in sets of six off the shelf and go to checkout
//Arrive every 100 ms, stay until they get what they want or the store closes
func spawnCustomers(checkout chan customer) {
	waitGroup.Add(1)
	for closed == false {
		fmt.Println("			Spawning Customer.")
		var cust customer
		go shop(cust, checkout)
		time.Sleep(75 * time.Millisecond)
	}
	//When the store closes, close the checkout channel
	fmt.Println("Closing the checkout queue")
	close(checkout)
	waitGroup.Done()
}

//Customers take sodas off the shelf in sets of 6, 12, 18, or 24
func shop(cust customer, checkout chan customer) {
	waitGroup.Add(1)
	targetNumber := (rand.Intn(4) * 6) + 6
	targetBrand := Brand(rand.Intn(2))
	sodas := 0
	if targetBrand == Coke {
		for ; sodas < targetNumber; sodas++ {
			//If we've closed
			if shelf.CokeChannel == nil {
				return
			}
			//Otherwise take a coke off the shelf
			<-shelf.CokeChannel
			//Decrement the coke counter
			shelf.Coke--
		}
		fmt.Printf("	Put %v Cokes in cart.\n", sodas)
	} else {
		for ; sodas < targetNumber; sodas++ {
			//If we've closed
			if shelf.PepsiChannel == nil {
				return
			}
			//Otherwise take a Pepsi off the shelf
			<-shelf.PepsiChannel
			//Decrement the Pepsi counter
			shelf.Pepsi--
		}
		fmt.Printf("	Put %v Pepsis in cart.\n", sodas)
	}
	cust.brand = targetBrand
	cust.sodas = sodas
	//This condition currently terminates too early.
	//We wind up with a bunch of customers with full carts in the store and no way to pay
	//Need to find a way to process them, or process customers with higher priority
	if closed == false {
		checkout <- cust
	}
	waitGroup.Done()
	return
}

//Checkout is not a go routine
//At checkout, record the amount of money
func checkOut(cust customer) {
	fmt.Println("Checking out customer")
	var sodaBrand string
	if cust.brand == Coke {
		books.cokeSold += cust.sodas
	} else {
		books.pepsiSold += cust.sodas
	}
	if cust.brand == Coke {
		sodaBrand = "Coke"
	} else {
		sodaBrand = "Pepsi"
	}
	fmt.Printf("Just sold %d %vs\n", cust.sodas, sodaBrand)
}

//go routine representing a checkout clerk
func clerk(checkout chan customer) {
	waitGroup.Add(1)
	//While checkout is not closed and empty
	for cust := range checkout {
		fmt.Println("Customer in line")
		checkOut(cust)
	}
	fmt.Println("Quitting time.")
	waitGroup.Done()
}

func main() {
	//var kind Brand
	shelf = Shelf{0, 0, make(chan int), make(chan int)}
	cokeDelivery := make(chan int)
	pepsiDelivery := make(chan int)
	checkoutLine := make(chan customer)
	go clerk(checkoutLine)
	//time.Sleep(10 * time.Millisecond)
	go spawnCustomers(checkoutLine)
	time.Sleep(1000 * time.Millisecond)
	go stocker(cokeDelivery, pepsiDelivery)
	//time.Sleep(10 * time.Millisecond)
	go deliverCoke(cokeDelivery)
	go deliverPepsi(pepsiDelivery)

	//time.Sleep(1 * time.Second)
	//closed = <-quittinTime
	waitGroup.Wait()

	balance := reconcileBill(books)
	fmt.Println("********************************************************")
	fmt.Printf("The store stocked %d Cokes and %d Pepsis.\n", books.pepsiDelivered, books.pepsiDelivered)
	fmt.Printf("And sold %d Cokes and %d Pepsis.\n", books.cokeSold, books.pepsiSold)
	fmt.Printf("The end of day balance was $%v.\n", balance)
	fmt.Println("********************************************************")
}
