package main

type PlayerData struct {
	name  string
	pid   int
	ready bool
	glory int
}

func (pd *PlayerData) setGlory(glory int) {
	pd.glory = glory
}

func (pd *PlayerData) toggleReady() {
	pd.ready = !pd.ready
}
