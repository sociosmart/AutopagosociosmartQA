package services

import "github.com/google/wire"

var ServicesSet = wire.NewSet(
	ProvideCustomerService,
	ProvideStripeService,
	ProvideSocioSmartService,
	ProvideSwitService,
	ProvideInvoicingService,
	ProvideMailService,
	ProvideDebitService,

	wire.Bind(new(CustomerService), new(*customerService)),
	wire.Bind(new(StripeService), new(*stripeService)),
	wire.Bind(new(SocioSmartService), new(*socioSmartService)),
	wire.Bind(new(SwitService), new(*switService)),
	wire.Bind(new(InvoicingService), new(*invoicingService)),
	wire.Bind(new(MailService), new(*mailService)),
	wire.Bind(new(DebitService), new(*debitService)),
)
