此工具用于以下一次性任务：

    从1月14日到 1月16日，每天一条volume24H的total 值，时间是UTC零点的，写入到sonex_vol_tvl_record.

相关任务: https://www.notion.so/sxteam/volume24H-total-PG-1a0dfef90e60808799b5c2728e2d45bf


## run

```sh
# 编译
go build -C tools/calc_volume_24h_14_15_16 -o /tmp/calc_volume_24h_14_15_16

# 执行
/tmp/calc_volume_24h_14_15_16 -f <配置文件>
```
