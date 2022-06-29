package models

type Orders struct {
	ByRequestId map[string]*Order
	ByTradeId   map[string]*Order
}
type Order struct {
	OpenTime int64
	OpenPrice float64
	ClosePrice float64
	RequestId string
	TradeId string

	UserId     string
	UserHash   string
	Option     string
	Duration   int
	Investment int

	BaseInvestment                int
	ProfitCoefficient             float64
	NextInvestmentInCaseOfFailure int
	InvestmentLineId              string
	SequenceInInvestmentLine      int
}