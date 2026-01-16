package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Synchronization struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	Type      string    `gorm:"column:type;type:enum('gas_pumps', 'gas_stations', 'customer_levels');not null;default:'gas_stations';"`
	Status    string    `gorm:"column:status;type:enum('running', 'done');default:'running';not null;"`
	Details   []SynchronizationDetail
	Errors    []SynchronizationError
	CreatedAt time.Time
}

func (s *Synchronization) TableName() string {
	return "synchronizations"
}

func (s *Synchronization) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()

	return
}

type SynchronizationDetail struct {
	ID                uuid.UUID `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	ExternalID        string    `gorm:"column:external_id;type:varchar(10);not null;"`
	Action            string    `gorm:"column:action;type:enum('created', 'updated', 'error');not null;default:'created'"`
	ErrorText         string    `gorm:"column:error_text;type:varchar(500);not null;default:'';"`
	SynchronizationID uuid.UUID
	Data              string `gorm:"column:data;type:varchar(500);not null;default:''"`
	CreatedAt         time.Time
}

func (sd *SynchronizationDetail) TableName() string {
	return "synchronization_details"
}

func (sd *SynchronizationDetail) BeforeCreate(tx *gorm.DB) (err error) {
	sd.ID = uuid.New()

	return
}

type SynchronizationError struct {
	ID                uuid.UUID `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	Text              string    `gorm:"column:text;type:varchar(300);not null;default:'';"`
	SynchronizationID uuid.UUID
	CreatedAt         time.Time
}

func (se *SynchronizationError) TableName() string {
	return "synchronization_errors"
}

func (se *SynchronizationError) BeforeCreate(tx *gorm.DB) (err error) {
	se.ID = uuid.New()

	return
}
