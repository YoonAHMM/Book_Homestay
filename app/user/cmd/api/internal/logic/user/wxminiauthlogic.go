package user

import (
	"context"

	"Book_Homestay/app/user/cmd/api/internal/svc"
	"Book_Homestay/app/user/cmd/api/internal/types"
	"Book_Homestay/app/user/cmd/rpc/user"
	"Book_Homestay/common/errx"
	"Book_Homestay/common/vars"

	"github.com/pkg/errors"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type WxMiniAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWxMiniAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WxMiniAuthLogic {
	return &WxMiniAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WxMiniAuthLogic) WxMiniAuth(req *types.WXMiniAuthReq) (resp *types.WXMiniAuthResp, err error) {
	

	program :=wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.Secret,
		Cache:     cache.NewMemory(),
	})

	//换取 用户唯一标识 OpenID 、 用户在微信开放平台账号下的唯一标识UnionID（若当前小程序已绑定到微信开放平台账号） 和 会话密钥 session_key
	seesion,err:=program.GetAuth().Code2Session(req.Code)

	if err != nil || seesion.ErrCode != 0 ||seesion.OpenID == "" {
		return nil,errors.Wrap(errx.NewErrCode(errx.WXMINI_ERROR),seesion.ErrMsg)
	}

	Data, err := program.GetEncryptor().Decrypt(seesion.SessionKey, req.EncryptedData, req.IV)
	if err != nil {
		return nil,errors.Wrap(errx.NewErrCode(errx.WXMINI_ERROR),seesion.ErrMsg)
	}

	//未绑定的则绑定，绑定的直接生成token
	if seesion.UnionID != ""{
		rpcRsp, err := l.svcCtx.UserRpc.GetUserAuthByAuthKey(l.ctx, &user.GetUserAuthByAuthKeyReq{
			AuthType: vars.UserAuthTypeSmallWX,
			AuthKey:  seesion.OpenID,
		})
		if err != nil {
			return nil, err
		}
		userId := rpcRsp.UserAuth.UserId
		tokenResp, err := l.svcCtx.UserRpc.GenerateToken(l.ctx, &user.GenerateTokenReq{
			UserId: userId,
		})

		if err != nil {
			return nil, err
		}
		return &types.WXMiniAuthResp{
			AccessToken:  tokenResp.AccessToken,
			AccessExpire: tokenResp.AccessExpire,
			RefreshAfter: tokenResp.RefreshAfter,
		}, nil
		}else{
			mobile := Data.PhoneNumber
			nickName := Data.NickName
			registerRsp, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterReq{
				AuthKey:  seesion.OpenID,
				AuthType: vars.UserAuthTypeSmallWX,
				Mobile:   mobile,
				Nickname: nickName,
			})
			if err != nil {
			   return nil,errors.Wrapf(errx.NewErrCode(errx.WXMINI_ERROR),"err :%v",err)
			}

			return &types.WXMiniAuthResp{
				AccessToken:  registerRsp.AccessToken,
				AccessExpire: registerRsp.AccessExpire,
				RefreshAfter: registerRsp.RefreshAfter,
			}, nil
		}
}

