package tasks

import (
	"errors"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/services"
	"smartgas-payment/internal/utils"
	"time"

	"github.com/jinzhu/now"

	"github.com/goccy/go-json"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

var RunningError = errors.New("There is a task already running")

//go:generate mockery --name SynchronizationTask --filename=mock_synchronization.go --inpackage=true
type SynchronizationTask interface {
	SyncGasStations() error
	SyncGasPumps() error
	GenerateElegibilityCustomers() error
}

type synchronizationTask struct {
	gasStationRepository      repository.GasStationRepository
	gasPumpRepository         repository.GasPumpRepository
	socioSmartService         services.SocioSmartService
	synchronizationRepository repository.SynchronizationRepository
	eleRepo                   repository.ElegibilityRepository
	cusRepo                   repository.CustomerRepository
	paymentRepository         repository.PaymentRepository
}

func ProvideSynchronizationTask(
	gasStationRepository repository.GasStationRepository,
	gasPumpRepository repository.GasPumpRepository,
	socioSmartService services.SocioSmartService,
	synchronizationRepository repository.SynchronizationRepository,
	eleRepo repository.ElegibilityRepository,
	cusRepo repository.CustomerRepository,
	paymentRepository repository.PaymentRepository,
) *synchronizationTask {
	return &synchronizationTask{
		gasStationRepository:      gasStationRepository,
		gasPumpRepository:         gasPumpRepository,
		socioSmartService:         socioSmartService,
		synchronizationRepository: synchronizationRepository,
		eleRepo:                   eleRepo,
		cusRepo:                   cusRepo,
		paymentRepository:         paymentRepository,
	}
}

func (st *synchronizationTask) SyncGasStations() error {
	sync, err := st.synchronizationRepository.GetLastByType("gas_stations")

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if sync != nil && sync.Status == "running" {
		return RunningError
	}

	// Create synchronization event
	syncModel := models.Synchronization{
		Type: "gas_stations",
	}

	err = st.synchronizationRepository.Create(&syncModel)

	if err != nil {
		// TODO: log error in sentry
		return err
	}

	stations, err := st.socioSmartService.GetGasStations()
	if err != nil {
		syncError := &models.SynchronizationError{
			Text:              err.Error(),
			SynchronizationID: syncModel.ID,
		}
		err = st.synchronizationRepository.CreateBatchErrors(
			[]*models.SynchronizationError{syncError},
		)

		if err != nil {
			// TODO: Log error on sentry
		}
	}

	syncDetails := make([]*models.SynchronizationDetail, 0)
	for _, station := range stations {
		var gasStation models.GasStation

		jsonData, _ := json.Marshal(station)

		copier.Copy(&gasStation, &station)

		created, err := st.gasStationRepository.GetByExternalIDOrCreate(
			station.ExternalID,
			&gasStation,
		)
		if err != nil {
			// TODO: Log the error on sentry

			syncDetails = append(syncDetails, &models.SynchronizationDetail{
				SynchronizationID: syncModel.ID,
				ExternalID:        gasStation.ExternalID,
				Action:            "error",
				Data:              string(jsonData),
				ErrorText:         err.Error(),
			})
			continue
		}

		action := "created"
		if !created {
			var gasStationToUpdate models.GasStation

			copier.Copy(&gasStationToUpdate, &station)

			_, err := st.gasStationRepository.UpdateByID(gasStation.ID, &gasStationToUpdate)
			if err != nil {
				syncDetails = append(syncDetails, &models.SynchronizationDetail{
					SynchronizationID: syncModel.ID,
					ExternalID:        gasStation.ExternalID,
					Action:            "error",
					ErrorText:         err.Error(),
					Data:              string(jsonData),
				})
				continue
			}
			action = "updated"
		}

		syncDetails = append(syncDetails, &models.SynchronizationDetail{
			SynchronizationID: syncModel.ID,
			ExternalID:        gasStation.ExternalID,
			Action:            action,
			Data:              string(jsonData),
		})
	}

	err = st.synchronizationRepository.CreateBatchDetails(syncDetails)

	if err != nil {
		// TODO: Log error on sentry
	}

	updated, err := st.synchronizationRepository.UpdateStatusByID(syncModel.ID, "done")
	if err != nil {
		// TODO: Log error on sentry
	}

	if !updated {
		// TODO: Log error on sentry
	}

	return nil
}

func (st *synchronizationTask) SyncGasPumps() error {
	sync, err := st.synchronizationRepository.GetLastByType("gas_pumps")

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if sync != nil && sync.Status == "running" {
		return RunningError
	}

	// Create synchronization event
	syncModel := models.Synchronization{
		Type: "gas_pumps",
	}

	err = st.synchronizationRepository.Create(&syncModel)

	if err != nil {
		// TODO: log error in sentry
		return err
	}

	stations, err := st.gasStationRepository.ListAll(map[string]any{})
	if err != nil {
		return err
	}

	syncErrors := make([]*models.SynchronizationError, 0)
	syncDetails := make([]*models.SynchronizationDetail, 0)

	// TODO: Log error
	for _, station := range stations {
		gasPumps, err := st.socioSmartService.GetGasPumpsByCrePermission(station.CrePermission)
		if err != nil {
			syncErrors = append(syncErrors, &models.SynchronizationError{
				SynchronizationID: syncModel.ID,
				Text:              err.Error(),
			})
			continue
		}

		for _, gasPump := range gasPumps {
			var pump models.GasPump

			copier.Copy(&pump, &gasPump)

			pump.GasStation = station

			jsonData, _ := json.Marshal(gasPump)

			created, err := st.gasPumpRepository.GetByExternalIDOrCreate(pump.ExternalID, &pump)
			if err != nil {
				syncDetails = append(syncDetails, &models.SynchronizationDetail{
					SynchronizationID: syncModel.ID,
					ExternalID:        gasPump.ExternalID,
					Action:            "error",
					Data:              string(jsonData),
					ErrorText:         err.Error(),
				})
				continue
			}

			action := "created"
			if !created {
				var pumpToUpdate models.GasPump
				pumpToUpdate.GasStationID = &station.ID

				copier.Copy(&pumpToUpdate, &gasPump)

				_, err := st.gasPumpRepository.UpdateByID(pump.ID, &pumpToUpdate)
				if err != nil {
					syncDetails = append(syncDetails, &models.SynchronizationDetail{
						SynchronizationID: syncModel.ID,
						ExternalID:        gasPump.ExternalID,
						Action:            "error",
						Data:              string(jsonData),
						ErrorText:         err.Error(),
					})
					continue
				}

				action = "updated"
			}

			syncDetails = append(syncDetails, &models.SynchronizationDetail{
				SynchronizationID: syncModel.ID,
				ExternalID:        gasPump.ExternalID,
				Action:            action,
				Data:              string(jsonData),
			})

		}

	}

	err = st.synchronizationRepository.CreateBatchErrors(syncErrors)

	if err != nil {
		// TODO: Log error on sentry
	}

	err = st.synchronizationRepository.CreateBatchDetails(syncDetails)

	if err != nil {
		// TODO: Log error on sentry
	}

	updated, err := st.synchronizationRepository.UpdateStatusByID(syncModel.ID, "done")
	if err != nil {
		// TODO: Log error on sentry
	}

	if !updated {
		// TODO: Log error on sentry
	}

	return nil
}

func (st *synchronizationTask) GenerateElegibilityCustomers() error {
	sync, err := st.synchronizationRepository.GetLastByType("customer_levels")

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if sync != nil && sync.Status == "running" {
		return RunningError
	}
	// Sub a month to current date in order to calculate
	// customer levels for prev month
	// example:
	// current date: 2023-11-11
	// will be: 2023-10-11
	loc, _ := time.LoadLocation("America/Mazatlan")
	n := now.With(time.Now().In(loc).AddDate(0, -1, 0))
	now := time.Now()
	lowDate := n.BeginningOfMonth()
	highDate := n.EndOfMonth()

	levels, err := st.eleRepo.LevelListAllActive()
	if err != nil {
		return err
	}

	customers, err := st.cusRepo.ListAll()
	if err != nil {
		return err
	}

	syncModel := models.Synchronization{
		Type: "customer_levels",
	}

	st.synchronizationRepository.Create(&syncModel)
	for _, customer := range customers {
		cusLevel, err := st.eleRepo.GetCustomerLevelByCriterias(map[string]any{
			"customer_id":    customer.ID,
			"validity_month": now.Month(),
			"validity_year":  now.Year(),
		})
		if err != nil {
			// TODO: Log Error
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
		}

		if cusLevel != nil && *cusLevel.ManuallyTouched {
			continue
		}

		opts := repository.StatsForCustomerOpts{
			HighDate: highDate,
			LowDate:  lowDate,
		}

		stats, err := st.paymentRepository.GetStatsForCustomer(customer.ID, opts)
		if err != nil {
			// TODO: Do something with the error
		}

		// Check closest level
		var levelReached *models.Level

		for _, level := range levels {
			if stats.TotalReported >= *level.MinAmount &&
				stats.TotalTransactions >= *level.MinCharges {
				levelReached = level
			}
		}

		if levelReached == nil {
			continue
		}

		if cusLevel != nil {
			// Update the level
			cusLevel.LevelID = &levelReached.ID
			st.eleRepo.UpdateCustomerLevelByID(cusLevel.ID, cusLevel)
		} else {
			// Create the level
			cusLevelM := models.CustomerLevel{
				CustomerID:    &customer.ID,
				LevelID:       &levelReached.ID,
				ValidityMonth: utils.IntAddr(int(now.Month())),
				ValidityYear:  utils.IntAddr(now.Year()),
			}
			st.eleRepo.CreateCustomerLevel(&cusLevelM)
		}

	}

	_, err = st.synchronizationRepository.UpdateStatusByID(syncModel.ID, "done")

	if err != nil {
		return err
	}
	return nil
}
