package main

import "log"

func main() {
	root := newRoot()
	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}
