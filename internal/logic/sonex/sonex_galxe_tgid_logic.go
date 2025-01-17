package sonex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexGalxeTgidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexGalxeTgidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexGalxeTgidLogic {
	return &SonexGalxeTgidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexGalxeTgidLogic) SonexGalxeTgid(req *types.GalxeTgIdRequest) (resp *types.GalxeTgIdResponse, err error) {

	_galxeurl := ""
	if l.svcCtx.Config.Galxe != nil {
		_galxeurl = l.svcCtx.Config.Galxe.Url
	}
	_tgId := req.Tgid

	_isExisted, err := GetTGInfo(_galxeurl, _tgId)
	if err != nil {
		logx.Errorf("SonexGalxeTgid topic failed %v", err.Error())
		return &types.GalxeTgIdResponse{IsPlayedOrNot: false}, errors.New("invalid param")
	} else {
		return &types.GalxeTgIdResponse{IsPlayedOrNot: _isExisted}, nil
	}
}

func GetTGInfo(_miniappUrl string, _tgId int64) (bool, error) {
	_miniappUrl_noslash := strings.TrimSuffix(_miniappUrl, "/")

	_reqString := fmt.Sprintf("%v/%v%v", _miniappUrl_noslash, "api/1/user/info?user_id=", _tgId)
	logx.Infof("GetTGInfo:%v ", _reqString)
	req, _ := http.NewRequest("GET", _reqString, nil)

	// 再设置个json
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-API-Key", utils.RoundDexapisApiToken())
	var err error

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		logx.Errorf("GetTGInfo topic failed %v", err.Error())
	} else {
		defer resp.Body.Close()

		respByte, _ := io.ReadAll(resp.Body)

		var _retFromMiniApp JSONDataFromMiniApp
		err = json.Unmarshal(respByte, &_retFromMiniApp)
		if err != nil {
			logx.Errorf("GetTGInfo topic failed %v", err.Error())
			return false, err
		} else {
			return _retFromMiniApp.Data.IsExist, nil
		}
	}
	return false, nil
}

type JSONDataFromMiniApp struct {
	Status int  `json:"status"`
	Data   Data `json:"data"`
}
type User struct {
	UserID        int64  `json:"userId"`
	UserName      string `json:"userName"`
	PhotoURL      string `json:"photoURL"`
	ReferralCode  string `json:"referralCode"`
	ReferralUID   int    `json:"referralUid"`
	ReferralScore int    `json:"referralScore"`
}
type Data struct {
	User    User `json:"user"`
	IsExist bool `json:"isExist"`
}
