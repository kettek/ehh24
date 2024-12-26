package main

import (
	. "github.com/kettek/gobl"
)

func main() {
	println("Hello, 世界")
	Task("build").
		Exec("go", "build", "./cmd/ehh24")

	Task("run").
		Exec("./ehh24")

	Task("watch").
		Watch("**/*.go", "**/*.png").
		Signaler(SigQuit).
		Run("build").
		Run("run")

	Go()
}
