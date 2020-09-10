package module

import (
	"fmt"
	"log"
	"strings"
	"time"

	"../db"
	"github.com/Shopify/sarama"
)

type exConsGrouphandler struct{}

var (
	exConsLogger *log.Logger
)

func (exConsGrouphandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (exConsGrouphandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h exConsGrouphandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Printf("ex start")
	// init configs
	configPath := "/home/jpjp/project/test/0602/config.ini"
	cfg = initConfig(configPath)
	redisAddr := cfg.Redis.Address
	///////////get features from config
	//idxImsi := cfg.GTPUFieldIndex.IdxImsi

	/////////initialize exCons
	exCons := InitexCons(redisAddr, preBCT, idxImsi, idxSynDirection,
		idxUserIp, idxSynCount)

	idx := 0
	// run session module
	for msg := range claim.Messages() {
		//fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		startTime := time.Now()
		record := strings.Split(string(msg.Value), ",")
		//
		//////consumer 처리
		//
		sess.MarkMessage(msg, "")
		idx++
		//exConsLogger.Println("markd message:", idx, "  partition:", id)
	}
	return nil
}

//ModexCons ..
type exCons struct {
	preBCT      string
	redisClient *db.RedisController
	// necessary gtpu index info

}

// init function
func InitexCons(redisAddr []string, preBCT string, idxImsi int, idxSynDirection int,
	idxUserIp int,
	idxSynCount int,
) *exCons {
	exCons := &exCons{
		redisClient: db.NewRedisController(redisAddr),
		preBCT:      preBCT,

		idxUserIp:   idxUserIp,
		idxSynCount: idxSynCount,
	}
	return exCons
}
