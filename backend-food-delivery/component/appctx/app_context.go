package appctx

import (
	cacheEngine "food_delivery/component/cache"
	"food_delivery/component/uploadprovider"
	"food_delivery/pubsub"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	UploadProvider() uploadprovider.UploadProvider
	GetSecretKey() string
	//NewTokenConfig() TokenConfig
	GetPubsub() pubsub.Pubsub
	GetCache() cacheEngine.Cache
}

type appCtx struct {
	db             *gorm.DB
	uploadProvider uploadprovider.UploadProvider
	secretKey      string
	pb             pubsub.Pubsub
	cache          cacheEngine.Cache
}

func NewAppContext(db *gorm.DB, uploadProvider uploadprovider.UploadProvider, secretKey string, pb pubsub.Pubsub, cache cacheEngine.Cache) *appCtx {
	return &appCtx{
		db:             db,
		uploadProvider: uploadProvider,
		secretKey:      secretKey,
		pb:             pb,
		cache:          cache,
	}
}

func (ctx *appCtx) GetMainDBConnection() *gorm.DB {
	return ctx.db
}
func (ctx *appCtx) UploadProvider() uploadprovider.UploadProvider {
	return ctx.uploadProvider
}

func (ctx *appCtx) GetSecretKey() string {
	return ctx.secretKey
}

func (ctx *appCtx) GetPubsub() pubsub.Pubsub {
	return ctx.pb
}

func (ctx *appCtx) GetCache() cacheEngine.Cache {
	return ctx.cache
}
