package main

import (
	"clientServer/client"
	"strconv"
	"sync"
)

func main()  {
	var wg sync.WaitGroup

	client1 := client.NewClient("127.0.0.1:1234")
	client2 := client.NewClient("127.0.0.1:1234")
	client3 := client.NewClient("127.0.0.1:1234")

	arr1 := makeRange(1, 10)

	wg.Add(6) //counter set to 6 as number of go func

	go func() {
		addItems(client1, arr1, 1, 3)
		wg.Done() //decrement counter each call
	}()
	go func() {
		addItems(client2, arr1, 4, 7)
		wg.Done()
	}()
	go func() {
		addItems(client3, arr1, 8, 9)
		wg.Done()
	}()
	go func() {
		client2.AddItem("KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK")
		wg.Done()
	}()
	go func() {
		popItems(client3, 3)
		wg.Done()
	}()
	go func() {
		popItems(client1, 4)
		wg.Done()
	}()

	wg.Wait() //block the main thread until the counter is 0
}

func popItems(currClient *client.Client, numPops int){
	for i := 0; i < numPops; i++ {
		currClient.PopItem()
	}
}

func addItems(currClient *client.Client, arr []int, from int, to int){
	for i := from; i < to; i++ {
		currClient.AddItem(strconv.Itoa(arr[i]))
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
