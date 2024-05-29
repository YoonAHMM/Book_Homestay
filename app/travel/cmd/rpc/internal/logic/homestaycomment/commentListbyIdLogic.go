package homestaycommentlogic

import (
	"context"

	"Book_Homestay/app/travel/cmd/rpc/internal/svc"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	"Book_Homestay/common/errx"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListbyIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentListbyIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListbyIdLogic {
	return &CommentListbyIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentListbyIdLogic) CommentListbyId(in *pb.CommentListbyIdReq) (*pb.CommentListbyIdResp, error) {
	var Builder squirrel.SelectBuilder
	
	Builder =l.svcCtx.HomestayCommentModel.SelectBuilder().Where(squirrel.Eq{"homestay_id": in.Id})
	

	list, err := l.svcCtx.HomestayCommentModel.FindPageListByIdDESC(l.ctx, Builder, in.Lastid, in.Pagesize)
	if err != nil {
		return nil, errors.Wrapf(errx.NewErrCode(errx.DB_ERROR), "err : %v , in : %+v", err, in)
	}

	var resp_list []*pb.HomestayComment

	if len(list) > 0 {
		for _, item := range list {
			var HomestayComment pb.HomestayComment
			_ = copier.Copy(&HomestayComment, item)

			resp_list = append(resp_list, &HomestayComment)
		}
	}

	return &pb.CommentListbyIdResp{
		CommentList: resp_list,
	}, nil
}
