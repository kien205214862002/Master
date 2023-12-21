package main

import (
	"fmt"
	"food_delivery/common"
	"food_delivery/component/appctx"
	cacheEngine "food_delivery/component/cache"
	"food_delivery/component/uploadprovider"
	"food_delivery/config"
	"food_delivery/middleware"
	"food_delivery/modules/upload/uploadtransport/ginupload"
	"food_delivery/modules/user/usertransport/ginuser"
	"food_delivery/pubsub/pblocal"
	"food_delivery/routes/restaurantroute"
	"food_delivery/skio"
	"food_delivery/subscriber"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	//refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	cfg := config.GetConfig()
	log.Println(cfg.Environment)

	s3BucketName := cfg.S3BucketName
	s3Region := cfg.S3Region
	s3APIKey := cfg.S3APIKey
	s3SecretKey := cfg.S3SecretKey
	s3Domain := cfg.S3Domain
	secretKey := cfg.SystemSecret

	s3Provider := uploadprovider.NewS3Provider(s3BucketName, s3Region, s3APIKey, s3SecretKey, s3Domain)

	db, err := gorm.Open(mysql.Open(cfg.DatabaseURI), &gorm.Config{})

	cache := cacheEngine.New(cacheEngine.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})

	fmt.Println(db, err)
	db = db.Debug()

	if err != nil {
		log.Fatalln(err)
	}

	if err := runService(db, s3Provider, secretKey, cache); err != nil {
		log.Fatalln(err)
	}
}

func runService(db *gorm.DB, provider uploadprovider.UploadProvider, secretKey string, cache cacheEngine.Cache) error {
	appCtx := appctx.NewAppContext(db, provider, secretKey, pblocal.NewPubsub(), cache)

	r := gin.Default()

	rtEngine := skio.NewEngine()

	if err := rtEngine.Run(appCtx, r); err != nil {
		log.Fatalln(err)
	}

	//deprecated
	//subscriber.Setup(appCtx)

	// use this line as an alternative for Setup
	if err := subscriber.NewEngine(appCtx, rtEngine).Start(); err != nil {
		log.Fatalln()
	}

	r.Use(middleware.Recover(appCtx))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.StaticFile("/demo/", "./demo.html")

	v1 := r.Group("/v1")
	v1.POST("/upload", ginupload.Upload(appCtx))
	//v1.GET("/presigned-upload-url", func(C *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{"data": s3Provider.GetUploadPresignedUrl(c.Request.Context())})
	//})

	v1.POST("/register", ginuser.Register(appCtx))
	v1.POST("/login", ginuser.Login(appCtx))
	v1.GET("/profile", middleware.RequireAuth(appCtx), ginuser.GetProfile(appCtx))

	v1.GET("/encode-uid", func(c *gin.Context) {
		type reqData struct {
			DbType int `form:"type"`
			RealId int `form:"id"`
		}

		var d reqData
		c.ShouldBind(&d)

		c.JSON(http.StatusOK, gin.H{"id": common.NewUID(uint32(d.RealId), d.DbType, 1)})
	})
	restaurantroute.Routes(v1, appCtx)

	admin := v1.Group(
		"/admin",
		middleware.RequireAuth(appCtx),
		middleware.RequireRoles(appCtx, "admin"),
	)
	{
		admin.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, common.SimpleSucessResponse("ok"))
		})
	}

	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
