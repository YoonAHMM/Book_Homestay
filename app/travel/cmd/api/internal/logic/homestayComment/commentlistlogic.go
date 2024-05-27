package homestayComment

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	"Book_Homestay/common/calculate"

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
	var Comment_list []*pb.HomestayComment

	if req.ByID{
		Comment_resp,err:=l.svcCtx.Comment_TravelRpc.CommentListbyId(l.ctx,&pb.CommentListbyIdReq{
			Id: req.Id,
			Lastid: req.LastId,
			Pagesize: req.PageSize,
		})
		if err!=nil{
			return nil,err
		}
		Comment_list=Comment_resp.CommentList
	}else{
		Comment_resp,err:=l.svcCtx.Comment_TravelRpc.CommentList(l.ctx,&pb.CommentListReq{
			Lastid: req.LastId,
			Pagesize: req.PageSize,
		})
		if err!=nil{
			return nil,err
		}
		Comment_list=Comment_resp.CommentList
	}

	var resq_list []types.HomestayComment

	for _,comment :=range Comment_list{
		
		var Comment types.HomestayComment
		_ = copier.Copy(&Comment, comment)

		Comment.Star=calculate.Ge2Qian(comment.Star)

		resq_list = append(resq_list, Comment)
	}
	
	return &types.CommentListResp{
		List: resq_list,
	},nil
}
