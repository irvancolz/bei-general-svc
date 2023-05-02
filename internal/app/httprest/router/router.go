package router

import (
	"be-idx-tsg/internal/app/helper"
	Announcement "be-idx-tsg/internal/app/httprest/handler/announcement"
	Guidances "be-idx-tsg/internal/app/httprest/handler/guidances"

	middlewares "be-idx-tsg/internal/global"
	"os"

	global "be-idx-tsg/internal/global"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "pong",
			"gin_mode":  os.Getenv("GIN_MODE"),
			"http_port": os.Getenv("HTTP_PORT"),
		})
	})
	globalRepo := middlewares.NewRepositorys()
	announcement := Announcement.NewHandler()
	guidances := Guidances.NewGuidanceHandler()

	v3noauth := r.Group("/api")
	bukuPetujukBerkasPengaturan := global.BukuPetunjukBerkasPengaturan

	UploadFile := v3noauth.Group("").Use(globalRepo.Authentication(&bukuPetujukBerkasPengaturan))
	{
		UploadFile.POST("/upload-file", helper.UploadFile)
		UploadFile.DELETE("/delete-file", helper.DeleteFile)
		UploadFile.GET("/uploaded/:filename", helper.GetFile)
	}

	// WithoutCheckPermission := v3noauth.Group("").Use(globalRepo.Authentication())
	// {

	// }

	// announcementRoute := v3noauth.Group("").Use(globalRepo.Authentication(nil))
	announcementRoute := v3noauth.Group("").Use(globalRepo.Authentication(nil))
	{
		announcementRoute.GET("/get-all-announcement", announcement.GetAllAnnouncement) // used
		announcementRoute.POST("/create-announcement", announcement.Create)             // used
		announcementRoute.GET("/get-by-id-announcement", announcement.GetById)          // used
		announcementRoute.PUT("/update-announcement", announcement.Update)
		announcementRoute.DELETE("/delete-announcement", announcement.Delete)
		// announcementRoute.GET("/get-all-an-for-form", announcement.GetAllMin)
		// announcementRoute.GET("/get-an-by-code", announcement.GetByCode)
		// announcementRoute.POST("/get-an-by-id-and-type", announcement.GetByIDandType)
		announcementRoute.GET("/get-an-by-filter", announcement.GetAllANWithFilter)
		announcementRoute.GET("/get-an-by-search", announcement.GetAllANWithSearch)

	}
	guidancesRoute := v3noauth.Group("").Use(globalRepo.Authentication(&bukuPetujukBerkasPengaturan))
	{
		guidancesRoute.POST("/create-new-guidances", guidances.CreateNewGuidance)
		guidancesRoute.PUT("/update-guidances", guidances.UpdateExistingGuidance)
		guidancesRoute.POST("/create-new-files", guidances.CreateNewFiles)
		guidancesRoute.PUT("/update-files", guidances.UpdateExistingFiles)
		guidancesRoute.POST("/create-new-regulation", guidances.CreateNewRegulation)
		guidancesRoute.PUT("/update-regulation", guidances.UpdateExistingRegulation)
		guidancesRoute.GET("/get-all-guidance-file-or-regulation-by-type", guidances.GetAllGuidanceBasedOnType)
		guidancesRoute.DELETE("/delete-guidance-file-or-regulation", guidances.DeleteGuidances)
	}

	return r
}
