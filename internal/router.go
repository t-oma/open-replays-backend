package internal

// func New() *gin.Engine {
// 	r := gin.Default()
// 	r.Use(cors.Default())
//
// 	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Hello, world!"}) })
//
// 	r.GET("/watch/:filename", func(c *gin.Context) {
// 		filename := c.Param("filename")
//
// 		files, _ := filepath.Glob("uploads/videos/" + filename)
// 		if len(files) == 0 {
// 			c.JSON(404, domain.APIResponse{
// 				Success: false,
// 				Code:    404,
// 				Error:   "File not found",
// 			})
// 			return
// 		}
//
// 		c.File(files[0])
// 	})
//
// 	r.GET("/preview/:filename", func(c *gin.Context) {
// 		filename := c.Param("filename")
//
// 		files, _ := filepath.Glob("uploads/previews/" + filename)
// 		if len(files) == 0 {
// 			c.JSON(404, domain.APIResponse{
// 				Success: false,
// 				Code:    404,
// 				Error:   "Preview not found",
// 			})
// 			return
// 		}
//
// 		c.File(files[0])
// 	})
//
// 	r.GET("/api/v1/videos", func(c *gin.Context) {
// 		files, err := filepath.Glob("uploads/videos/*.mp4")
// 		if err != nil {
// 			c.JSON(500, domain.APIResponse{
// 				Success: false,
// 				Code:    500,
// 				Error:   "Failed to read videos",
// 			})
// 			return
// 		}
//
// 		videos := []sqlc.Video{}
// 		for _, file := range files {
// 			var uploadTime time.Time
// 			if stat, err := os.Stat(file); err == nil {
// 				uploadTime = stat.ModTime()
// 			}
//
// 			replay := sqlc.Video{
// 				Title:      stringutil.TrimExt(filepath.Base(file)),
// 				Filename:   filepath.Base(file),
// 				UploadedAt: uploadTime,
// 			}
//
// 			videos = append(videos, replay)
// 		}
//
// 		sort.Slice(videos, func(i, j int) bool {
// 			return videos[i].UploadedAt.After(videos[j].UploadedAt)
// 		})
//
// 		c.JSON(200, domain.APIResponse{
// 			Success: true,
// 			Data:    gin.H{"videos": videos},
// 		})
// 	})
//
// 	r.POST("/api/v1/videos/upload", func(c *gin.Context) {
// 		file, err := c.FormFile("video")
// 		if err != nil {
// 			c.JSON(400, domain.APIResponse{
// 				Success: false,
// 				Code:    400,
// 				Error:   "No file uploaded",
// 			})
// 			return
// 		}
//
// 		ext := filepath.Ext(file.Filename)
// 		if ext != ".mp4" && ext != ".webm" && ext != ".mov" {
// 			c.JSON(400, domain.APIResponse{
// 				Success: false,
// 				Code:    400,
// 				Error:   "Invalid file type",
// 			})
// 			return
// 		}
//
// 		if file.Size > 100*MEGABYTE { // 100MB
// 			c.JSON(400, domain.APIResponse{
// 				Success: false,
// 				Code:    400,
// 				Error:   "File too large",
// 			})
// 			return
// 		}
//
// 		uploadTime := time.Now()
//
// 		filename := fmt.Sprintf("%s-%d%s", stringutil.TrimExt(file.Filename), uploadTime.Unix(), ext)
// 		dst := filepath.Join("uploads/videos", filename)
//
// 		if err := c.SaveUploadedFile(file, dst); err != nil {
// 			c.JSON(500, domain.APIResponse{
// 				Success: false,
// 				Code:    500,
// 				Error:   "Failed to save file",
// 			})
// 			return
// 		}
//
// 		preview, err := c.FormFile("preview")
// 		if err == nil {
// 			previewExt := filepath.Ext(preview.Filename)
// 			previewFilename := fmt.Sprintf("preview-%s-%d%s", stringutil.TrimExt(file.Filename), uploadTime.Unix(), previewExt)
// 			previewDst := filepath.Join("uploads/previews", previewFilename)
//
// 			if err := c.SaveUploadedFile(preview, previewDst); err != nil {
// 				c.JSON(500, domain.APIResponse{
// 					Success: false,
// 					Code:    500,
// 					Error:   "Failed to save preview",
// 				})
// 				return
// 			}
// 		}
//
// 		c.JSON(200, domain.APIResponse{
// 			Success: true,
// 			Data:    gin.H{"filename": filename},
// 			Message: "File uploaded successfully",
// 		})
// 	})
//
// 	return r
// }
