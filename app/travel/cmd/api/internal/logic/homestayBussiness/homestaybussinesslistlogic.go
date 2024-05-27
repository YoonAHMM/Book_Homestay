package homestayBussiness

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/common/errx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBussinessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomestayBussinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBussinessListLogic {
	return &HomestayBussinessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayBussinessListLogic) HomestayBussinessList(req *types.HomestayBussinessListReq) (resp *types.HomestayBussinessListResp, err error) {
	Builder := l.svcCtx.HomestayBusinessModel.SelectBuilder()
	list, err := l.svcCtx.HomestayBusinessModel.FindPageListByIdDESC(l.ctx, Builder, req.LastId, req.PageSize)
	if err != nil {
		return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
	}

	var resp_list []types.HomestayBusinessListInfo
	if len(list) > 0 {
		for _, item := range list {
			var HomestayBusinessListInfo types.HomestayBusinessListInfo
			_ = copier.Copy(&HomestayBusinessListInfo, item)

			resp_list = append(resp_list, HomestayBusinessListInfo)
		}
	}

	return &types.HomestayBussinessListResp{
		List: resp_list,
	}, nil

}
