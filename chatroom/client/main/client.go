package main

import (
	"chatroom/client/processors"
)

func main() {
	processors.GetTheMainProcessor().Process(true)
}
