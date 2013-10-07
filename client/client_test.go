package glock

import (
	"fmt"
	"testing"
	"time"
)

func TestLockUnlock(t *testing.T) {
	client1, err := NewClient("localhost:45625", 10)
	if err != nil {
		t.Error("Unexpected new client error: ", err)
	}

	fmt.Println("1 getting lock")
	id1, err := client1.Lock("x", 10*time.Second)
	if err != nil {
		t.Error("Unexpected lock error: ", err)
	}
	fmt.Println("1 got lock")

	go func() {
		fmt.Println("2 getting lock")
		id2, err := client1.Lock("x", 10*time.Second)
		if err != nil {
			t.Error("Unexpected lock error: ", err)
		}
		fmt.Println("2 got lock")

		time.Sleep(1 * time.Second)
		fmt.Println("2 releasing lock")
		err = client1.Unlock("x", id2)
		if err != nil {
			t.Error("Unexpected Unlock error: ", err)
		}
		fmt.Println("2 released lock")
	}()

	fmt.Println("sleeping")
	time.Sleep(2 * time.Second)
	fmt.Println("finished sleeping")

	fmt.Println("1 releasing lock")
	err = client1.Unlock("x", id1)
	if err != nil {
		t.Error("Unexpected Unlock error: ", err)
	}

	fmt.Println("1 released lock")

	time.Sleep(5 * time.Second)
}

func TestConnectionDrop(t *testing.T) {
	client1, err := NewClient("localhost:45625", 10)
	if err != nil {
		t.Error("Unexpected new client error: ", err)
	}

	fmt.Println("closing connection")
	client1.testClose()
	fmt.Println("closed connection")

	fmt.Println("1 getting lock")
	id1, err := client1.Lock("x", 1*time.Second)
	if err != nil {
		t.Error("Unexpected lock error: ", err)
	}
	fmt.Println("1 got lock")

	fmt.Println("1 releasing lock")
	err = client1.Unlock("x", id1)
	if err != nil {
		t.Error("Unexpected Unlock error: ", err)
	}
	fmt.Println("1 released lock")

	client1.testClose()

}

// This is used to simulate dropped out or bad connections in the connection pool
func (c *Client) testClose() {
	size := len(c.connectionPool)
	for x := 0; x < size; x++ {
		connection := <-c.connectionPool
		connection.Close()
		c.connectionPool <- connection
	}
}
