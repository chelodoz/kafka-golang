package payment

import (
	"fmt"
	"strings"

	fake "github.com/brianvoe/gofakeit/v6"
)

type PaymentMessage struct {
	ClientReferenceInformation ClientReferenceInformation `json:"clientReferenceInformation"`
	PaymentInformation         PaymentInformation         `json:"paymentInformation"`
	OrderInformation           OrderInformation           `json:"orderInformation"`
}

type ClientReferenceInformation struct {
	Code string `json:"code"`
}
type Card struct {
	Number          string `json:"number"`
	ExpirationMonth string `json:"expirationMonth"`
	ExpirationYear  string `json:"expirationYear"`
}
type PaymentInformation struct {
	Card Card `json:"card"`
}
type AmountDetails struct {
	TotalAmount string `json:"totalAmount"`
	Currency    string `json:"currency"`
}
type BillTo struct {
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Address1           string `json:"address1"`
	Locality           string `json:"locality"`
	AdministrativeArea string `json:"administrativeArea"`
	PostalCode         string `json:"postalCode"`
	Country            string `json:"country"`
	Email              string `json:"email"`
	PhoneNumber        string `json:"phoneNumber"`
}
type OrderInformation struct {
	AmountDetails AmountDetails `json:"amountDetails"`
	BillTo        BillTo        `json:"billTo"`
}

func NewMessage() *PaymentMessage {
	return &PaymentMessage{
		ClientReferenceInformation: ClientReferenceInformation{
			Code: fake.UUID(),
		},
		PaymentInformation: PaymentInformation{
			Card: Card{
				Number:          fake.CreditCard().Number,
				ExpirationMonth: strings.Split(fake.CreditCard().Exp, "/")[0],
				ExpirationYear:  strings.Split(fake.CreditCard().Exp, "/")[1],
			},
		},
		OrderInformation: OrderInformation{
			AmountDetails: AmountDetails{
				TotalAmount: fmt.Sprintf("%f", fake.Price(100, 2000)),
				Currency:    fake.CurrencyShort(),
			},
			BillTo: BillTo{
				FirstName:          fake.Person().FirstName,
				LastName:           fake.Person().LastName,
				Address1:           fake.Address().Address,
				Locality:           fake.Address().City,
				AdministrativeArea: fake.Address().State,
				PostalCode:         fake.Address().Zip,
				Country:            fake.Address().Country,
				Email:              fake.Person().Contact.Email,
				PhoneNumber:        fake.Person().Contact.Phone,
			},
		},
	}
}
