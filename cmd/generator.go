package main

import (
	"fmt"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	amount := os.Args[1]
	intAm, err := strconv.Atoi(amount)
	if err != nil {
		panic(err)
	}
	for i := 0; i < intAm; i++ {
		fmt.Println(primitive.NewObjectID())
	}
}
