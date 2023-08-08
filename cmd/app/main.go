package main

import "app/cmd/app/commands"

func main() {
	if err := commands.Root.Execute(); err != nil {
		panic(err)
	}
}
