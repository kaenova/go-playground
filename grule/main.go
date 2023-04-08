/*
rule ValidCoupons "Check if the coupons is valid" salience 1 {
	when
}

rule Year2022 "Check the default values" salience 10 {
    when
        MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
    then
        MF.WhatToSay = MF.GetWhatToSay("Hello Grule");
        Retract("CheckValues");
}

rule CheckValues2 "Check the default values 21" salience 10 {
	when
			MF.IntAttribute == 123
	then
			MF.WhatToSay = MF.GetWhatToSay("Hello Grule anjir");
			Retract("CheckValues2");
}
*/

package main

import (
	"fmt"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type Data struct {
	CouponUsed string

	Ammount float64

	PaymentMethod string

	UserCreatedYear  int
	IntAttribute     int64
	StringAttribute  string
	BooleanAttribute bool
	FloatAttribute   float64
	TimeAttribute    time.Time
	WhatToSay        string
}

type Response struct {
	ReducedPrice float64
	AddedPrice   float64
	ErrorMessage string
}

func (R *Response) AddValue(val float64) {
	R.AddedPrice = val
}

func (R *Response) ReducePrice(val float64) {
	R.ReducedPrice = val
}

func (R *Response) Error(msg string) {
	R.ErrorMessage = msg
}

func NewResponse() *Response {
	return &Response{}
}

func NewData() *Data {
	return &Data{
		UserCreatedYear: 2019,

		Ammount: 100000,

		PaymentMethod: "credit",

		IntAttribute:     123,
		StringAttribute:  "Some string value",
		BooleanAttribute: true,
		FloatAttribute:   1.234,
		TimeAttribute:    time.Now(),
	}
}

func main() {

	myFact := NewData()

	myResponse := NewResponse()

	dataCtx := ast.NewDataContext()

	err := dataCtx.Add("DATA", myFact)
	if err != nil {
		panic(err)
	}

	err = dataCtx.Add("RES", myResponse)
	if err != nil {
		panic(err)
	}

	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	drlsKupon1 := `
rule Year2018 "Coupon for 2018" {
	when
		DATA.UserCreatedYear == 2018
	then
		RES.ReducePrice(DATA.Ammount);
		Retract("Year2018");
}

rule Year2020 "Coupon for 2020" {
	when
		DATA.UserCreatedYear == 2020
	then
		RES.ReducePrice(0.0);
		Retract("Year2020");
}

rule CreditCard "Coupon with credit card" salience 1 {
	when
		DATA.PaymentMethod == 'credit'
	then
		RES.ReducePrice(DATA.Ammount * 0.4);
		Retract("CreditCard");
}

rule Year2019 "Coupon for 2019" salience 1 {
	when
		DATA.UserCreatedYear == 2019
	then
		RES.ReducePrice(DATA.Ammount / 2);
		Retract("Year2019");
}

rule CreditCard2 "Coupon with credit card" salience -1 {
	when
		DATA.PaymentMethod == 'credit' && DATA.UserCreatedYear == 2019
	then
		RES.ReducePrice(DATA.Ammount * 0.7);
		Retract("CreditCard2");
}
`

	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	bs := pkg.NewBytesResource([]byte(drlsKupon1))
	err = ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", bs)
	if err != nil {
		panic(err)
	}

	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

	gruleEngine := engine.NewGruleEngine()
	err = gruleEngine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	fmt.Println("Actual harga", myFact.Ammount)
	fmt.Println("Final harga", myFact.Ammount-myResponse.ReducedPrice+myResponse.AddedPrice)
}
