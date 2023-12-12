package api

import (
	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/golanguzb70/simple-post-app/api/docs" // docs
	v1 "github.com/golanguzb70/simple-post-app/api/handlers/v1"
	"github.com/golanguzb70/simple-post-app/api/middleware"
	t "github.com/golanguzb70/simple-post-app/api/tokens"
	"github.com/golanguzb70/simple-post-app/config"
	"github.com/golanguzb70/simple-post-app/pkg/logger"
	"github.com/golanguzb70/simple-post-app/storage"
	"github.com/golanguzb70/simple-post-app/storage/redisrepo"
	"github.com/gomodule/redigo/redis"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Option ...
type Option struct {
	Conf       config.Config
	Logger     *logger.Logger
	Postgres   storage.StorageI
	JWTHandler t.JWTHandler
	Redis      redisrepo.InMemoryStorageI
}

// New ...
// @title           Monolithic project API Endpoints
// @version         1.0
// @description     Here QA can test and frontend or mobile developers can get information of API endpoints.
// @BasePath  /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(log *logger.Logger, cfg config.Config, strg storage.StorageI) *gin.Engine {
	casbinEnforcer, err := casbin.NewEnforcer(cfg.AuthConfigPath, cfg.CSVFilePath)
	if err != nil {
		log.Error("casbin enforcer error", err)
		panic(err)
	}
	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		log.Error("casbin error load policy", err)
		panic(err)
	}

	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch3", util.KeyMatch3)

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.RedisHost+":"+cfg.RedisPort)
		},
	}

	jwtHandler := t.JWTHandler{
		SigninKey: cfg.SignInKey,
		Log:       log,
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	h := v1.New(&v1.HandlerV1Config{
		Logger:     log,
		Cfg:        cfg,
		Postgres:   strg,
		JWTHandler: jwtHandler,
		Redis:      redisrepo.NewRedisRepo(pool),
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.NewAuth(casbinEnforcer, jwtHandler, cfg))

	api := router.Group("/v1")

	user := api.Group("/user")
	user.GET("/check/:email", h.UserCheck)
	user.GET("/otp", h.OtpCheck)
	user.POST("", h.UserRegister)
	user.POST("/login", h.LoginUser)
	user.GET("/forgot-password/:user_name_or_email", h.UserForgotPassword)
	user.POST("/forgot-password/verify", h.UserForgotPasswordVerify)
	user.GET("/profile", h.UserGet)
	user.PUT("", h.UserUpdate)
	user.DELETE("", h.UserDelete)

	post := api.Group("/post")
	post.POST("", h.PostCreate)
	post.GET("/:id", h.PostGet)
	post.GET("/list", h.PostFind)
	post.PUT("", h.PostUpdate)
	post.DELETE(":id", h.PostDelete)

	media := api.Group("/media")
	api.Static("/media", "./media")
	media.POST("/photo", h.UploadMedia)

	comment := api.Group("/comment")
	comment.POST("", h.CommentCreate)
	comment.GET("/:id", h.CommentGet)
	comment.GET("/list", h.CommentFind)
	comment.PUT("", h.CommentUpdate)
	comment.DELETE(":id", h.CommentDelete)

	// Don't delete this line, it is used to modify the file automatically

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
