package usecase

import (
	"strings"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cronjob"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"github.com/rs/zerolog/log"
)

type MarketplaceAnalyticsUsecase interface {
	GetMarketplaceDashboardActiveUserStatistics(input dto.MarketplaceAnalyticsActiveUserReqBody) ([]dto.MarketplaceAnalyticsActiveUserResBody, error)
	GetMarketplaceDashboardUserConversionStatistics(input dto.MarketplaceAnalyticsUserConversionReqBody) ([]dto.MarketplaceAnalyticsUserConversionResBody, error)
	GetMarketplaceDashboardSalesStatistics(input dto.MarketplaceAnalyticsSalesReqBody) ([]dto.MarketplaceAnalyticsSalesResBody, error)
	GetMarketplaceDashboardCustomerSatisfactionStatistics(input dto.MarketplaceAnalyticsCustomerSatisfactionReqBody) ([]dto.MarketplaceAnalyticsCustomerSatisfactionResBody, error)
	UpdateMarketplaceDashboard(input *dto.MarketplaceAnalyticsUpdateReqBody) error
}

type MarketplaceAnalyticsUsecaseConfig struct {
	MarketplaceAnalyticsRepository repository.MarketplaceAnalyticsRepository
}

type marketplaceAnalyticsUsecaseImpl struct {
	marketplaceAnalyticsRepository repository.MarketplaceAnalyticsRepository
}

const (
	dateFormat = "2006-01-02"
)

func NewMarketplaceAnalyticsUsecase(c MarketplaceAnalyticsUsecaseConfig) MarketplaceAnalyticsUsecase {
	cr := cronjob.GetCron()
	_, err := cr.AddJob("0 0 * * *", func() {
		errExe := c.MarketplaceAnalyticsRepository.UpdateMarketplaceDailyAnalytics(time.Now().Format(dateFormat))
		if errExe != nil {
			log.Error().Msg("error executing marketplace analytics daily update")
		} else {
			log.Info().Msg("marketplace analytics daily update executed")
		}
	})
	if err != nil {
		log.Error().Msg("error scheduling marketplace analytics daily update")
	} else {
		log.Info().Msg("marketplace analytics daily update scheduled")
	}

	return &marketplaceAnalyticsUsecaseImpl{
		marketplaceAnalyticsRepository: c.MarketplaceAnalyticsRepository,
	}
}

func parseInputDate(input dto.DashboardReqBody) (time.Time, time.Time, error) {
	strStartDate := strings.Trim(string(input.StartDate), "\"")
	startDate, err := time.Parse(dateFormat, strStartDate)
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrInvalidStartDateFormat
	}

	strEndDate := strings.Trim(string(input.EndDate), "\"")
	endDate, err := time.Parse(dateFormat, strEndDate)
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrInvalidStartDateFormat
	}

	return startDate, endDate, err
}

func (u *marketplaceAnalyticsUsecaseImpl) GetMarketplaceDashboardActiveUserStatistics(input dto.MarketplaceAnalyticsActiveUserReqBody) ([]dto.MarketplaceAnalyticsActiveUserResBody, error) {
	startDate, endDate, err := parseInputDate(dto.DashboardReqBody(input))
	if err != nil {
		return nil, err
	}

	mpAnalytics, err := u.marketplaceAnalyticsRepository.GetMarketplaceDailyAnalytics(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var activeUserStatistics []dto.MarketplaceAnalyticsActiveUserResBody
	for _, mp := range mpAnalytics {
		activeUserStatistics = append(activeUserStatistics, dto.MarketplaceAnalyticsActiveUserResBody{
			Type:  "MAU",
			Date:  mp.DatePartition.Format(dateFormat),
			Value: mp.CountDailyMau,
		})
		activeUserStatistics = append(activeUserStatistics, dto.MarketplaceAnalyticsActiveUserResBody{
			Type:  "MTU",
			Date:  mp.DatePartition.Format(dateFormat),
			Value: mp.CountDailyMtu,
		})
	}

	return activeUserStatistics, nil
}

func (u *marketplaceAnalyticsUsecaseImpl) GetMarketplaceDashboardUserConversionStatistics(input dto.MarketplaceAnalyticsUserConversionReqBody) ([]dto.MarketplaceAnalyticsUserConversionResBody, error) {
	startDate, endDate, err := parseInputDate(dto.DashboardReqBody(input))
	if err != nil {
		return nil, err
	}

	mpAnalytics, err := u.marketplaceAnalyticsRepository.GetMarketplaceDailyAnalytics(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var userConversionStatistics []dto.MarketplaceAnalyticsUserConversionResBody
	for _, mp := range mpAnalytics {
		userConversionStatistics = append(userConversionStatistics, dto.MarketplaceAnalyticsUserConversionResBody{
			Date:  mp.DatePartition.Format(dateFormat),
			Value: mp.UserConversionRate,
		})
	}

	return userConversionStatistics, nil
}

func (u *marketplaceAnalyticsUsecaseImpl) GetMarketplaceDashboardSalesStatistics(input dto.MarketplaceAnalyticsSalesReqBody) ([]dto.MarketplaceAnalyticsSalesResBody, error) {
	startDate, endDate, err := parseInputDate(dto.DashboardReqBody(input))
	if err != nil {
		return nil, err
	}

	mpAnalytics, err := u.marketplaceAnalyticsRepository.GetMarketplaceDailyAnalytics(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var salesStatistics []dto.MarketplaceAnalyticsSalesResBody
	for _, mp := range mpAnalytics {
		salesStatistics = append(salesStatistics, dto.MarketplaceAnalyticsSalesResBody{
			Date: mp.DatePartition.Format(dateFormat),
			Rev:  mp.Revenue,
			Trx:  mp.TrxCount,
		})
	}

	return salesStatistics, nil
}

func (u *marketplaceAnalyticsUsecaseImpl) GetMarketplaceDashboardCustomerSatisfactionStatistics(input dto.MarketplaceAnalyticsCustomerSatisfactionReqBody) ([]dto.MarketplaceAnalyticsCustomerSatisfactionResBody, error) {
	startDate, endDate, err := parseInputDate(dto.DashboardReqBody(input))
	if err != nil {
		return nil, err
	}

	mpAnalytics, err := u.marketplaceAnalyticsRepository.GetMarketplaceDailyAnalytics(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var salesStatistics []dto.MarketplaceAnalyticsCustomerSatisfactionResBody
	for _, mp := range mpAnalytics {
		salesStatistics = append(salesStatistics, dto.MarketplaceAnalyticsCustomerSatisfactionResBody{
			Date:   mp.DatePartition.Format(dateFormat),
			Review: mp.AvgReview,
			Count:  mp.CountReview,
		})
	}

	return salesStatistics, nil
}

func (u *marketplaceAnalyticsUsecaseImpl) UpdateMarketplaceDashboard(input *dto.MarketplaceAnalyticsUpdateReqBody) error {
	err := u.marketplaceAnalyticsRepository.UpdateMarketplaceDailyAnalytics(input.DatePartition)
	return err
}
