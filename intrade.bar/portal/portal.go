package portal

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"intrade.bar/models"
	"intrade.bar/portal/request"
	"intrade.bar/trading_session"
	"log"
)

type Portal struct {
	ctx context.Context
	allocCancel context.CancelFunc
	ctxCancel context.CancelFunc

	username string
	password string

	TS tradingScreen
}



func (p *Portal) Initialize(username, password string){

	p.username = username
	p.password = password
	p.TS = tradingScreen{
		p: p,
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
	)

	var allocCtx context.Context

	allocCtx, p.allocCancel = chromedp.NewExecAllocator(context.Background(), opts...)
	p.ctx, p.ctxCancel = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
}

func (p Portal) OpenTradingScreen(){
	listenForNetworkEvent(p.ctx)
	chromedp.Run(p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {

		//c := chromedp.FromContext(ctx)

		chromedp.Navigate("https://intrade.bar").Do(ctx)
		chromedp.Click(`body > header > div > a.main-header__btn.btn.btn--black.p-open`).Do(ctx)
		chromedp.WaitVisible(`#log-in > form > p.pop-up-log-in__title`).Do(ctx)
		chromedp.SendKeys(`#input-log-in1`, p.username, chromedp.ByID).Do(ctx)
		chromedp.SendKeys(`#input-log-in2`, p.password, chromedp.ByID).Do(ctx)
		chromedp.Click(`#log-in > form > button`, chromedp.ByID).Do(ctx)
		chromedp.WaitVisible(`body > div.profile.profile--black > div.menu-profile > div > a > img`).Do(ctx)

		network.Enable().Do(ctx)
		chromedp.Navigate("https://intrade.bar/").Do(ctx)

		//chromedp.WaitVisible(`body > div.js-rootresizer__contents > div.layout__area--center`).Do(ctx)
		//time.Sleep(2 * time.Second)

		return nil
	}))
}

func (p Portal) CheckTimeAndInvestment(){
	chromedp.Run(p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var t, i string

		chromedp.Value(timeSelectorId, &t, chromedp.ByID).Do(ctx)
		chromedp.Value(investmentSelectorId, &i, chromedp.ByID).Do(ctx)

		fmt.Printf("investment: %s | time: %s\n", i, t)

		return nil
	}))
}

func listenForNetworkEvent(ctx context.Context) {
	go chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {

		case *network.EventRequestWillBeSent:
			if ev.Request.Method == "POST" {
				if yes, rt := request.GENERIC.IsAnyOfTargetType(ev.Request.URL); yes {
					switch rt {
					case request.NEW_TRADE:
						trading_session.TS.HandleNewOrderRequest(ev.Request.PostData, ev.RequestID.String(), nil)
						for _, e := range ev.Request.PostDataEntries {
							fmt.Println(e.Bytes)
						}
						trading_session.TS.AddOrder(models.Order{
							RequestId:  ev.RequestID.String(),
						})

					case request.PROFIT_PERCENT:
						trading_session.TS.ProfitCoefficient.RequestId = ev.RequestID.String()
					}



					fmt.Println("> RequestId: ", ev.RequestID)
					fmt.Println("DocumentUrl: ", ev.DocumentURL)
					fmt.Println("Url: ", ev.Request.URL)
				}
			}

		case *network.EventLoadingFinished:
			order := trading_session.TS.GetOrder(ev.RequestID.String())

			if order != nil {
				//response, err := network.GetResponseBody(ev.RequestID).Do(ctx)
				//fmt.Println("> RequestId: ", ev.RequestID)
				//fmt.Println("Response:", response, err)

				go func() {
					c := chromedp.FromContext(ctx)
					rbp := network.GetResponseBody(ev.RequestID)
					body, err := rbp.Do(cdp.WithExecutor(ctx, c.Target))
					if err != nil {
						fmt.Println(err)
					}
					if err == nil {
						fmt.Printf("%s\n", body)
						trading_session.TS.HandleNewOrderResponse(string(body), ev.RequestID.String())
					}
				}()
			}
		//case *network.EventLoadingFinished:
		//	res := ev
		//	var data []byte
		//	var e error
		//	if reqId1 == res.RequestID {
		//		data, e = network.GetResponseBody(reqId1).Do(ctx)
		//	} else if reqId2 == res.RequestID {
		//		data, e = network.GetResponseBody(reqId2).Do(ctx)
		//	}
		//	if e != nil {
		//		panic(e)
		//	}
		//	if len(data) > 0 {
		//		fmt.Printf("=========data: %+v\n", string(data))
		//	}

		case *network.EventResponseReceived:
			resp := ev.Response
			if len(resp.Headers) != 0 {
				//log.Printf("response url: %s", resp.URL)
				//log.Printf("response payload: %s", ev.Response.)
			}
		//case *network.EventWebSocketFrameReceived:
		//	payload := ev.Response.PayloadData
		//	fmt.Println(payload)
		}
		// other needed network Event
	})
}