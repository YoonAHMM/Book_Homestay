package homestayComment

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/common/errx"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	var Builder squirrel.SelectBuilder
	
	if req.ByID{
		Builder =l.svcCtx.HomestayCommentModel.SelectBuilder().Where(squirrel.Eq{"homestay_id": req.Id})
	}else{
		Builder = l.svcCtx.HomestayCommentModel.SelectBuilder()
	}

	list, err := l.svcCtx.HomestayCommentModel.FindPageListByIdDESC(l.ctx, Builder, req.LastId, req.PageSize)
	if err != nil {
		return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
	}

	var resp_list []types.HomestayComment

	if len(list) > 0 {
		for _, item := range list {
			var HomestayComment types.HomestayComment
			_ = copier.Copy(&HomestayComment, item)

			resp_list = append(resp_list, HomestayComment)
		}
	}
	return &types.CommentListResp{
		List: resp_list,
	},nil
}
