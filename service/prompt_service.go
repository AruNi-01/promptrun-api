package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"promptrun-api/common/errs"
	"promptrun-api/model"
	"promptrun-api/third_party"
	"promptrun-api/utils"
	"strconv"
	"time"
)

// SortBy 前端传入的排序方式
const (
	SortByHot        = "hot"         // 按照热度
	SortByTime       = "time"        // 按照时间
	SortBySellAmount = "sell_amount" // 按照销量
)

var (
	LatestBrowseAmountUpdateTime time.Time // 最近一次更新提示词浏览量的时间，防止短时间内多次更新
)

const (
	BrowseAmountUpdateInterval = 1 * time.Second // 更新提示词浏览量的时间间隔，1s
)

// PromptListReq 获取提示词列表请求
type PromptListReq struct {
	Paginate      *utils.Page `json:"paginate"`
	ModelId       *int        `json:"modelId,omitempty"`
	CategoryTypes []int       `json:"categoryTypes,omitempty"`
	SortBy        string      `json:"sortBy,omitempty"`

	PublishStatus []int `json:"publishStatus,omitempty"`
	AuditStatus   []int `json:"auditStatus,omitempty"`

	UserId *int `json:"userId,omitempty"` // 用户 ID，用于获取某个用户（seller）发布的提示词
}

// PromptMasterImgListReq 获取主图列表请求
type PromptMasterImgListReq struct {
	PromptIds []int `json:"promptIds" binding:"required"`
}

// PromptListByBuyerIdReq 根据买家 id 获取提示词列表请求
type PromptListByBuyerIdReq struct {
	Paginate *utils.Page `json:"paginate"`
	BuyerId  int         `json:"buyerId"`
}

// PromptListBySellerIdReq 根据卖家 id 获取提示词列表请求
type PromptListBySellerIdReq struct {
	Paginate *utils.Page `json:"paginate"`
	SellerId int         `json:"sellerId"`
}

// PromptPublishReq 发布提示词请求
type PromptPublishReq struct {
	UserId             int      `json:"userId"`
	PromptTitle        string   `json:"promptTitle"`
	PromptModelId      int      `json:"promptModelId"`
	PromptCategoryType int      `json:"promptCategoryType"`
	PromptIntro        string   `json:"promptIntro"`
	PromptContent      string   `json:"promptContent"`
	UseSuggestion      string   `json:"useSuggestion"`
	InputExample       string   `json:"inputExample"`
	OutputExample      string   `json:"outputExample"`
	MasterImgBase64    string   `json:"masterImgBase64"`
	ImgBase64List      []string `json:"imgBase64List"`
	PromptPrice        float64  `json:"promptPrice"`
}

// PromptListResp 获取提示词列表响应
type PromptListResp struct {
	Prompts []model.Prompt `json:"prompts"`
	Rows    int            `json:"rows"`
}

type PromptDetailResp struct {
	Prompt        model.Prompt      `json:"prompt"`
	Seller        model.Seller      `json:"seller"`
	SellerUser    model.User        `json:"sellerUser"`
	Model         model.Model       `json:"model"`
	PromptImgList []model.PromptImg `json:"promptImgList"`
}

// PromptList 获取提示词列表
func (r *PromptListReq) PromptList(c *gin.Context) ([]model.Prompt, *errs.Errs) {
	var prompts []model.Prompt
	query := model.DB.Model(model.Prompt{})

	if r.ModelId != nil {
		query = query.Where("model_id = ?", *r.ModelId)
	}
	if r.CategoryTypes != nil {
		query = query.Where("category_type IN (?)", r.CategoryTypes)
	}
	switch r.SortBy {
	case SortByHot:
		query = query.Order("score DESC")
	case SortBySellAmount:
		query = query.Order("sell_amount DESC")
	case SortByTime:
		query = query.Order("create_time DESC")
	default:
		query = query.Order("create_time DESC")
	}

	if r.PublishStatus != nil {
		query = query.Where("publish_status IN (?)", r.PublishStatus)
	}
	if r.AuditStatus != nil {
		query = query.Where("audit_status IN (?)", r.AuditStatus)
	}

	if r.UserId != nil {
		seller, e := FindSellerByUserId(c, *r.UserId)
		if e != nil {
			return nil, e
		}
		query = query.Where("seller_id = ?", seller.Id)
	}

	// 注意：Count 要放在分页查询之前，否则会导致 count 为空
	query.Count(&r.Paginate.Rows)

	if query.Scopes(
		utils.Paginate(r.Paginate)).
		Find(&prompts).
		Error != nil {
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示词列表失败"))
	}

	return prompts, nil
}

// FindPromptFullInfoById 根据 ID 查找提示词的详细信息（提示词详情页使用借口）
func FindPromptFullInfoById(c *gin.Context, promptId int) (PromptDetailResp, *errs.Errs) {
	// TODO: 暂时用此方法更新提示词浏览量，后续考虑使用 Redis
	go func(promptId int) {
		updatePromptBrowseAmountAsync(c, promptId)
	}(promptId)

	prompt, e := FindPromptById(c, promptId)
	if e != nil {
		return PromptDetailResp{}, e
	}
	seller, e := FindSellerById(c, prompt.SellerId)
	if e != nil {
		return PromptDetailResp{}, e
	}
	sellerUser, e := FindUserById(c, seller.UserId)
	if e != nil {
		return PromptDetailResp{}, e
	}
	promptModel, e := FindModelById(c, prompt.ModelId)
	if e != nil {
		return PromptDetailResp{}, e
	}
	promptImgList, e := FindPromptImgListByPromptId(c, promptId)
	if e != nil {
		return PromptDetailResp{}, e
	}
	return PromptDetailResp{
		Prompt:        prompt,
		Seller:        seller,
		SellerUser:    sellerUser,
		Model:         promptModel,
		PromptImgList: promptImgList,
	}, nil

}

func updatePromptBrowseAmountAsync(c *gin.Context, promptId int) {
	if LatestBrowseAmountUpdateTime.Add(BrowseAmountUpdateInterval).Before(time.Now()) {
		LatestBrowseAmountUpdateTime = time.Now()
		if err := model.DB.Model(model.Prompt{}).Where("id = ?", promptId).UpdateColumn("browse_amount", gorm.Expr("browse_amount + 1")).Error; err != nil {
			utils.Log().Error(c.FullPath(), "异步更新提示词浏览量失败")
		}
	}
}

// FindPromptById 根据 ID 查找提示词
func FindPromptById(c *gin.Context, promptId int) (model.Prompt, *errs.Errs) {
	var prompt model.Prompt
	if model.DB.First(&prompt, promptId).RecordNotFound() {
		utils.Log().Error(c.FullPath(), "提示词不存在")
		return model.Prompt{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("提示词不存在"))
	}
	return prompt, nil
}

func FindPromptImgListByPromptId(c *gin.Context, id int) ([]model.PromptImg, *errs.Errs) {
	var promptImgList []model.PromptImg
	if model.DB.Where("prompt_id = ?", id).Find(&promptImgList).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取提示词图片列表失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示词图片列表失败"))
	}
	return promptImgList, nil
}

func FindPromptMasterImgByPromptId(promptId int) (model.PromptImg, *errs.Errs) {
	var promptImg model.PromptImg
	if model.DB.Where("prompt_id = ?", promptId).Where("is_master = ?", 1).First(&promptImg).RecordNotFound() {
		return model.PromptImg{}, errs.NewErrs(errs.ErrRecordNotFound, errors.New("未找到主图"))
	}
	return promptImg, nil
}

func (r *PromptMasterImgListReq) FindMasterImgListByPromptIds(c *gin.Context) ([]model.PromptImg, *errs.Errs) {
	var promptImgList []model.PromptImg
	if model.DB.Where("prompt_id IN (?)", r.PromptIds).Where("is_master = ?", 1).Find(&promptImgList).Error != nil {
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取主图列表失败"))
	}
	return promptImgList, nil
}

// FindListByBuyerId 根据买家 ID 获取该买家买入的提示词列表
func (r *PromptListByBuyerIdReq) FindListByBuyerId(c *gin.Context) ([]model.Prompt, *errs.Errs) {
	orderList, err := FindOrderListByBuyerId(c, r.BuyerId)
	if err != nil {
		return nil, err
	}

	var promptIds = func() []int {
		var ids []int
		for _, order := range orderList {
			ids = append(ids, order.PromptId)
		}
		return ids
	}()

	var prompts []model.Prompt
	query := model.DB.
		Model(model.Prompt{}).
		Where("id in (?)", promptIds).
		Order("create_time DESC").
		Count(&r.Paginate.Rows)

	if query.Scopes(utils.Paginate(r.Paginate)).Find(&prompts).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取提示词列表失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示词列表失败"))
	}
	return prompts, nil
}

// FindListBySellerId 根据卖家 ID 获取该卖家发布的提示词列表
func (r *PromptListBySellerIdReq) FindListBySellerId(c *gin.Context) ([]model.Prompt, *errs.Errs) {
	var prompts []model.Prompt
	query := model.DB.
		Model(model.Prompt{}).
		Where("seller_id = ?", r.SellerId).
		Where("publish_status = ?", model.PublishStatusOn).
		Where("audit_status = ?", model.AuditStatusPass).
		Order("create_time DESC").
		Count(&r.Paginate.Rows)

	if query.Scopes(utils.Paginate(r.Paginate)).Find(&prompts).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取提示词列表失败")
		return nil, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示词列表失败"))
	}
	return prompts, nil
}

func FindPromptCountBySellerId(c *gin.Context, sellerId int) (int, *errs.Errs) {
	var count int
	if model.DB.
		Model(model.Prompt{}).
		Where("seller_id = ?", sellerId).
		Where("publish_status = ?", model.PublishStatusOn).
		Where("audit_status = ?", model.AuditStatusPass).
		Count(&count).
		Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取提示词数量失败")
		return 0, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取提示词数量失败"))
	}
	return count, nil
}

func UpdatePromptBrowseAmountById(c *gin.Context, promptId int) (bool, *errs.Errs) {
	if err := model.DB.Model(model.Prompt{}).Where("id = ?", promptId).UpdateColumn("browse_amount", gorm.Expr("browse_amount + 1")).Error; err != nil {
		utils.Log().Error(c.FullPath(), "更新提示词浏览量失败")
		return false, errs.NewErrs(errs.ErrDBError, errors.New("更新提示词浏览量失败"))
	}
	return true, nil
}
func (r *PromptPublishReq) PromptPublish(c *gin.Context) (bool, *errs.Errs) {
	var promptModel model.Model
	if model.DB.First(&promptModel, r.PromptModelId).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取模型失败")
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取模型失败"))
	}
	var seller model.Seller
	if model.DB.Where("user_id = ?", r.UserId).First(&seller).Error != nil {
		utils.Log().Error(c.FullPath(), "DB 获取卖家失败")
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 获取卖家失败"))
	}

	prompt := model.Prompt{
		SellerId:      seller.Id,
		Title:         r.PromptTitle,
		ModelId:       r.PromptModelId,
		CategoryType:  r.PromptCategoryType,
		Intro:         r.PromptIntro,
		InputExample:  r.InputExample,
		OutputExample: r.OutputExample,
		Rating:        model.Ratting5, // 默认评分为 5.0，后续根据有买家评分时再根据评分计算
		Price:         r.PromptPrice,
		AuditStatus:   model.AuditStatusPass,
		PublishStatus: model.PublishStatusOn,
		CreateTime:    time.Now(),
	}
	if err := model.DB.Create(&prompt).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 创建提示词失败")
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 创建提示词失败"))
	}

	switch promptModel.MediaType {
	case model.ModelMediaTypeText:
		return r.handleTextPromptPublish(c, prompt.Id)
	case model.ModelMediaTypeImage:
		return r.handleImagePromptPublish(c, prompt.Id)
	case model.ModelMediaTypeVideo:
		return r.handleVideoPromptPublish(c, prompt.Id)
	default:
		utils.Log().Error(c.FullPath(), "模型类型错误")
		return false, errs.NewErrs(errs.ErrParam, errors.New("模型类型错误"))
	}
}

func (r *PromptPublishReq) handleTextPromptPublish(c *gin.Context, promptId int) (bool, *errs.Errs) {
	promptImg := model.PromptImg{
		PromptId: promptId,
		ImgUrl: func(promptId int) string {
			objectName := third_party.OSSPrefixPromptImg + strconv.Itoa(promptId) + "-" + time.Now().Format("2006-01-02_150405")
			headerUrl, err := third_party.UploadBase64ImgToOSS(objectName, r.MasterImgBase64)
			if err != nil {
				utils.Log().Error(c.FullPath(), "OSS 上传 Banner 图片失败，errMsg: %s", err.Error())
			}
			return headerUrl
		}(promptId),
		IsMaster: model.PromptImgIsMaster,
	}
	if err := model.DB.Create(&promptImg).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 创建提示词图片失败")
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 创建提示词图片失败"))
	}

	promptDetail := model.PromptDetail{
		PromptId:      promptId,
		MediaType:     model.ModelMediaTypeText,
		Content:       r.PromptContent,
		UseSuggestion: r.UseSuggestion,
		CreateTime:    time.Now(),
	}
	if err := model.DB.Create(&promptDetail).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 创建提示词详情失败")
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 创建提示词详情失败"))
	}
	return true, nil
}

func (r *PromptPublishReq) handleImagePromptPublish(c *gin.Context, promptId int) (bool, *errs.Errs) {
	promptMasterImg := model.PromptImg{
		PromptId: promptId,
		ImgUrl: func(promptId int) string {
			objectName := third_party.OSSPrefixPromptImg + strconv.Itoa(promptId) + "-" + time.Now().Format("2006-01-02_150405")
			headerUrl, err := third_party.UploadBase64ImgToOSS(objectName, r.MasterImgBase64)
			if err != nil {
				utils.Log().Error(c.FullPath(), "OSS 上传 Banner 图片失败，errMsg: %s", err.Error())
			}
			return headerUrl
		}(promptId),
		IsMaster: model.PromptImgIsMaster,
	}
	if err := model.DB.Create(&promptMasterImg).Error; err != nil {
		utils.Log().Error(c.FullPath(), "DB 创建提示词图片失败")
		return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 创建提示词图片失败"))
	}

	for _, imgBase64 := range r.ImgBase64List {
		promptImg := model.PromptImg{
			PromptId: promptId,
			ImgUrl: func(promptId int) string {
				objectName := third_party.OSSPrefixPromptImg + strconv.Itoa(promptId) + "-" + time.Now().Format("2006-01-02_150405")
				headerUrl, err := third_party.UploadBase64ImgToOSS(objectName, imgBase64)
				if err != nil {
					utils.Log().Error(c.FullPath(), "OSS 上传提示词图片失败，errMsg: %s", err.Error())
				}
				return headerUrl
			}(promptId),
			IsMaster: model.PromptImgNotMaster,
		}
		if err := model.DB.Create(&promptImg).Error; err != nil {
			utils.Log().Error(c.FullPath(), "DB 创建提示词图片失败")
			return false, errs.NewErrs(errs.ErrDBError, errors.New("DB 创建提示词图片失败"))
		}
	}
	return true, nil
}

func (r *PromptPublishReq) handleVideoPromptPublish(c *gin.Context, promptId int) (bool, *errs.Errs) {
	// TODO: 实现视频提示词发布
	return false, errs.NewErrs(errs.ErrParam, errors.New("暂不支持视频提示词发布"))
}

func UpdatePromptPublishStatusById(c *gin.Context, promptId int, status int) (bool, *errs.Errs) {
	if err := model.DB.Model(model.Prompt{}).Where("id = ?", promptId).UpdateColumn("publish_status", status).Error; err != nil {
		utils.Log().Error(c.FullPath(), "更新提示词发布状态失败")
		return false, errs.NewErrs(errs.ErrDBError, errors.New("更新提示词发布状态失败"))
	}
	return true, nil
}
