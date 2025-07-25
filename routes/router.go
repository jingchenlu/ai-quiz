package routes

import (
	"aiquiz/controllers"
	"aiquiz/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppDependencies interface {
	GetAuthController() *controllers.AuthController
	GetUserController() *controllers.UserController
	GetQuestionController() *controllers.QuestionController
	GetPaperController() *controllers.PaperController
	GetStatisticController() *controllers.StatisticController
	GetDB() *gorm.DB
}

func InitRouter(deps AppDependencies) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		// 创建控制器实例
		authController := deps.GetAuthController()
		userController := deps.GetUserController()
		questionController := deps.GetQuestionController()
		paperController := deps.GetPaperController()
		statisticController := deps.GetStatisticController()
		DB := deps.GetDB()

		// 认证相关路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
		}

		// 需要认证的路由
		authorized := api.Group("/", middlewares.JWTAuth())
		{
			// 用户相关路由
			users := authorized.Group("/users")
			{
				users.GET("/", userController.ListUsers)
				users.PUT("/:id", userController.UpdateUser)
				users.DELETE("/:id", userController.DeleteUser)
			}

			// 题目相关路由
			questions := authorized.Group("/questions")
			{
				questions.POST("/generate", questionController.GenerateQuestion)
				questions.POST("/confirm", questionController.ConfirmQuestions)
				questions.GET("/", questionController.ListQuestions)
				// 需要判断是否为该用户的题目，由于方法较少故未抽象为中间件
				questions.PUT("/:question_id", questionController.UpdateQuestion)
				questions.DELETE("/:question_id", questionController.DeleteQuestion)
			}
			// 试卷相关路由
			papers := authorized.Group("/papers")
			{
				papers.POST("/", paperController.CreatePaper)
				papers.GET("/", paperController.ListPapers)
				// 中间件: 校验用户是否有该试卷权限
				paperAuth := papers.Group("/:paper_id", middlewares.PaperAuth(DB))
				{
					paperAuth.GET("/", paperController.GetPaper)
					paperAuth.PUT("/", paperController.UpdatePaper)
					paperAuth.DELETE("/", paperController.DeletePaper)
					// 试卷题目相关
					paperQuestion := paperAuth.Group("/questions")
					{
						paperQuestion.POST("/", paperController.AddPaperQuestions)
						paperQuestion.DELETE("/:question_id", paperController.DeletePaperQuestions)
						paperQuestion.PUT("/order", paperController.UpdatePaperQuestionOrder)
					}
				}
			}
			// 统计相关路由
			statistics := authorized.Group("/statistics", middlewares.AdminMiddleware())
			{
				// 用户统计路由
				statistics.GET("/users/:user_id", statisticController.GetUserStatistics)
				// 整体统计路由
				statistics.GET("/overview", statisticController.GetSystemStatistics)
			}
		}

	}
	return r
}
