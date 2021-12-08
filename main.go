package main

import (
	"renamosy/bpanda"
	"renamosy/smongo"
	"time"
)

func main() {
	time.Sleep(time.Second * 2)

	filePath := "./mockfiles"

	// bnats.SendToNats(filePath)
	// smongo.FetchFromNats()

	bpanda.SendToRedpanda(filePath)
	smongo.FetchFromRedpanda()

	// sscylla.FetchFromNats()
	// sscylla.FetchFromRedpanda()

	time.Sleep(time.Minute)
}
