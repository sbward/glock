package glock

import (
	"fmt"
	"testing"
	"time"
)

func TestLockUnlock(t *testing.T) {
	client1, err := NewClient("localhost:45625")
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
		client2, err := NewClient("localhost:45625")
		if err != nil {
			t.Error("Unexpected new client error: ", err)
		}
		fmt.Println("2 getting lock")
		id2, err := client2.Lock("x", 10*time.Second)
		if err != nil {
			t.Error("Unexpected lock error: ", err)
		}
		fmt.Println("2 got lock")

		fmt.Println("2 releasing lock")
		ok, err := client2.Unlock("x", id2)
		if err != nil {
			t.Error("Unexpected Unlock error: ", err)
		}
		if !ok {
			t.Error("Unexpected Unlock error: unable to Unlock")
		}
		fmt.Println("2 released lock")
		err = client2.Close()
		if err != nil {
			t.Error("Unexpected connection close error: ", err)
		}
	}()

	fmt.Println("sleeping")
	time.Sleep(5 * time.Second)
	fmt.Println("finished sleeping")

	fmt.Println("1 releasing lock")
	ok, err := client1.Unlock("x", id1)
	if err != nil {
		t.Error("Unexpected Unlock error: ", err)
	}
	if !ok {
		t.Error("Unexpected Unlock error: unable to Unlock")
	}

	fmt.Println("1 released lock")

	err = client1.Close()
	if err != nil {
		t.Error("Unexpected connection close error: ", err)
	}

	time.Sleep(5 * time.Second)
}
