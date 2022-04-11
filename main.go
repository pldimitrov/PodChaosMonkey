package main

import (
	"context"
	"podchaosmonkey/podchaosmonkey"
	"podchaosmonkey/util"
)

func main() {

	annotations, namespace, intervalSeconds, gracePeriodSeconds, err := util.ParseArgs()

	if err != nil {
		panic(err.Error())
	}

	podChaosMonkey := podchaosmonkey.NewPodChaosMonkey(
		podchaosmonkey.InitClientSet(),
		annotations,
		namespace,
		intervalSeconds,
		gracePeriodSeconds,
	)

	podChaosMonkey.Run(context.TODO())
}
