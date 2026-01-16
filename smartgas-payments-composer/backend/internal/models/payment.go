package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID                    uuid.UUID `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	ExternalTransactionID string    `gorm:"column:external_transaction_id;type:varchar(255);not null;unique;"`
	Amount                float32   `gorm:"column:amount;type:float;not null;default:0;check:amount > -1;"`
	TotalLiter            float32   `gorm:"column:total_liter;type:float;not null;default:0;check:total_liter > -1;"`
	ChargeType            string    `gorm:"column:charge_type;type:enum('by_total', 'by_liter');not null;default:'by_total';"`
	RefundedAmount        float32   `gorm:"column:refunded_amount;type:float;not null;default:0;check:refunded_amount > -1;"`
	RealAmountReported    float32   `gorm:"column:real_amount_reported;type:float;"`
	ChargeFee             float32   `gorm:"column:charge_fee;type:float;default:0;not null;check:charge_fee > -1;"`

	FuelType         string  `gorm:"column:fuel_type;type:enum('regular', 'premium', 'diesel');not null;default:'regular';"`
	Price            float64 `gorm:"column:price;type:double;not null;default:0;check:price > -1;"`
	DiscountPerLiter float64 `gorm:"column:discount_per_liter;type:double;not null;default:0;check:discount_per_liter > -1;"`
	DiscountType     string  `gorm:"column:discount_type;type:enum('campaign','elegibility','none');not null;default:'none';"`

	Status          string `gorm:"column:status;type:enum('pending', 'paid', 'canceled', 'failed');not null;default:'pending';"`
	PaymentProvider string `gorm:"column:payment_provider;type:enum('stripe', 'swit', 'debit');not null;default:'stripe';"`

	GasPumpID *uuid.UUID `gorm:"column:gas_pump_id;type:varchar(36);"`
	GasPump   *GasPump   `gorm:"constraint:OnDelete:SET NULL;"`

	CustomerID      *uuid.UUID `gorm:"column:customer_id;type:varchar(36);"`
	Customer        *Customer  `gorm:"constraint:OnDelete:SET NULL;"`
	GMPoints        float32    `gorm:"column:gm_points;type:float;not null;default:0;check:gm_points > -1;"`
	GMID            string     `gorm:"column:external_gm_points_id;type:varchar(25);not null;default:''"`
	Invoiced        bool       `gorm:"column:invoiced;type:boolean;default:false;"`
	InvoiceID       string     `gorm:"column:external_invoice_id;type:varchar(100);not null;default:''"`
	CampaignID      *uuid.UUID `gorm:"column:campaign_id;type:varchar(36);"`
	Campaign        *Campaign  `gorm:"constraint:OnDelete:SET NULL;"`
	LevelID         *uuid.UUID `gorm:"column:elegibility_level_id;type:varchar(36);"`
	Level           *Level     `gorm:"foreignKey:LevelID;constraint:OnDelete:SET NULL;"`
	FromOperations  *bool      `gorm:"column:from_operations;type:boolean;not null;default:false;"`
	GiftCardKey     *string    `gorm:"column:gift_card_key;type:varchar(40);"`
	SetByEmployeeID *string    `gorm:"column:set_by_employee_id;type:varchar(20);"`

	Events []PaymentEvent

	gorm.Model
}

func (p *Payment) TableName() string {
	return "payments"
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()

	return
}

type PaymentEvent struct {
	ID                      uuid.UUID `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	Type                    string    `gorm:"column:type;type:enum('paid','funds_reserved', 'failed', 'canceled', 'pending', 'serving', 'serving_paused', 'served', 'partial_refund', 'pump_ready', 'internal_cancellation', 'manual_action');not null;default:'pending'"`
	PaymentID               uuid.UUID
	AuthorizedApplicationID *uuid.UUID
	AuthorizedApplication   *AuthorizedApplication
	CreatedAt               time.Time
}

func (pe *PaymentEvent) TableName() string {
	return "payment_events"
}

func (pe *PaymentEvent) BeforeCreate(tx *gorm.DB) (err error) {
	pe.ID = uuid.New()

	return
}
