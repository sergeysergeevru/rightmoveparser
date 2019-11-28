package types

import "strings"

type RightMoveResponse struct {
	Properties []Property `json:"properties"`
}

type Property struct {
	Id uint32 `json:"id"`
	Bedrooms int `json:"bedrooms"`
	Price Price `json:"price"`
	PropertySubType string `json:"propertySubType"`
}

func (p *Property) IsShare() bool {
	return strings.Contains(p.PropertySubType, "Share")
}

type Price struct {
	Amount int `json:"amount"`
	Frequency string `json:"frequency"`
}

func (p *Price) GetMonthPrice() int {
	switch p.Frequency {
	case "monthly":
		return p.Amount
	case "weekly":
		return p.Amount*52/12
	default:
		panic("unexpected value")
	}
}