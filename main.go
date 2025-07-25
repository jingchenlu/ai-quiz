package main

import (
	"aiquiz/config"
	"aiquiz/controllers"
	"aiquiz/dao"
	"aiquiz/migrations"
	"aiquiz/routes"
	"aiquiz/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type AppDependencies struct {
	DB          *gorm.DB
	UserDAO     *dao.UserDao
	QuestionDAO *dao.QuestionDao
	PaperDAO    *dao.PaperDao
	statsDAO    *dao.UserStatisticsDao
	systemDAO   *dao.SystemStatisticsDao

	UserService     *services.UserService
	QuestionService *services.QuestionService
	PaperService    *services.PaperService
	statsService    *services.StatisticsService

	AuthController      *controllers.AuthController
	UserController      *controllers.UserController
	QuestionController  *controllers.QuestionController
	PaperController     *controllers.PaperController
	StatisticController *controllers.StatisticController
}

// GetAuthController 获取认证控制器
func (d *AppDependencies) GetAuthController() *controllers.AuthController {
	if d.AuthController == nil {
		d.AuthController = controllers.NewAuthController(d.UserService)
	}
	return d.AuthController
}
func (d *AppDependencies) GetUserController() *controllers.UserController {
	if d.UserController == nil {
		d.UserController = controllers.NewUserController(d.UserService)
	}
	return d.UserController
}
func (d *AppDependencies) GetQuestionController() *controllers.QuestionController {
	if d.QuestionController == nil {
		d.QuestionController = controllers.NewQuestionController(d.QuestionService)
	}
	return d.QuestionController
}
func (d *AppDependencies) GetPaperController() *controllers.PaperController {
	if d.PaperController == nil {
		d.PaperController = controllers.NewPaperController(d.PaperService)
	}
	return d.PaperController
}
func (d *AppDependencies) GetStatisticController() *controllers.StatisticController {
	if d.StatisticController == nil {
		d.StatisticController = controllers.NewStatisticController(d.statsService)
	}
	return d.StatisticController
}

func (d *AppDependencies) GetDB() *gorm.DB {
	return d.DB
}

// 程序入口
func main() {
	// 获取配置
	appConfig := config.GetConfig(true)

	// 设置Gin模式
	gin.SetMode(appConfig.Mode)

	// 初始化数据库连接
	db := migrations.InitDB(appConfig.DBPath)
	log.Println("数据库初始化完成")

	// 初始化依赖
	deps := initDependencies(db)

	// 设置路由
	router := routes.InitRouter(deps)

	log.Printf("服务器启动于 %s 端口，运行模式: %s\n", appConfig.ServerPort, appConfig.Mode)

	for language := range appConfig.SupportedLanguages {
		log.Printf("支持语言: %s\n", language)
	}

	// 启动服务器
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

}

// 初始化依赖
func initDependencies(db *gorm.DB) *AppDependencies {
	// 初始化DAO
	userDAO := dao.NewUserDAO(db)
	questionDao := dao.NewQuestionDAO(db)
	paperDao := dao.NewPaperDAO(db)
	statsDao := dao.NewUserStatisticsDao(db)
	systemStatisticsDao := dao.NewSystemStatisticsDao(db)

	// 初始化服务
	userService := services.NewUserService(userDAO, questionDao, paperDao)
	questionService := services.NewQuestionService(questionDao)
	paperService := services.NewPaperService(paperDao, questionDao)
	statsService := services.NewStatisticService(userDAO, statsDao, systemStatisticsDao)

	return &AppDependencies{
		DB:              db,
		UserDAO:         userDAO,
		QuestionDAO:     questionDao,
		PaperDAO:        paperDao,
		statsDAO:        statsDao,
		UserService:     userService,
		QuestionService: questionService,
		PaperService:    paperService,
		statsService:    statsService,
	}
}
