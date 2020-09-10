package config

type Kafka struct {
	Log_folder_name string
	Brokers         []string
	Topic_prefix    string
	Partitions      []int
}

type Redis struct {
	Address []string
}

type Session struct {
	Log_folder_path string
}

type Sigdos struct {
	Log_folder_path      string
	A                    float64
	H                    float64
	PacketCountThreshold int
	PreBCT               string
	PreBCTT              string
}

type Scan struct {
	PreBCT   string
	Var2     float64
	Var1     float64
	UpTrw    float64
	DownTrw  float64
	Boundary int
}

type SynStat struct {
	PreBCT string
}

type NAS struct {
	PreBCT string
}

type StreamData struct {
	DataGTPC string
	DataGTPU string
	DataB    string
	DataNAS  string
	DataUU   string
	DataCCC  string
	DataUD   string
	DataUT   string
}

type GTPCFieldIndex struct {
	IdxImsi              int
	IdxSessoniID         int
	IdxRequestType       int
	IdxResponseTime      int
	IdxResponseDirection int
	IdxSummaryCreateTime int
}

type GTPUFieldIndex struct {
	IdxImsi              int
	IdxSessionID         int
	IdxBearerCreatedTime int
	IdxPacketCount       int
	IdxSynDirection      int
	IdxRstDirection      int
	IdxAppIp             int
	IdxAppPort           int
	IdxUserIp            int
	IdxuserPort          int
	IdxSynCount          int
	IdxSynAckCount       int
	IdxRstCount          int
	IdxuplinkFirsttime   int
	IdxSummaryCreateTime int
	IdxSynAckFirsttime   int
	IdxRstFirsttime      int
}

type GTPBFieldIndex struct {
	IdxImsi              int
	IdxSessionID         int
	IdxBearerCreatedTime int
	IdxPacketCount       int
	IdxSummaryCreateTime int
	IdxUpPacketCount     int
}

type NASFieldIndex struct {
	IdxSummaryCreateTime  int
	IdxImsi               int
	IdxSecurityModeReject int
	IdxCallType           int
	IdxUEcapability       int
	IdxAttachRejectCause  int
}

type CCCFieldIndex struct {
	IdxImsi              int
	IdxSessoniID         int
	IdxRequestType       int
	IdxResponseTime      int
	IdxResponseDirection int
}

type UTFieldIndex struct {
	IdxImsi              int
	IdxSessionID         int
	IdxBearerCreatedTime int
	IdxPacketCount       int
	IdxSynDirection      int
	IdxRstDirection      int
	IdxAppIp             int
	IdxAppPort           int
	IdxUserIp            int
	IdxuserPort          int
	IdxSynCount          int
	IdxSynAckCount       int
	IdxRstCount          int
	IdxuplinkFirsttime   int
}

type UDFieldIndex struct {
	IdxImsi              int
	IdxSessionID         int
	IdxBearerCreatedTime int
	IdxPacketCount       int
	IdxSynDirection      int
	IdxRstDirection      int
	IdxAppIp             int
	IdxAppPort           int
	IdxUserIp            int
	IdxuserPort          int
	IdxSynCount          int
	IdxSynAckCount       int
	IdxRstCount          int
	IdxuplinkFirsttime   int
}

type Config struct {
	ABS_PATH string
	Kafka
	Redis
	Session
	Sigdos
	Scan
	StreamData

	GTPCFieldIndex
	GTPUFieldIndex
	GTPBFieldIndex
	NASFieldIndex

	CCCFieldIndex
	UTFieldIndex
	UDFieldIndex
}
