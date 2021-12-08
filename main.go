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
	bpanda.SendToRedpanda(filePath)

	smongo.FetchFromRedpanda()
	// smongo.FetchFromNats()

	// sscylla.FetchFromNats()

	time.Sleep(time.Minute)
}
