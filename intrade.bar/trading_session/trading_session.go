package trading_session

import (
	"intrade.bar/helpers"
	. "intrade.bar/models"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var TS = newTradingSession()

func newTradingSession() tradingSession{
	ts := tradingSession{
		InvestmentLines: map[string]*Order{},
		Orders: Orders{
			ByRequestId: map[string]*Order{},
			ByTradeId:   map[string]*Order{},
		},
	}

	return ts
}

type tradingSession struct {

	BaseInvestment int
	ProfitCoefficient struct{
		RequestId string
		Value float64
	}
	InvestmentLines map[string]*Order //Uuid of investment line => Last order
	Orders Orders
}

func (ts *tradingSession) HandleProfitCoefficient(percent string) {
	percentInt, _ := strconv.Atoi(percent)
	ts.ProfitCoefficient.Value = float64(1 + percentInt / 100)
}

func (ts *tradingSession) HandleNewOrderResponse(data string, requestId string){
	//[ invoice ]
	//data-id="[0-9]*?"
	//
	//[ open price ]
	//data-rate="[0-9]*\.[0-9]*?"
	//
	//[ time open ]
	//data-timeopen="[0-9]*?"

	invoiceRegex := regexp.MustCompile(`data-id="[0-9]*?"`)
	openPriceRegex := regexp.MustCompile(`data-rate="[0-9]*\.[0-9]*?"`)
	timeOpenRegex := regexp.MustCompile(`data-timeopen="[0-9]*?"`)

	invoiceRaw := invoiceRegex.FindString(data)
	invoiceRaw = strings.ReplaceAll(invoiceRaw, "data-id=", "")
	invoiceString := strings.Trim(invoiceRaw, `"`)

	openPriceRaw := openPriceRegex.FindString(data)
	openPriceRaw = strings.ReplaceAll(openPriceRaw, `data-rate=`, "")
	openPriceString := strings.Trim(openPriceRaw, `"`)

	timeOpenRaw := timeOpenRegex.FindString(data)
	timeOpenRaw = strings.ReplaceAll(timeOpenRaw, `data-timeopen=`, "")
	timeOpenString := strings.Trim(timeOpenRaw, `"`)

	order := ts.GetOrder(requestId)
	order.TradeId = invoiceString
	order.OpenPrice, _ = strconv.ParseFloat(openPriceString, 64)
	order.OpenTime, _ = strconv.ParseInt(timeOpenString, 10, 64)
	order.BaseInvestment = ts.BaseInvestment
	order.ProfitCoefficient = ts.ProfitCoefficient.Value

	ts.Orders.ByTradeId[order.TradeId] = order

}

func (ts *tradingSession) HandleNewOrderRequest(postData string, requestId string, previousOrder *Order){
	var o Order
	for _, entry := range strings.Split(postData, "&"){
		keyValue := strings.Split(entry, "=")
		key := keyValue[0]
		value := keyValue[1]

		switch key {
		case "user_id":
			o.UserId = value
		case "user_hash":
			o.UserHash = value
		case "option":
			o.Option = value
		case "investment":
			o.Investment, _ = strconv.Atoi(value)
		case "time":
			o.Duration, _ = strconv.Atoi(value)
		}
	}

	o.RequestId = requestId

	o.ProfitCoefficient = ts.ProfitCoefficient.Value

	if previousOrder == nil {
		o.InvestmentLineId = strconv.FormatInt(time.Now().Unix(), 10)
		o.SequenceInInvestmentLine = 0
	} else {
		o.InvestmentLineId = previousOrder.InvestmentLineId
		o.SequenceInInvestmentLine = previousOrder.SequenceInInvestmentLine + 1
		o.NextInvestmentInCaseOfFailure = ts.calculateNextInvestmentInCaseOfFailure(previousOrder.Investment)
	}
	ts.InvestmentLines[o.InvestmentLineId] = &o

	if _, ok := ts.Orders.ByRequestId[o.RequestId]; !ok {
		ts.Orders.ByRequestId[o.RequestId] = &o
	}
}

func (ts *tradingSession) AddOrder(o Order){
	if o.RequestId != "" {
		if _, ok := ts.Orders.ByRequestId[o.RequestId]; !ok {
			ts.Orders.ByRequestId[o.RequestId] = &o
		}
	}

	if o.TradeId != "" {
		if _, ok := ts.Orders.ByTradeId[o.TradeId]; !ok {
			ts.Orders.ByTradeId[o.TradeId] = &o
		}
	}
}

func (ts *tradingSession) GetOrder(id string) *Order{
	if order, ok := ts.Orders.ByTradeId[id]; ok {
		return order
	}
	if order, ok := ts.Orders.ByRequestId[id]; ok {
		return order
	}

	return nil
}

func (ts tradingSession) calculateNextInvestmentInCaseOfFailure(previousInvestment int) int{
	multiplier := ts.ProfitCoefficient.Value / (ts.ProfitCoefficient.Value - 1)
	return helpers.RoundUp(float64(previousInvestment) * multiplier)
}