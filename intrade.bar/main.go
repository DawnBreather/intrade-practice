package main

import (
	"fmt"
	"intrade.bar/portal"
	"time"
)

var (
	username = "dawnbreather@gmail.com"
	password = "43Ete3DO77cMV"
)


func main(){
	fmt.Println("Hello")

	p := portal.Portal{}

	p.Initialize(username, password)
	p.OpenTradingScreen()

	p.TS.SetInstrument("1")
	p.TS.SetInvestment("2")
	p.TS.SetTime("1")

	p.CheckTimeAndInvestment()

	for {
		time.Sleep(5 * time.Second)
	}

	//chromedp.lis
}



