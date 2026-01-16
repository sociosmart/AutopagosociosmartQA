package services

import (
	"fmt"
	"smartgas-payment/internal/schemas"
	"strconv"
	"strings"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/refund"
)

//go:generate mockery --name StripeService --filename=mock_stripe.go --inpackage=true
type StripeService interface {
	CreateCustomer(*schemas.Customer) (string, error)
	CreatePaymentIntent(float64, string) (*stripe.PaymentIntent, error)
	CancelPaymentIntent(string) error
	ListPaymenthMethodsByCustomer(string) []*stripe.PaymentMethod
	MakeARefund(string, float64) (*stripe.Refund, error)
	DeletePaymentMethod(string) error
}

type stripeService struct{}

func ProvideStripeService() *stripeService {
	return &stripeService{}
}

func (ss *stripeService) DeletePaymentMethod(paymentMethodID string) error {
	_, err := paymentmethod.Detach(
		paymentMethodID,
		nil,
	)

	return err
}

func (ss *stripeService) CreateCustomer(cus *schemas.Customer) (string, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(cus.Fullname()),
		Phone: stripe.String(cus.PhoneNumber),
	}

	c, err := customer.New(params)
	if err != nil {
		// TODO: Log error on sentry
		return "", err
	}
	return c.ID, nil
}

func (ss *stripeService) CreatePaymentIntent(
	amount float64,
	customerID string,
) (*stripe.PaymentIntent, error) {
	amount2Decimals := fmt.Sprintf("%0.2f", amount)

	// replacing decimal
	amountToInt, _ := strconv.Atoi(strings.Replace(amount2Decimals, ".", "", -1))

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amountToInt)),
		Customer: stripe.String(customerID),
		Currency: stripe.String(string(stripe.CurrencyMXN)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		// TODO: Log error on sentry
		return nil, err
	}

	return pi, nil
}

func (ss *stripeService) CancelPaymentIntent(paymentIntentID string) error {
	_, err := paymentintent.Cancel(paymentIntentID, nil)
	return err
}

func (ss *stripeService) ListPaymenthMethodsByCustomer(customerID string) []*stripe.PaymentMethod {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String("card"),
	}

	i := paymentmethod.List(params)

	paymentMethods := make([]*stripe.PaymentMethod, 0)

	for i.Next() {
		paymentMethods = append(paymentMethods, i.PaymentMethod())
	}

	return paymentMethods
}

func (ss *stripeService) MakeARefund(transactionID string, amount float64) (*stripe.Refund, error) {
	amount2Decimals := fmt.Sprintf("%0.2f", amount)

	// replacing decimal
	amountToInt, _ := strconv.Atoi(strings.Replace(amount2Decimals, ".", "", -1))

	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(transactionID),
	}

	if amount > 0 {
		params.Amount = stripe.Int64(int64(amountToInt))
	}

	if amount == -3 {
		params.AddMetadata("status", "manual_action")
	}

	if amount == -2 {
		params.AddMetadata("status", "internal_cancellation")
	}

	return refund.New(params)
}
