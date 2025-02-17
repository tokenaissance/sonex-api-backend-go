package points

import (
	"context"
	"fmt"
	"net/http"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/util"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpc"
)

func IncreasePoints(ctx context.Context, conf config.PointsServiceConfig, user string, points int64, subTitle string, extraId string, extra map[string]string) error {
	req := struct {
		User      string            `json:"user_address"`
		Points    int64             `json:"points"`
		Source    string            `json:"source"`
		SubTitle  string            `json:"sub_title"`
		ExtraId   string            `json:"extra_id"`
		Extra     map[string]string `json:"extra"`
		Signature string            `json:"signature"`
	}{
		User:     user,
		Points:   points,
		Source:   conf.Source,
		SubTitle: subTitle,
		ExtraId:  extraId,
		Extra:    extra,
	}
	signature, err := util.Sign(&req, conf.PrivateKey)
	if err != nil {
		return fmt.Errorf("%s", fmt.Sprintf("Sign failed: %s", err))
	}
	req.Signature = signature

	url := fmt.Sprintf("%s/extra/increase_points", strings.TrimRight(conf.Host, "/"))
	resp, err := httpc.Do(ctx, http.MethodPost, url, req)
	if err != nil {
		return fmt.Errorf("%s", fmt.Sprintf("Increase points failed: %s %v", err, req))
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", fmt.Sprintf("Increase points failed: %s %v", resp.Status, req))
	}

	return nil
}
