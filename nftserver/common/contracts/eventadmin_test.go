package contracts

import (
	"fmt"
	"testing"
)

func TestGetAdminList(t *testing.T) {
	addr, err := AdminList()
	if err != nil {
		fmt.Println("get admin error.")
	}
	fmt.Println(addr)
}