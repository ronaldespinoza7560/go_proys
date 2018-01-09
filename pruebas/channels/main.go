package main                                                                                                                                                           

import (
    "fmt"
    "strconv"
	"sync"
)
var wg sync.WaitGroup
func makeCakeAndSend(cs chan string) {
    for i := 1; i<=3; i++ {
        cakeName := "Strawberry Cake " + strconv.Itoa(i)
        fmt.Println("Making a cake and sending ...", cakeName)
        cs <- cakeName //send a strawberry cake
	}   
	wg.Done()
	close(cs)
}

func receiveCakeAndPack(cs chan string) {
    for i := 1; i<=3; i++ {
        s := <-cs //get whatever cake is on the channel
        fmt.Println("Packing received cake: ", s)
	}   
	wg.Done()
}

func main() {
	cs := make(chan string)
	wg.Add(1)
	go makeCakeAndSend(cs)
	wg.Add(1)
    go receiveCakeAndPack(cs)

    //sleep for a while so that the program doesn’t exit immediately
    wg.Wait()
}