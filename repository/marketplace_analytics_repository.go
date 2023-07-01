package repository

import (
	"strings"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type MarketplaceAnalyticsRepository interface {
	GetMarketplaceDailyAnalytics(startDate, endDate time.Time) ([]entity.MarketplaceDailyAnalyticsHist, error)
	UpdateMarketplaceDailyAnalytics(datePartition string) error
}

type MarketplaceAnalyticsRepositoryConfig struct {
	DB *gorm.DB
}

type mpAnalyticsRepositoryImpl struct {
	db *gorm.DB
}

func NewMarketplaceAnalyticsRepository(c MarketplaceAnalyticsRepositoryConfig) MarketplaceAnalyticsRepository {
	return &mpAnalyticsRepositoryImpl{
		db: c.DB,
	}
}

func (r *mpAnalyticsRepositoryImpl) GetMarketplaceDailyAnalytics(startDate, endDate time.Time) ([]entity.MarketplaceDailyAnalyticsHist, error) {
	var mpAnalytics []entity.MarketplaceDailyAnalyticsHist
	err := r.db.Model(&mpAnalytics).
		Where("date_partition BETWEEN ? AND ?", startDate, endDate).
		Order("date_partition asc").
		Find(&mpAnalytics).Error
	if err != nil {
		return nil, domain.ErrGetMarketplaceAnalytics
	}

	return mpAnalytics, nil
}

func (r *mpAnalyticsRepositoryImpl) UpdateMarketplaceDailyAnalytics(datePartition string) error {
	query := strings.ReplaceAll(`
	with cte_daily_mau as (
		select count(distinct ula.user_id) as count_daily_mau
		from user_login_activities ula 
		where ula.created_at 
			between DATE('{{selected_date}}') - interval '30 day'
			and DATE('{{selected_date}}')
	), cte_daily_mtu as (
		select count(distinct ula.user_id) as count_daily_mtu
		from user_login_activities ula 
		join transactions t 
		on ula.user_id = t.user_id 
		and t.created_at
			between DATE('{{selected_date}}') - interval '30 day'
			and DATE('{{selected_date}}')
		where ula.created_at 
			between DATE('{{selected_date}}') - interval '30 day'
			and DATE('{{selected_date}}')
	), cte_trx_and_revenue as (
		select 
			SUM(cast(t.payment_details->>'subtotal' as decimal)) as revenue, 
			COUNT(t.id) as trx_count
		from transactions t
		where t.created_at
			between DATE('{{selected_date}}') - interval '30 day'
			and DATE('{{selected_date}}')
	), cte_review as (
		select avg(pr.rating) as avg_review, count(pr.id) as count_review
		from product_reviews pr
		join transactions t
		on pr.transaction_id = t.id 
		and pr.created_at
			between DATE('{{selected_date}}') - interval '30 day'
			and DATE('{{selected_date}}')
	), cte_daily_analytics as (
		select 
			DATE('{{selected_date}}') as date_partition, 
			(select count_daily_mau from cte_daily_mau), 
			(select count_daily_mtu from cte_daily_mtu),
			(select revenue from cte_trx_and_revenue),
			(select trx_count from cte_trx_and_revenue),
			(select avg_review from cte_review),
			(select count_review from cte_review)
	)

	insert into marketplace_daily_analytics_hists (
		date_partition,
		count_daily_mau, 
		count_daily_mtu,
		user_conversion_rate,
		revenue,
		trx_count,
		avg_review,
		count_review
	) 
	select 
		date_partition, 
		count_daily_mau, 
		count_daily_mtu,
		cast(count_daily_mtu as decimal)/count_daily_mau as user_conversion_rate,
		revenue,
		trx_count,
		avg_review,
		count_review
	from cte_daily_analytics`, "{{selected_date}}", datePartition)

	var mpAnalytics entity.MarketplaceDailyAnalyticsHist
	err := r.db.
		Raw(query).
		Scan(&mpAnalytics).
		Error
	if err != nil {
		return domain.ErrUpdateMarketplaceAnalytics
	}

	return nil
}
