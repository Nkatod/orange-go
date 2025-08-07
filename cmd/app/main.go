package main

import (
	"orange-go/internal/storage/memory"
)

func main() {
	memory.Put("key", "value")
	println(memory.Store["key"])
}
