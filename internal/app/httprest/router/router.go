package router

import (
	Announcement "be-idx-tsg/internal/app/httprest/handler/announcement"
	Guidances "be-idx-tsg/internal/app/httprest/handler/guidances"
	JsonToXml "be-idx-tsg/internal/app/httprest/handler/jsontoxml"
	Pkp "be-idx-tsg/internal/app/httprest/handler/pkp"
	Topic "be-idx-tsg/internal/app/httprest/handler/topic"
	UploadFiles "be-idx-tsg/internal/app/httprest/handler/upload"

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
	globalRepo := global.NewRepositorys()
	announcement := Announcement.NewHandler()
	guidances := Guidances.NewGuidanceHandler()
	upload := UploadFiles.NewHandler()
	pkp := Pkp.NewHandler()
	jsonToXml := JsonToXml.NewHandler()
	topic := Topic.NewHandler()

	v3noauth := r.Group("/api")
	bukuPetujukBerkasPengaturan := global.BukuPetunjukBerkasPengaturan
	ParameterAdmin := global.ParameterAdmin

	UploadFile := v3noauth.Group("").Use(globalRepo.Authentication(&bukuPetujukBerkasPengaturan))
	{
		UploadFile.POST("/upload-form-file", upload.UploadForm)
		UploadFile.POST("/upload-admin-file", upload.UploadAdmin)
		UploadFile.POST("/upload-user-file", upload.UploadUser)
		UploadFile.POST("/upload-pkp-file", upload.UploadPkp)
		UploadFile.POST("/upload-report-file", upload.UploadReport)
		UploadFile.POST("/upload-guidances-files-regulation-file", upload.UploadGuidebook)

	}

	WithoutCheckPermission := v3noauth.Group("").Use(globalRepo.Authentication(nil))
	{
		WithoutCheckPermission.GET("/download-existing-file", upload.Download)
		WithoutCheckPermission.DELETE("/delete-existing-file", upload.Remove)
	}

	announcementRoute := v3noauth.Group("").Use(globalRepo.Authentication(nil))
	{
		announcementRoute.GET("/get-all-announcement", announcement.GetAllAnnouncement) // used
		announcementRoute.POST("/create-announcement", announcement.Create)             // used
		announcementRoute.GET("/get-by-id-announcement", announcement.GetById)          // used
		announcementRoute.PUT("/update-announcement", announcement.Update)
		announcementRoute.DELETE("/delete-announcement", announcement.Delete)
		announcementRoute.GET("/get-an-by-filter", announcement.GetAllANWithFilter)
		announcementRoute.POST("/get-an-by-search", announcement.GetAllANWithSearch)
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
		guidancesRoute.GET("/get-all-guidance-file-or-regulation", guidances.GetAllData)
		guidancesRoute.DELETE("/delete-guidance-file-or-regulation", guidances.DeleteGuidances)
	}
	pkpRoute := v3noauth.Group("").Use(globalRepo.Authentication(nil))
	{
		pkpRoute.GET("/get-all-pkp", pkp.GetAllPKuser)
		pkpRoute.POST("/create-pkp", pkp.CreatePKuser)
		pkpRoute.PUT("/update-pkp", pkp.UpdatePKuser)
		pkpRoute.DELETE("/delete-pkp", pkp.Delete)
		pkpRoute.GET("/get-pkp-by-filter", pkp.GetAllWithFilter)
		pkpRoute.GET("/get-pkp-by-search", pkp.GetAllWithSearch)
	}

	parameterAdminRoute := v3noauth.Group("").Use(globalRepo.Authentication(&ParameterAdmin))
	{
		parameterAdminRoute.POST("/upload-parameter-admin-file", upload.UploadParameterAdminFile)
		parameterAdminRoute.POST("/upload-parameter-admin-image", upload.UploadParameterAdminImage)
	}
	jsonToXmlRoute := v3noauth.Group("").Use(globalRepo.Authentication(nil))
	{
		jsonToXmlRoute.POST("/to-xml", jsonToXml.ToXml)
	}

	WithoutToken := v3noauth.Group("")
	{
		WithoutToken.GET("/download-existing-file-without-token", upload.Download)
	}

	topicRoute := v3noauth.Group("").Use(globalRepo.Authentication(nil))
	{
		topicRoute.GET("/get-all-topic", topic.GetAll)
		topicRoute.GET("/get-by-id-topic", topic.GetById)
		topicRoute.POST("/create-topic", topic.CreateTopicWithMessage)
		topicRoute.PUT("/update-handler", topic.UpdateHandler)
		topicRoute.POST("/create-message", topic.CreateMessage)
		topicRoute.DELETE("/delete-topic", topic.DeleteTopic)
		topicRoute.POST("/archive-topic", topic.ArchiveTopicToFAQ)
	}

	return r
}
