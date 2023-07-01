-- marketplace daily analyticsc
create table marketplace_daily_analytics_hist (
	date_partition date,
	count_daily_mau int,
	count_daily_mtu int,
	user_convertion_rate decimal,
	revenue numeric,
	trx_count int,
	avg_review decimal,
	count_review int
);

with cte_daily_mau as (
	select count(distinct ula.user_id) as count_daily_mau
	from user_login_activities ula 
	where ula.created_at 
		between DATE('2023-02-23') - interval '30 day'
		and DATE('2023-02-23')
), cte_daily_mtu as (
	select count(distinct ula.user_id) as count_daily_mtu
	from user_login_activities ula 
	join transactions t 
	on ula.user_id = t.user_id 
	and t.created_at
		between DATE('2023-02-23') - interval '30 day'
		and DATE('2023-02-23')
	where ula.created_at 
		between DATE('2023-02-23') - interval '30 day'
		and DATE('2023-02-23')
), cte_trx_and_revenue as (
	select 
		SUM(cast(t.payment_details->>'subtotal' as decimal)) as revenue, 
		COUNT(t.id) as trx_count
	from transactions t
	where t.created_at
		between DATE('2023-02-23') - interval '30 day'
		and DATE('2023-02-23')
), cte_review as (
	select avg(pr.rating) as avg_review, count(pr.id) as count_review
	from product_reviews pr
	join transactions t
	on pr.transaction_id = t.id 
	and pr.created_at
		between DATE('2023-02-23') - interval '30 day'
		and DATE('2023-02-23')
), cte_daily_analytics as (
	select 
		DATE('2023-02-23') as date_partition, 
		(select count_daily_mau from cte_daily_mau), 
		(select count_daily_mtu from cte_daily_mtu),
		(select revenue from cte_trx_and_revenue),
		(select trx_count from cte_trx_and_revenue),
		(select avg_review from cte_review),
		(select count_review from cte_review)
)

insert into marketplace_daily_analytics_hist (
	date_partition,
	count_daily_mau, 
	count_daily_mtu,
	user_convertion_rate,
	revenue,
	trx_count,
	avg_review,
	count_review
) 
select 
	date_partition, 
	count_daily_mau, 
	count_daily_mtu,
	cast(count_daily_mtu as decimal)/count_daily_mau as user_convertion_rate,
	revenue,
	trx_count,
	avg_review,
	count_review
from cte_daily_analytics; 


-- merchant daily analytics
create table merchant_daily_analytics_hist (
	date_partition date,
	domain varchar,
	revenue numeric,
	trx_count int,
	avg_review decimal,
	count_review int,
	oad interval,
	osd interval
);

with cte_merchant_trx_and_revenue as (
	select 
		t.merchant_domain,
		SUM(cast(t.payment_details->>'subtotal' as decimal)) as revenue, 
		COUNT(t.id) as trx_count
	from transactions t 
	where t.created_at
		between DATE('2023-02-23') - interval '30 day'
		and DATE('2023-02-23')
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
		between DATE('2023-02-23') - interval '30 day'
		and DATE('2023-02-23')
	group by t.merchant_domain
), cte_merchant_responsiveness as (
	select
		t.merchant_domain,
		avg(ts.on_processed_at - ts.on_waited_at) as oad,
		avg(ts.on_delivered_at - ts.on_waited_at) as osd
	from transaction_statuses ts
	join transactions t
	on t.id = ts.transaction_id 
	and t.created_at
		between DATE('2023-02-23') - interval '30 day'
		and DATE('2023-02-23')
	where ts.on_processed_at is not null and ts.on_delivered_at is not null
	group by t.merchant_domain 
)
INSERT INTO public.merchant_daily_analytics_hist
(date_partition, "domain", revenue, trx_count, avg_review, count_review, oad, osd)

select 
	DATE('2023-02-23') as date_partition, 
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
on m.domain = cmr2.merchant_domain;

