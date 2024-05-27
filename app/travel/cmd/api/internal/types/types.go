// Code generated by goctl. DO NOT EDIT.
package types

type BusinessListReq struct {
	LastId             int64 `json:"lastId"`
	PageSize           int64 `json:"pageSize"`
	HomestayBusinessId int64 `json:"homestayBusinessId"`
}

type BusinessListResp struct {
	List []Homestay `json:"list"`
}

type CommentListReq struct {
	ByID     bool  `json:"By_id"`
	Id       int64 `json:"Id"`
	LastId   int64 `json:"lastId"`
	PageSize int64 `json:"pageSize"`
}

type CommentListResp struct {
	List []HomestayComment `json:"list"`
}

type GoodBossReq struct {
}

type GoodBossResp struct {
	List []HomestayBusinessBoss `json:"list"`
}

type GuessListReq struct {
}

type GuessListResp struct {
	List []Homestay `json:"list"`
}

type Homestay struct {
	Id                  int64   `json:"id"`
	Title               string  `json:"title"`
	SubTitle            string  `json:"subTitle"`
	Banner              string  `json:"banner"`
	Info                string  `json:"info"`
	PeopleNum           int64   `json:"peopleNum"`           //容纳人的数量
	HomestayBusinessId  int64   `json:"homestayBusinessId"`  //店铺id
	UserId              int64   `json:"userId"`              //房东id
	RowState            int64   `json:"rowState"`            //0:下架 1:上架
	RowType             int64   `json:"rowType"`             //售卖类型0：按房间出售 1:按人次出售
	FoodInfo            string  `json:"foodInfo"`            //餐食标准
	FoodPrice           float64 `json:"foodPrice"`           //餐食价格
	HomestayPrice       float64 `json:"homestayPrice"`       //民宿价格
	MarketHomestayPrice float64 `json:"marketHomestayPrice"` //民宿市场价格
}

type HomestayBusiness struct {
	Id        int64   `json:"id"`
	Title     string  `json:"title"` //店铺名称
	Info      string  `json:"info"`  //店铺介绍
	Tags      string  `json:"tags"`  //标签，多个用“,”分割
	Cover     string  `json:"cover"`
	Star      float64 `json:"star"`
	IsFav     int64   `json:"isFav"`     //是否收藏
	HeaderImg string  `json:"headerImg"` //店招门头图片
}

type HomestayBusinessBoss struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"userId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Info     string `json:"info"` //房东介绍
	Rank     int64  `json:"rank"` //排名
}

type HomestayBusinessListInfo struct {
	HomestayBusiness
	SellMonth     int64 `json:"sellMonth"`     //月销售
	PersonConsume int64 `json:"personConsume"` //个人消费
}

type HomestayBussinessDetailReq struct {
	Id int64 `json:"id"`
}

type HomestayBussinessDetailResp struct {
	Boss HomestayBusinessBoss `json:"boss"`
}

type HomestayBussinessListReq struct {
	LastId   int64 `json:"lastId"`
	PageSize int64 `json:"pageSize"`
}

type HomestayBussinessListResp struct {
	List []HomestayBusinessListInfo `json:"list"`
}

type HomestayComment struct {
	Id         int64   `json:"id"`         //评论id
	HomestayId int64   `json:"homestayId"` //民宿id
	Content    string  `json:"content"`    //评论
	Star       float64 `json:"star"`       //评论点赞
	UserId     int64   `json:"userId"`     //用户id
	Nickname   string  `json:"nickname"`   //用户名
	Avatar     string  `json:"avatar"`     //用户头像
}

type HomestayDetailReq struct {
	Id int64 `json:"id"`
}

type HomestayDetailResp struct {
	Homestay Homestay `json:"homestay"`
}

type HomestayListReq struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
}

type HomestayListResp struct {
	List []Homestay `json:"list"`
}
