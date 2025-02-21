package task

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"testing"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func TestHandleACSSwapReward(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	err := global.Initialize(c)
	if err != nil {
		t.Error(err)
		return
	}
	global.InitIdGenerater(0)

	start := time.Now()
	HandleACSSwapReward(&c)
	elapsed := time.Since(start)
	fmt.Printf("Time token: %s\n", elapsed)
}

func TestSendACS(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	sqlx.DisableLog()
	err := global.Initialize(c)
	if err != nil {
		t.Error(err)
		return
	}
	global.InitIdGenerater(0)

	chainConf := c.GetChainConfig(1868)
	if chainConf == nil {
		t.Error(fmt.Errorf("HandleACSSwapReward error: %s", "Miss chain configure"))
		return
	}

	start := time.Now()
	conn := repository.NewConn(c.PostgreSQL.DataSource())
	acsSwapAcsRecordModel := repository.NewAcsSwapAcsRecordModel(conn)
	unsends, err := acsSwapAcsRecordModel.QueryUnsend(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	totalACS := linq.From(unsends).SelectT(func(item repository.AcsSwapAcsRecord) int64 { return item.Acs }).SumInts()
	logx.Infof("Total ACS: %d", totalACS)

	batchSize := 100
	batchCount := len(unsends) / int(batchSize)
	if len(unsends)%int(batchSize) > 0 {
		batchCount++
	}

	for i := 0; i < batchCount; i++ {
		start := i * int(batchSize)
		end := (i + 1) * int(batchSize)
		if end > len(unsends) {
			end = len(unsends)
		}
		currentUnsends := unsends[start:end]

		err := sendACS(context.Background(), c, currentUnsends)
		if err != nil {
			logx.Errorf("SendACS error: %s", err)
			continue
		}

	}

	elapsed := time.Since(start)
	fmt.Printf("Time token: %s\n", elapsed)
}
