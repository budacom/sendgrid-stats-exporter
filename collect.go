package main

import (
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jinzhu/now"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	logger log.Logger

	blocks           *prometheus.Desc
	bounceDrops      *prometheus.Desc
	bounces          *prometheus.Desc
	clicks           *prometheus.Desc
	deferred         *prometheus.Desc
	delivered        *prometheus.Desc
	invalidEmails    *prometheus.Desc
	opens            *prometheus.Desc
	processed        *prometheus.Desc
	requests         *prometheus.Desc
	spamReportDrops  *prometheus.Desc
	spamReports      *prometheus.Desc
	uniqueClicks     *prometheus.Desc
	uniqueOpens      *prometheus.Desc
	unsubscribeDrops *prometheus.Desc
	unsubscribes     *prometheus.Desc
}

func collector(logger log.Logger) *Collector {
	var labels = []string{"user_name"}
	if *byCategoryMetrics {
		labels = append(labels, "category")
	}

	return &Collector{
		logger: logger,

		blocks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "blocks"),
			"blocks",
			labels,
			nil,
		),
		bounceDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bounce_drops"),
			"bounce_drops",
			labels,
			nil,
		),
		bounces: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bounces"),
			"bounces",
			labels,
			nil,
		),
		clicks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "clicks"),
			"clicks",
			labels,
			nil,
		),
		deferred: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "deferred"),
			"deferred",
			labels,
			nil,
		),
		delivered: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "delivered"),
			"delivered",
			labels,
			nil,
		),
		invalidEmails: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "invalid_emails"),
			"invalid_emails",
			labels,
			nil,
		),
		opens: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "opens"),
			"opens",
			labels,
			nil,
		),
		processed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "processed"),
			"processed",
			labels,
			nil,
		),
		requests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "requests"),
			"requests",
			labels,
			nil,
		),
		spamReportDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "spam_report_drops"),
			"spam_report_drops",
			labels,
			nil,
		),
		spamReports: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "spam_reports"),
			"spam_reports",
			labels,
			nil,
		),
		uniqueClicks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unique_clicks"),
			"unique_clicks",
			labels,
			nil,
		),
		uniqueOpens: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unique_opens"),
			"unique_opens",
			labels,
			nil,
		),
		unsubscribeDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unsubscribe_drops"),
			"unsubscribe_drops",
			labels,
			nil,
		),
		unsubscribes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unsubscribes"),
			"unsubscribes",
			labels,
			nil,
		),
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var today time.Time

	if *location != "" && *timeOffset != 0 {
		loc := time.FixedZone(*location, *timeOffset)
		today = time.Now().In(loc)
	} else {
		today = time.Now()
	}

	queryDate := today
	if *accumulatedMetrics {
		queryDate = now.With(today).BeginningOfMonth()
	}

	if !*byCategoryMetrics {
		statistics, err := collectByDate(queryDate, today)
		if err != nil {
			level.Error(c.logger).Log(err)
			return
		}

		for _, stats := range statistics[0].Stats {
			ch <- prometheus.MustNewConstMetric(
				c.blocks,
				prometheus.GaugeValue,
				float64(stats.Metrics.Blocks),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.bounceDrops,
				prometheus.GaugeValue,
				float64(stats.Metrics.BounceDrops),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.bounces,
				prometheus.GaugeValue,
				float64(stats.Metrics.Bounces),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.clicks,
				prometheus.GaugeValue,
				float64(stats.Metrics.Clicks),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deferred,
				prometheus.GaugeValue,
				float64(stats.Metrics.Deferred),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.delivered,
				prometheus.GaugeValue,
				float64(stats.Metrics.Delivered),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.invalidEmails,
				prometheus.GaugeValue,
				float64(stats.Metrics.InvalidEmails),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.opens,
				prometheus.GaugeValue,
				float64(stats.Metrics.Opens),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.processed,
				prometheus.GaugeValue,
				float64(stats.Metrics.Processed),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.requests,
				prometheus.GaugeValue,
				float64(stats.Metrics.Requests),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.spamReportDrops,
				prometheus.GaugeValue,
				float64(stats.Metrics.SpamReportDrops),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.spamReports,
				prometheus.GaugeValue,
				float64(stats.Metrics.SpamReports),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.uniqueClicks,
				prometheus.GaugeValue,
				float64(stats.Metrics.UniqueClicks),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.uniqueOpens,
				prometheus.GaugeValue,
				float64(stats.Metrics.UniqueOpens),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.unsubscribeDrops,
				prometheus.GaugeValue,
				float64(stats.Metrics.UnsubscribeDrops),
				*sendGridUserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.unsubscribes,
				prometheus.GaugeValue,
				float64(stats.Metrics.Unsubscribes),
				*sendGridUserName,
			)
		}
	} else {
		category_statistics, err := collectByCategory(today)

		if err != nil {
			level.Error(c.logger).Log(err)
			return
		}
		for _, stats := range category_statistics.Stats {
			ch <- prometheus.MustNewConstMetric(
				c.blocks,
				prometheus.GaugeValue,
				float64(stats.Metrics.Blocks),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.bounceDrops,
				prometheus.GaugeValue,
				float64(stats.Metrics.BounceDrops),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.bounces,
				prometheus.GaugeValue,
				float64(stats.Metrics.Bounces),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.clicks,
				prometheus.GaugeValue,
				float64(stats.Metrics.Clicks),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deferred,
				prometheus.GaugeValue,
				float64(stats.Metrics.Deferred),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.delivered,
				prometheus.GaugeValue,
				float64(stats.Metrics.Delivered),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.invalidEmails,
				prometheus.GaugeValue,
				float64(stats.Metrics.InvalidEmails),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.opens,
				prometheus.GaugeValue,
				float64(stats.Metrics.Opens),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.processed,
				prometheus.GaugeValue,
				float64(stats.Metrics.Processed),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.requests,
				prometheus.GaugeValue,
				float64(stats.Metrics.Requests),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.spamReportDrops,
				prometheus.GaugeValue,
				float64(stats.Metrics.SpamReportDrops),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.spamReports,
				prometheus.GaugeValue,
				float64(stats.Metrics.SpamReports),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.uniqueClicks,
				prometheus.GaugeValue,
				float64(stats.Metrics.UniqueClicks),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.uniqueOpens,
				prometheus.GaugeValue,
				float64(stats.Metrics.UniqueOpens),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.unsubscribeDrops,
				prometheus.GaugeValue,
				float64(stats.Metrics.UnsubscribeDrops),
				*sendGridUserName,
				stats.Category,
			)
			ch <- prometheus.MustNewConstMetric(
				c.unsubscribes,
				prometheus.GaugeValue,
				float64(stats.Metrics.Unsubscribes),
				*sendGridUserName,
				stats.Category,
			)
		}
	}
}
