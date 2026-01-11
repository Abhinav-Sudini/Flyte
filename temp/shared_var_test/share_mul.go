package main

import (
	"fmt"
	"test/share"
)

func Hi(){
	fmt.Println(share.I)
	share.Say()
	share.I = 20
	fmt.Println(share.I)
}
