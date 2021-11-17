package lknpd

import (
	"github.com/shopspring/decimal"
	"time"
)

type PaymentType string
type IncomeType string

const (
	Cash PaymentType = "CASH"

	Individual  IncomeType = "FROM_INDIVIDUAL"
	LegalEntity IncomeType = "FROM_LEGAL_ENTITY"

	TimeFormat string = "2006-01-02T15:04:05+00:00"
)

type Income struct {
	PaymentType                     PaymentType     `json:"paymentType,omitempty"`
	IgnoreMaxTotalIncomeRestriction bool            `json:"ignoreMaxTotalIncomeRestriction,omitempty"`
	Client                          Client          `json:"client"`
	RequestTime                     string          `json:"requestTime,omitempty"`
	OperationTime                   string          `json:"operationTime,omitempty"`
	Services                        []Service       `json:"services,omitempty"`
	TotalAmount                     decimal.Decimal `json:"totalAmount,omitempty"`
}

type Client struct {
	ContactPhone string     `json:"contactPhone,omitempty"`
	DisplayName  string     `json:"displayName,omitempty"`
	IncomeType   IncomeType `json:"incomeType,omitempty"`
	INN          string     `json:"inn,omitempty"`
}

type Service struct {
	Name     string          `json:"name,omitempty"`
	Amount   decimal.Decimal `json:"amount,omitempty"`
	Quantity int             `json:"quantity,omitempty"`
}

type Receipt struct {
	UUID     string
	InfoURL  string
	PrintURL string
}

func ParseTime(t time.Time) string {
	return t.UTC().Format(TimeFormat)
}

func (i *Income) SetClientType(t IncomeType) {
	i.Client.IncomeType = t
}

func (i *Income) SetClientName(s string) {
	i.Client.DisplayName = s
}

func (i *Income) SetClientINN(inn string) {
	i.Client.INN = inn
}

func (i *Income) SetOperationTime(t time.Time) {
	i.OperationTime = ParseTime(t)
}

func (i *Income) AddService(s Service) {
	i.Services = append(i.Services, s)
	i.TotalAmount.Add(s.Amount.Mul(decimal.NewFromInt32(int32(s.Quantity))))
}

func DefaultIncome() Income {
	t := ParseTime(time.Now())

	return Income{
		PaymentType:                     Cash,
		IgnoreMaxTotalIncomeRestriction: false,
		Client: Client{
			IncomeType: Individual,
		},
		RequestTime:   t,
		OperationTime: t,
		Services:      []Service{},
		TotalAmount:   decimal.New(0, 1),
	}
}
