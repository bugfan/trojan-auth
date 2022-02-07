package main

import (
	"github.com/bugfan/trojan-auth/env"
	"github.com/bugfan/trojan-auth/srv"
)

func main() {
	// run it
	srv.NewAPIBackend().Run(env.Get("api_addr"))
}
