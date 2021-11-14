package main

import (
	"flag"
	"fmt"
)

func main() {
	wordPtr := flag.String("word", value:"foo", usage:"a string")
	numPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	//call flag.Parse() to perform command-line parsing.
	flag.Parse()

	fmt.Println("word", *wordPtr)
	fmt.Println("numb", *numPtr)
	fmt.Println("bool", *boolPtr)

	flag.Usage()

}
