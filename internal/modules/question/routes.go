package question

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/middleware"
	"github.com/topboyasante/pitstop/internal/modules/question/handler"
)

// RegisterRoutes registers all question-related routes
func RegisterRoutes(app fiber.Router, questionHandler *handler.QuestionHandler, answerHandler *handler.AnswerHandler) {
	SetupRoutes(app, questionHandler, answerHandler)
}

// SetupRoutes sets up all question-related routes
func SetupRoutes(app fiber.Router, questionHandler *handler.QuestionHandler, answerHandler *handler.AnswerHandler) {
	// Question routes
	questions := app.Group("/questions")

	// Public question routes
	questions.Get("/", questionHandler.GetAllQuestions)
	questions.Get("/tag", questionHandler.GetQuestionsByTag)
	questions.Get("/:id", questionHandler.GetQuestion)

	// Protected question routes (require authentication)
	protected := questions.Group("", middleware.JWTMiddleware(config.Get()))
	protected.Post("/", questionHandler.CreateQuestion)
	protected.Put("/:id", questionHandler.UpdateQuestion)
	protected.Delete("/:id", questionHandler.DeleteQuestion)

	// Answer routes for questions
	questionAnswers := questions.Group("/:question_id/answers")

	// Public answer routes
	questionAnswers.Get("/", answerHandler.GetAnswersByQuestionID)
	questionAnswers.Get("/:answer_id", answerHandler.GetAnswer)

	// Protected answer routes (require authentication)
	protectedAnswers := questionAnswers.Group("", middleware.JWTMiddleware(config.Get()))
	protectedAnswers.Post("/", answerHandler.CreateAnswer)
	protectedAnswers.Put("/:answer_id", answerHandler.UpdateAnswer)
	protectedAnswers.Delete("/:answer_id", answerHandler.DeleteAnswer)

	// Answer acceptance routes (require authentication)
	protectedAnswers.Post("/:answer_id/accept", answerHandler.AcceptAnswer)
	protectedAnswers.Post("/:answer_id/unaccept", answerHandler.UnacceptAnswer)
}