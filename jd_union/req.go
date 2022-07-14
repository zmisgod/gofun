package jd_union

type JdUnionOpenGoodsMaterialQueryTopLevel struct {
	JdUnionOpenGoodsMaterialQueryResponse struct {
		Result string `json:"queryResult"`
		Code   string `json:"code"`
	} `json:"jd_union_open_goods_material_query_responce"`
}

type JdUnionOpenSellingPromotionGetTopLevel struct {
	JdUnionOpenSellingPromotionGetResponse struct {
		Result string `json:"getResult"`
		Code   string `json:"code"`
	} `json:"jd_union_open_promotion_common_get_responce"`
}

type JdUnionOpenGoodsQueryTopLevel struct {
	JdUnionOpenGoodsQueryResponse struct {
		Result string `json:"getResult"`
		Code   string `json:"code"`
	} `json:"jd_union_open_goods_query_response"`
}

type JdUnionOpenGoodsMaterialQueryResult struct {
	Code       int64   `json:"code"`
	Data       []Goods `json:"data"`
	Message    string  `json:"message"`
	RequestID  string  `json:"requestId"`
	TotalCount int64   `json:"totalCount"`
}

type JdUnionOpenSellingPromotionGetResult struct {
	Code int `json:"code"`
	Data struct {
		ClickURL string `json:"clickURL"`
		ShortURL string `json:"shortURL"`
	} `json:"data"`
	Message string `json:"message"`
}

type GoodCategoryInfo struct {
	Cid1     int    `json:"cid1"`
	Cid1Name string `json:"cid1Name"`
	Cid2     int    `json:"cid2"`
	Cid2Name string `json:"cid2Name"`
	Cid3     int    `json:"cid3"`
	Cid3Name string `json:"cid3Name"`
}

type GoodCommissionInfo struct {
	Commission          float64 `json:"commission"`
	CommissionShare     float32 `json:"commissionShare"`
	CouponCommission    float64 `json:"couponCommission"`
	PlusCommissionShare float32 `json:"plusCommissionShare"`
}

type GoodCouponInfo struct {
	CouponList []struct {
		BindType     int     `json:"bindType"`
		Discount     float32 `json:"discount"`
		GetEndTime   int64   `json:"getEndTime"`
		GetStartTime int64   `json:"getStartTime"`
		IsBest       int     `json:"isBest"`
		Link         string  `json:"link"`
		PlatformType int     `json:"platformType"`
		Quota        float32 `json:"quota"`
		UseEndTime   int64   `json:"useEndTime"`
		UseStartTime int64   `json:"useStartTime"`
	} `json:"couponList"`
}

type GoodImageInfo struct {
	ImageList []struct {
		URL string `json:"url"`
	} `json:"imageList"`
	WhiteImage string `json:"whiteImage"`
}

type GoodPriceInfo struct {
	LowestCouponPrice float64 `json:"lowestCouponPrice"`
	LowestPrice       float64 `json:"lowestPrice"`
	LowestPriceType   int     `json:"lowestPriceType"`
	Price             float64 `json:"price"`
}

type GoodPromotionInfo struct {
	ClickURL string `json:"clickURL"`
}

type GoodResourceInfo struct {
	EliteID   int    `json:"eliteId"`
	EliteName string `json:"eliteName"`
}

type GoodShopInfo struct {
	AfsFactorScoreRankGrade       string  `json:"afsFactorScoreRankGrade"`
	AfterServiceScore             string  `json:"afterServiceScore"`
	CommentFactorScoreRankGrade   string  `json:"commentFactorScoreRankGrade"`
	LogisticsFactorScoreRankGrade string  `json:"logisticsFactorScoreRankGrade"`
	LogisticsLvyueScore           string  `json:"logisticsLvyueScore"`
	ScoreRankRate                 string  `json:"scoreRankRate"`
	ShopID                        int     `json:"shopId"`
	ShopLabel                     string  `json:"shopLabel"`
	ShopLevel                     float64 `json:"shopLevel"`
	ShopName                      string  `json:"shopName"`
	UserEvaluateScore             string  `json:"userEvaluateScore"`
}

type GoodPinGouPrice struct {
	PingouPrice     float32 `json:"pingouPrice"`
	PingouTmCount   int     `json:"pingouTmCount"`
	PingouStartTime int64   `json:"pingouStartTime"`
	PingouEndTime   int64   `json:"pingouEndTime"`
}

type GoodVideoInfo struct {
	VideoList []GoodVideoInfoVideoList `json:"videoList"`
}

type GoodVideoInfoVideoList struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	ImageUrl  string `json:"imageUrl"`
	VideoType int    `json:"videoType"`
	PlayType  string `json:"playType"`
	Duration  int    `json:"duration"`
	PlayUrl   string `json:"playUrl"`
}

type Goods struct {
	BrandCode             string             `json:"brandCode"`
	BrandName             string             `json:"brandName"`
	CategoryInfo          GoodCategoryInfo   `json:"categoryInfo"`
	Comments              int                `json:"comments"`
	CommissionInfo        GoodCommissionInfo `json:"commissionInfo"`
	CouponInfo            GoodCouponInfo     `json:"couponInfo"`
	DeliveryType          int                `json:"deliveryType"`
	ForbidTypes           []int              `json:"forbidTypes"`
	GoodCommentsShare     float32            `json:"goodCommentsShare"`
	ImageInfo             GoodImageInfo      `json:"imageInfo"`
	InOrderCount30Days    int                `json:"inOrderCount30Days"`
	InOrderCount30DaysSku int                `json:"inOrderCount30DaysSku"`
	IsHot                 int                `json:"isHot"`
	MaterialURL           string             `json:"materialUrl"`
	Owner                 string             `json:"owner"`
	PinGouInfo            *GoodPinGouPrice   `json:"pinGouInfo"`
	PriceInfo             GoodPriceInfo      `json:"priceInfo"`
	PromotionInfo         GoodPromotionInfo  `json:"promotionInfo"`
	ResourceInfo          GoodResourceInfo   `json:"resourceInfo"`
	ShopInfo              GoodShopInfo       `json:"shopInfo"`
	SkuID                 int64              `json:"skuId"`
	SkuName               string             `json:"skuName"`
	Spuid                 int64              `json:"spuid"`
	VideoInfo             GoodVideoInfo      `json:"videoInfo"`
}
