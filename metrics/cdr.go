package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/webitel/cdr_metrics/gateway"
	"github.com/webitel/cdr_metrics/model"
	"github.com/webitel/wlog"
)

type MetricCdr struct {
	gw *gateway.Gateway

	inAudioMos          *prometheus.GaugeVec
	inRawBytes          *prometheus.GaugeVec
	inMediaBytes        *prometheus.GaugeVec
	inQualityPercentage *prometheus.GaugeVec
	inJitterLossRate    *prometheus.GaugeVec
	inJitterMaxVariance *prometheus.GaugeVec
	inJitterMinVariance *prometheus.GaugeVec

	outRawBytes   *prometheus.GaugeVec
	outMediaBytes *prometheus.GaugeVec

	billSec  *prometheus.GaugeVec
	duration *prometheus.GaugeVec

	hangupCode *prometheus.GaugeVec
}

func (m *MetricCdr) reset() {
	m.inAudioMos.Reset()
	m.inRawBytes.Reset()
	m.inMediaBytes.Reset()
	m.inQualityPercentage.Reset()
	m.inJitterLossRate.Reset()
	m.inJitterMaxVariance.Reset()
	m.inJitterMinVariance.Reset()

	m.outRawBytes.Reset()
	m.outMediaBytes.Reset()

	m.billSec.Reset()
	m.duration.Reset()

	m.hangupCode.Reset()

}

func NewCdr(space string, gw *gateway.Gateway) *MetricCdr {
	m := &MetricCdr{
		gw: gw,
	}

	m.inAudioMos = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_stats_audio_inbound_mos",
			Help:      "Cdr stats audio inbound mos",
		},
		model.Cdr{}.Labels(),
	)

	m.inRawBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_stats_audio_inbound_raw_bytes",
			Help:      "Cdr stats audio inbound raw_bytes",
		},
		model.Cdr{}.Labels(),
	)

	m.inMediaBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_stats_audio_inbound_media_bytes",
			Help:      "Cdr stats audio inbound media_bytes",
		},
		model.Cdr{}.Labels(),
	)

	m.inQualityPercentage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_stats_audio_inbound_quality_percentage",
			Help:      "Cdr stats audio inbound quality_percentage",
		},
		model.Cdr{}.Labels(),
	)

	m.inJitterLossRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_stats_audio_inbound_jitter_loss_rate",
			Help:      "Cdr stats audio inbound jitter_loss_rate",
		},
		model.Cdr{}.Labels(),
	)

	m.inJitterMaxVariance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_stats_audio_inbound_jitter_max_variance",
			Help:      "Cdr stats audio inbound jitter_max_variance",
		},
		model.Cdr{}.Labels(),
	)

	m.inJitterMinVariance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_stats_audio_inbound_jitter_min_variance",
			Help:      "Cdr stats audio inbound jitter_min_variance",
		},
		model.Cdr{}.Labels(),
	)

	m.outRawBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_stats_audio_outbound_raw_bytes",
			Help:      "Cdr stats audio outbound raw_bytes",
		},
		model.Cdr{}.Labels(),
	)

	m.outMediaBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_stats_audio_outbound_media_bytes",
			Help:      "Cdr stats audio outbound media_bytes",
		},
		model.Cdr{}.Labels(),
	)

	m.billSec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_bill_sec",
			Help:      "Cdr bill sec",
		},
		model.Cdr{}.Labels(),
	)

	m.duration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_duration_sec",
			Help:      "Cdr duration sec",
		},
		model.Cdr{}.Labels(),
	)

	m.hangupCode = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: space,
			Name:      "cdr_hangup_code",
			Help:      "Cdr hangup code",
		},
		model.Cdr{}.Labels(),
	)

	return m
}

func (m *MetricCdr) Push(cdr *model.Cdr) error {
	m.inAudioMos.WithLabelValues(cdr.Names()...).Add(cdr.InboundMos())
	m.inRawBytes.WithLabelValues(cdr.Names()...).Add(float64(cdr.InboundRawBytes()))
	m.inMediaBytes.WithLabelValues(cdr.Names()...).Add(float64(cdr.InboundMediaBytes()))
	m.inQualityPercentage.WithLabelValues(cdr.Names()...).Add(float64(cdr.InboundQualityPercentage()))

	m.inJitterLossRate.WithLabelValues(cdr.Names()...).Add(float64(cdr.InboundJitterLossRate()))
	m.inJitterMaxVariance.WithLabelValues(cdr.Names()...).Add(float64(cdr.InboundJitterMaxVariance()))
	m.inJitterMinVariance.WithLabelValues(cdr.Names()...).Add(float64(cdr.InboundJitterMinVariance()))

	m.outRawBytes.WithLabelValues(cdr.Names()...).Add(float64(cdr.OutboundRawBytes()))
	m.outMediaBytes.WithLabelValues(cdr.Names()...).Add(float64(cdr.OutboundMediaBytes()))

	m.billSec.WithLabelValues(cdr.Names()...).Add(float64(cdr.BillSec()))
	m.duration.WithLabelValues(cdr.Names()...).Add(float64(cdr.Duration()))

	m.hangupCode.WithLabelValues(cdr.Names()...).Add(float64(cdr.HangupCode()))

	if err := m.gw.Push(cdr.GetInstance(), m.inAudioMos, m.inRawBytes, m.inMediaBytes, m.inQualityPercentage, m.inJitterLossRate, m.inJitterMaxVariance,
		m.inJitterMinVariance, m.outRawBytes, m.outMediaBytes, m.billSec, m.duration, m.hangupCode); err != nil {
		wlog.Error(fmt.Sprintf("[cdr_metric] could not push completion time to gateway: %s", err.Error()))
	} else {
		wlog.Debug("[cdr_metric] send data - success")
	}

	m.reset()
	return nil
}
