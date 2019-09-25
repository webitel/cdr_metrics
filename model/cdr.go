package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	EMPTY_VALUE = "<Unknown>"
)

type CdrAudioInbound struct {
	Mos      float64 `json:"mos"`
	RawBytes int     `json:"raw_bytes"`
	//FlawTotal         int     `json:"flaw_total"`
	MediaBytes int `json:"media_bytes"`
	//PacketCount       int     `json:"packet_count"`
	//MeanInterval      float32 `json:"mean_interval"`
	//LargestJbSize     int     `json:"largest_jb_size"`
	//CngPacketCount    int     `json:"cng_packet_count"`
	JitterLossRate int `json:"jitter_loss_rate"`
	//DtmfPacketCount   int     `json:"dtmf_packet_count"`
	//JitterBurstRate   int     `json:"jitter_burst_rate"`
	//SkipPacketCount   int     `json:"skip_packet_count"`
	//FlushPacketCount  int     `json:"flush_packet_count"`
	//MediaPacketCount  int     `json:"media_packet_count"`
	QualityPercentage int     `json:"quality_percentage"`
	JitterMaxVariance float64 `json:"jitter_max_variance"`
	JitterMinVariance float64 `json:"jitter_min_variance"`
	//JitterPacketCount int     `json:"jitter_packet_count"`
}

type CdrAudioOutbound struct {
	RawBytes   int `json:"raw_bytes"`
	MediaBytes int `json:"media_bytes"`
	//PacketCount      int `json:"packet_count"`
	//CngPacketCount   int `json:"cng_packet_count"`
	//RtcpOctetCount   int `json:"rtcp_octet_count"`
	//DtmfPacketCount  int `json:"dtmf_packet_count"`
	//RtcpPacketCount  int `json:"rtcp_packet_count"`
	//SkipPacketCount  int `json:"skip_packet_count"`
	//MediaPacketCount int `json:"media_packet_count"`
}

type CdrAudio struct {
	Inbound  CdrAudioInbound  `json:"inbound"`
	Outbound CdrAudioOutbound `json:"outbound"`
}

type CdrStats struct {
	Audio CdrAudio `json:"audio"`
}

type CdrVariables map[string]interface{}

func (c CdrVariables) GetString(name string) (string, bool) {
	if tmp, ok := c[name]; ok {
		return fmt.Sprintf("%s", tmp), true
	}
	return "", false
}

func (c CdrVariables) GetInt(name string) (int, bool) {
	if tmp, ok := c.GetString(name); ok {
		i, _ := strconv.Atoi(tmp)
		return i, true
	}
	return 0, false
}

type Cdr struct {
	Stats     CdrStats     `json:"callStats"`
	Variables CdrVariables `json:"variables"`
}

func (cdr *Cdr) InboundMos() float64 {
	return cdr.Stats.Audio.Inbound.Mos
}

func (cdr *Cdr) InboundRawBytes() int {
	return cdr.Stats.Audio.Inbound.RawBytes
}

func (cdr *Cdr) InboundMediaBytes() int {
	return cdr.Stats.Audio.Inbound.MediaBytes
}

func (cdr *Cdr) InboundQualityPercentage() int {
	return cdr.Stats.Audio.Inbound.QualityPercentage
}

func (cdr *Cdr) InboundJitterLossRate() int {
	return cdr.Stats.Audio.Inbound.JitterLossRate
}

func (cdr *Cdr) InboundJitterMaxVariance() float64 {
	return cdr.Stats.Audio.Inbound.JitterMaxVariance
}

func (cdr *Cdr) InboundJitterMinVariance() float64 {
	return cdr.Stats.Audio.Inbound.JitterMinVariance
}

func (cdr *Cdr) OutboundRawBytes() int {
	return cdr.Stats.Audio.Outbound.RawBytes
}

func (cdr *Cdr) OutboundMediaBytes() int {
	return cdr.Stats.Audio.Outbound.MediaBytes
}

func (cdr *Cdr) BillSec() int {
	if i, ok := cdr.Variables.GetInt("billsec"); ok {
		return i
	}
	return 0
}

func (cdr *Cdr) Duration() int {
	if i, ok := cdr.Variables.GetInt("duration"); ok {
		return i
	}
	return 0
}

func (cdr *Cdr) HangupCode() int {
	if i, ok := cdr.Variables.GetInt("hangup_cause_q850"); ok {
		return i
	}
	return 0
}

func CdrFromJson(data []byte) *Cdr {
	var cdr *Cdr
	json.Unmarshal(data, &cdr)
	return cdr
}

func (cdr *Cdr) GetInstance() string {
	var tmp string
	var ok bool

	if tmp, ok = cdr.Variables.GetString("sip_local_network_addr"); ok {
		return tmp
	}

	if tmp, ok = cdr.Variables.GetString("local_media_ip"); ok {
		return tmp
	}
	return EMPTY_VALUE
}

func (c *Cdr) UserAgent() string {
	var tmp string
	var ok bool
	if tmp, ok = c.Variables.GetString("sip_user_agent"); ok {
		return tmp
	}
	if tmp, ok = c.Variables.GetString("verto_user_agent"); ok {
		return tmp
	}
	return EMPTY_VALUE
}

func (c Cdr) Labels() []string {
	return []string{"uuid", "call_uuid", "sip_network_ip", "remote_media_ip", "user_agent"}
}

func (c *Cdr) Names() []string {
	uuid, _ := c.Variables.GetString("uuid")
	callUUid, _ := c.Variables.GetString("call_uuid")
	sipNetworkIp, _ := c.Variables.GetString("sip_network_ip")
	remoteAudioIp, _ := c.Variables.GetString("remote_media_ip")
	return []string{uuid, callUUid, sipNetworkIp, remoteAudioIp, c.UserAgent()}
}
