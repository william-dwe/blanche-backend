package repository

import (
	"strings"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type MerchantAnalyticsRepository interface {
	GetMerchantDailyAnalytics(merchantDomain string, startDate, endDate time.Time) ([]entity.MerchantDailyAnalyticsHist, error)
	UpdateMerchantDailyAnalytics(datePartition string) error
}

type MerchantAnalyticsRepositoryConfig struct {
	DB *gorm.DB
}

type mAnalyticsRepositoryImpl struct {
	db *gorm.DB
}

func NewMerchantAnalyticsRepository(c MerchantAnalyticsRepositoryConfig) MerchantAnalyticsRepository {
	return &mAnalyticsRepositoryImpl{
		db: c.DB,
	}
}

func (r *mAnalyticsRepositoryImpl) GetMerchantDailyAnalytics(merchantDomain string, startDate, endDate time.Time) ([]entity.MerchantDailyAnalyticsHist, error) {
	var mpAnalytics []entity.MerchantDailyAnalyticsHist
	err := r.db.Model(&mpAnalytics).
		Where("date_partition BETWEEN ? AND ?", startDate, endDate).
		Where("domain = ?", merchantDomain).
		Order("date_partition asc").
		Find(&mpAnalytics).Error
	if err != nil {
		return nil, domain.ErrGetMerchantAnalytics
	}

	return mpAnalytics, nil
}

func (r *mAnalyticsRepositoryImpl) UpdateMerchantDailyAnalytics(datePartition string) error {
	query := strings.ReplaceAll(`
	with cte_merchant_trx_and_revenue as (
		select 
			t.merchant_domain,
			SUM(cast(t.payment_details->>'subtotal' as decimal)) as revenue, 
			COUNT(t.id) as trx_count
		from transactions t 
		where t.created_at
			between DATE('{{selected_date}}') - interval '30 day'
			and DATE('{{selected_date}}')
		group by t.merchant_domain
	), cte_merchant_review as (
		select
			t.merchant_domain,
			avg(pr.rating) as avg_review, 
			count(pr.id) as count_review
		from product_reviews pr
		join transactions t
		on pr.transaction_id = t.id 
		and pr.created_at
			between DATE('{{selected_date}}') - interval '30 day'
			and DATE('{{selected_date}}')
		group by t.merchant_domain
	), cte_merchant_responsiveness as (
		select
			t.merchant_domain,
			extract(epoch from avg(ts.on_processed_at - ts.on_waited_at))/3600 as oad,
			extract(epoch from avg(ts.on_delivered_at - ts.on_waited_at))/3600 as osd
		from transaction_statuses ts
		join transactions t
		on t.id = ts.transaction_id 
		and t.created_at
			between DATE('{{selected_date}}') - interval '30 day'
			and DATE('{{selected_date}}')
		where ts.on_processed_at is not null and ts.on_delivered_at is not null
		group by t.merchant_domain 
	)
	INSERT INTO merchant_daily_analytics_hists
	(date_partition, "domain", revenue, trx_count, avg_review, count_review, oad, osd)
	
	select 
		DATE('{{selected_date}}') as date_partition, 
		m.domain,
		cmtar.revenue,
		cmtar.trx_count,
		cmr.avg_review,
		cmr.count_review,
		cmr2.oad,
		cmr2.osd
	from merchants m
	left join cte_merchant_trx_and_revenue cmtar
	on m.domain = cmtar.merchant_domain
	left join cte_merchant_review cmr
	on m.domain = cmr.merchant_domain
	left join cte_merchant_responsiveness cmr2
	on m.domain = cmr2.merchant_domain;`, "{{selected_date}}", datePartition)

	err := r.db.
		Exec(query).
		Error
	if err != nil {
		return domain.ErrUpdateMerchantAnalytics
	}

	return nil
}
