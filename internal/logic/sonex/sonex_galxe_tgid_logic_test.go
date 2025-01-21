package sonex

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_GetTGInfo(t *testing.T) {
	GetTGInfo("https://sonex-devnet-cloudrun-sonex-mini-app-backend-acti-i4lu42aqsq-as.a.run.app///", 7103978900)
}

func Test_changeToCheckSumaddress(t *testing.T) {

	_userWithChecksum := "0xdA44a1ed127691decB80d8a6072e8B7957Fa6ab6"
	_userToLower := strings.ToLower(_userWithChecksum)
	_user_addressWithCheckSum := common.HexToAddress(_userToLower).Hex()
	logx.Info("_userWithChecksum:", _userWithChecksum, "_userToLower:", _userToLower, "_user_addressWithCheckSum:", _user_addressWithCheckSum)
}
