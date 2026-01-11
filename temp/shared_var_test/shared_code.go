package main

import (
	"test/share"
)

func main(){
	share.Say()
	share.I = 10

	share.Say()
	Hi()
}
