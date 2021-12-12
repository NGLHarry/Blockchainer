package main

import (
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

func main() {

	// 160
	// b66140b4bfd22da44399f352b07182864098123f
	// bit 160 20
	hasher := ripemd160.New()

	hasher.Write([]byte("http://liyuechun.org"))

	bytes := hasher.Sum(nil)

	fmt.Printf("%x\n", bytes)

}
