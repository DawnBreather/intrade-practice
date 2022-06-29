package redis_client

import (
	"context"
	"fmt"
	"intrade.bar/models"
)

func (c c) saveOrder(o models.Order){ //(option string, investmentLineId string, duration int, investment int, baseInvestment int, profitCoefficient float64, openTimeUnix int64, openPrice, closePrice float64, sequenceNumber int, tradeId string, userId string, userHash string){
	// i.e. /orders/eurusd/9012832109380/1/16932047102
	key := fmt.Sprintf(`/orders/%s/%s/%d/%d`, o.Option, o.InvestmentLineId, o.SequenceInInvestmentLine, o.OpenTime)
	c.Set(context.Background(), key, o, 0)
}

func (c c) saveQuotes(option string, timeUnix int64, price float64){
	// i.e. /quotes/eurusd/16120839721
	key := fmt.Sprintf(`/%s/%s/%d`, option, timeUnix)
	c.Set(context.Background(), key, price, 0)
}