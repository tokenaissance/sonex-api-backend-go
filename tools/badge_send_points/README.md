此工具用于 Badge 活动结束时发放积分奖励。

Badge PRD: https://www.notion.so/sxteam/SONEX-OG-Badge-1abdfef90e608084a800fcab30747cde

## 程序行为

计算积分奖励时，如果有奖励，则会记录到数据库 `sonex_badge_rewards` 表中。
如果此程序重新运行，会跳过已经记录奖励的用户。

积分奖励记录到数据库后， 会发送积分到积分微服务中。
如果发送失败，此程序会终止。

**此程序重新运行时，不会重新发送积分到积分微服务中！所以要记录并保存此程序的日志，便于以后对账。**

要空投 NFT 时，可以查询 `sonex_badge_rewards` (过滤条件：积分等于 `8000`).


## Run

```sh
# 编译
go build -C tools/badge_send_points -o /tmp/badge_send_points

# 执行
CHAIN_ID=1868 /tmp/badge_send_points -f <配置文件>
```
