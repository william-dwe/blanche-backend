package usecase

import (
	"fmt"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cronjob"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"github.com/rs/zerolog/log"
)

type MerchantAnalyticsUsecase interface {
	GetMerchantDashboardMerchantResponsivenessStatistics(username string, input dto.MerchantAnalyticsMerchantResponsivenessReqBody) ([]dto.MerchantAnalyticsMerchantResponsivenessResBody, error)
	GetMerchantDashboardSalesStatistics(username string, input dto.MerchantAnalyticsSalesReqBody) ([]dto.MerchantAnalyticsSalesResBody, error)
	GetMerchantDashboardCustomerSatisfactionStatistics(username string, input dto.MerchantAnalyticsCustomerSatisfactionReqBody) ([]dto.MerchantAnalyticsCustomerSatisfactionResBody, error)
	UpdateMerchantDashboard(input *dto.MerchantAnalyticsUpdateReqBody) error
}

type MerchantAnalyticsUsecaseConfig struct {
	MerchantAnalyticsRepository repository.MerchantAnalyticsRepository
	MerchantRepository          repository.MerchantRepository
}

type merchantAnalyticsUsecaseImpl struct {
	merchantAnalyticsRepository repository.MerchantAnalyticsRepository
	merchantRepository          repository.MerchantRepository
}

func NewMerchantAnalyticsUsecase(c MerchantAnalyticsUsecaseConfig) MerchantAnalyticsUsecase {
	cr := cronjob.GetCron()
	_, err := cr.AddJob("30 0 * * *", func() {
		errExe := c.MerchantAnalyticsRepository.UpdateMerchantDailyAnalytics(time.Now().Format(dateFormat))
		if errExe != nil {
			log.Error().Msg("error executing merchant analytics daily update")
		} else {
			log.Info().Msg("merchant analytics daily update executed")
		}
	})
	if err != nil {
		log.Error().Msg("error scheduling merchant analytics daily update")
	} else {
		log.Info().Msg("merchant analytics daily update scheduled")
	}

	return &merchantAnalyticsUsecaseImpl{
		merchantAnalyticsRepository: c.MerchantAnalyticsRepository,
		merchantRepository:          c.MerchantRepository,
	}
}

func (u *merchantAnalyticsUsecaseImpl) GetMerchantDashboardMerchantResponsivenessStatistics(username string, input dto.MerchantAnalyticsMerchantResponsivenessReqBody) ([]dto.MerchantAnalyticsMerchantResponsivenessResBody, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	startDate, endDate, err := parseInputDate(dto.DashboardReqBody(input))
	if err != nil {
		return nil, err
	}

	mAnalytics, err := u.merchantAnalyticsRepository.GetMerchantDailyAnalytics(merchant.Domain, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var activeUserStatistics []dto.MerchantAnalyticsMerchantResponsivenessResBody
	for _, m := range mAnalytics {
		activeUserStatistics = append(activeUserStatistics, dto.MerchantAnalyticsMerchantResponsivenessResBody{
			Type:  "OAD",
			Date:  m.DatePartition.Format(dateFormat),
			Value: fmt.Sprintf("%.2f", m.OAD),
		})
		activeUserStatistics = append(activeUserStatistics, dto.MerchantAnalyticsMerchantResponsivenessResBody{
			Type:  "OSD",
			Date:  m.DatePartition.Format(dateFormat),
			Value: fmt.Sprintf("%.2f", m.OSD),
		})
	}

	return activeUserStatistics, nil
}

func (u *merchantAnalyticsUsecaseImpl) GetMerchantDashboardSalesStatistics(username string, input dto.MerchantAnalyticsSalesReqBody) ([]dto.MerchantAnalyticsSalesResBody, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	startDate, endDate, err := parseInputDate(dto.DashboardReqBody(input))
	if err != nil {
		return nil, err
	}

	mAnalytics, err := u.merchantAnalyticsRepository.GetMerchantDailyAnalytics(merchant.Domain, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var salesStatistics []dto.MerchantAnalyticsSalesResBody
	for _, m := range mAnalytics {
		salesStatistics = append(salesStatistics, dto.MerchantAnalyticsSalesResBody{
			Date: m.DatePartition.Format(dateFormat),
			Rev:  m.Revenue,
			Trx:  m.TrxCount,
		})
	}

	return salesStatistics, nil
}

func (u *merchantAnalyticsUsecaseImpl) GetMerchantDashboardCustomerSatisfactionStatistics(username string, input dto.MerchantAnalyticsCustomerSatisfactionReqBody) ([]dto.MerchantAnalyticsCustomerSatisfactionResBody, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	startDate, endDate, err := parseInputDate(dto.DashboardReqBody(input))
	if err != nil {
		return nil, err
	}

	mAnalytics, err := u.merchantAnalyticsRepository.GetMerchantDailyAnalytics(merchant.Domain, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var salesStatistics []dto.MerchantAnalyticsCustomerSatisfactionResBody
	for _, m := range mAnalytics {
		salesStatistics = append(salesStatistics, dto.MerchantAnalyticsCustomerSatisfactionResBody{
			Date:   m.DatePartition.Format(dateFormat),
			Review: m.AvgReview,
			Count:  m.CountReview,
		})
	}

	return salesStatistics, nil
}

func (u *merchantAnalyticsUsecaseImpl) UpdateMerchantDashboard(input *dto.MerchantAnalyticsUpdateReqBody) error {
	err := u.merchantAnalyticsRepository.UpdateMerchantDailyAnalytics(input.DatePartition)
	return err
}
