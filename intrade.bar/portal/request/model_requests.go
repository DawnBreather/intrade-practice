package request

import "strings"

type RequestType string
type requestTypeGeneric string

const (
	GENERIC        requestTypeGeneric = "generic"
	NEW_TRADE      RequestType        = "ajax5_new.php"
	TRADE_CHECK    RequestType        = "trade_check2.php"
	PROFIT_PERCENT RequestType        = "ajax_percent.php"
	NONE           RequestType        = ""
)

func (r requestTypeGeneric) IsAnyOfTargetType(url string) (yes bool, requestType RequestType){
	if strings.HasSuffix(url, NEW_TRADE.ToString()){
		return true, NEW_TRADE
	}

	if strings.HasSuffix(url, TRADE_CHECK.ToString()){
		return true, TRADE_CHECK
	}

	if strings.HasSuffix(url, PROFIT_PERCENT.ToString()){
		return true, PROFIT_PERCENT
	}

	return false, NONE
}

func (nt RequestType) ToString() string{
	return string(nt)
}
