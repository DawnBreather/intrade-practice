package portal

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type tradingScreen struct {
	p *Portal
}

func (t tradingScreen) Call() error {

	return chromedp.Run(t.p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return chromedp.Click(callSelector).Do(ctx)
	}))
}

func (t tradingScreen) Put() error {
	return chromedp.Run(t.p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		network.Enable().Do(ctx)
		return chromedp.Click(putSelector).Do(ctx)
	}))
}

func (t tradingScreen) SetInvestment(value string) error{
	return chromedp.Run(t.p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return chromedp.SetValue(investmentSelectorId, value, chromedp.ByID).Do(ctx)
	}))
}

func (t tradingScreen) SetTime(value string) error{
	return chromedp.Run(t.p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return chromedp.SetValue(timeSelectorId, value, chromedp.ByID).Do(ctx)
	}))
}

func (t tradingScreen) SetInstrument(elementNumber string) error{
	return chromedp.Run(t.p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return chromedp.Click(instrumentSelector + elementNumber).Do(ctx)
	}))
}