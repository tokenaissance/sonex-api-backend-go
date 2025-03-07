package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/module/points"
	"sonex_api_backend_go/internal/repository"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/sonex_api_backend_go.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	c.Timeout = 0

	if time.Now().Before(repository.BadgeCampaignStop) {
		logx.Errorf("Cannot process before campaign stop.\n")
		return
	}

	accCondsSwaps := [3]uint32{
		repository.MilestoneCondsSwaps[0],
		repository.MilestoneCondsSwaps[0] + repository.MilestoneCondsSwaps[1],
		repository.MilestoneCondsSwaps[0] + repository.MilestoneCondsSwaps[1] + repository.MilestoneCondsSwaps[2],
	}

	accCondsStakes := [3]uint64{
		repository.MilestoneCondsStakes[0],
		repository.MilestoneCondsStakes[0] + repository.MilestoneCondsStakes[1],
		repository.MilestoneCondsStakes[0] + repository.MilestoneCondsStakes[1] + repository.MilestoneCondsStakes[2],
	}

	accPoints := [3]uint64{
		repository.MilestonePoints[0],
		repository.MilestonePoints[0] + repository.MilestonePoints[1],
		repository.MilestonePoints[0] + repository.MilestonePoints[1] + repository.MilestonePoints[2],
	}

	// get chain ID
	chainId, err := strconv.ParseInt(os.Getenv("CHAIN_ID"), 10, 64)
	if err != nil {
		logx.Errorf("CHAIN_ID invalid\n")
		return
	}

	conn := repository.NewConn(c.PostgreSQL.DataSource())
	rewardsModel := repository.NewSonexBadgeRewardsModel(conn)
	swapModel := repository.NewEventSonexPoolSwapModel(conn)
	stakedModel := repository.NewSonexBadgeStakedPositionModel(conn)
	stakes, err := stakedModel.AggregateValueByUser(context.Background())
	if err != nil {
		logx.Errorf("AggregateValueByUser: %v\n", err)
		return
	}
	logx.Infof("stakes=%v\n", len(stakes))

	for i, stake := range stakes {
		logx.Infof("#%d  user=%s\n", i+1, stake.UserAddress)
		swaps, err := swapModel.CountBySenderAndTimeRange(context.Background(), chainId, stake.UserAddress, repository.BadgeCampaignStart, repository.BadgeCampaignStop)
		if err != nil {
			logx.Errorf("CountBySenderAndTimeRange: %v\n", err)
			return
		}

		value, err := strconv.ParseInt(stake.Value, 10, 64)
		if err != nil {
			logx.Errorf("fails to parse: %v\n", err)
			return
		}

		for i := 2; i >= 0; i-- {
			if swaps < accCondsSwaps[i] || uint64(value) < accCondsStakes[i] {
				continue
			}
			p := int64(accPoints[i])
			_, err := rewardsModel.Insert(context.Background(), &repository.SonexBadgeRewards{
				UserAddress: stake.UserAddress,
				Points:      p,
				NftId:       0,
			})
			if err != nil {
				logx.Errorf("fails to insert SonexBadgeRewards: %v\n", err)
				return
			}

			// send point
			err = points.IncreasePoints(
				context.Background(),
				*c.PointsService,
				stake.UserAddress,
				p,
				c.Badge.SubTitle,
				fmt.Sprintf("badge-%s", stake.UserAddress),
				map[string]string{
					"user":      stake.UserAddress,
					"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
				},
			)
			if err != nil {
				logx.Errorf("fails to send points: %v\n", err)
				return
			}

			break
		}
	}
}
