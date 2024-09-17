package main

import "log"

func must(err error) {
	if err != nil {
		log.Panicf(err.Error())
	}
}

func mustNot(err error) {
	must(err)
}
