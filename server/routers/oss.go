package routers

// func UploadFileToOSS(ctx *gin.Context) {
// 	f, file, err := ctx.Request.FormFile("file")
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"code":    2,
// 			"message": "获取数据失败",
// 		})
// 		return
// 	}
// 	defer f.Close()

// 	if err != nil {
// 		fmt.Println("获取数据失败")
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"code":    1,
// 			"message": "获取数据失败",
// 		})
// 	} else {
// 		UploadFileToOSS()
// 		// buf := make([]byte, 512)

// 		// length, err := f.Read(buf)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	ctx.JSON(http.StatusBadRequest, gin.H{
// 		// 		"code":    1,
// 		// 		"message": "获取数据失败",
// 		// 	})
// 		// 	return
// 		// }
// 		// newOffset, err := f.Seek(0, 0)
// 		// if err != nil {
// 		// 	fmt.Println(err, newOffset)
// 		// 	ctx.JSON(http.StatusBadRequest, gin.H{
// 		// 		"code":    1,
// 		// 		"message": "获取数据失败",
// 		// 	})
// 		// 	return
// 		// }
// 		// contentType := http.DetectContentType(buf[:length])
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	ctx.JSON(http.StatusBadRequest, gin.H{
// 		// 		"code":    1,
// 		// 		"message": "获取数据失败",
// 		// 	})
// 		// 	return
// 		// }

// 		// ossHandler := utils.GetOSSHandler()
// 		// splitted := strings.Split(file.Filename, ".")
// 		// ext := ""
// 		// if len(splitted) > 1 {
// 		// 	ext = splitted[len(splitted)-1]
// 		// }
// 		// fileObjectName := uuid.New().String() + "." + ext
// 		// err = utils.PutObject(ossHandler.Buckets.Files, fileObjectName, contentType, f, -1)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	ctx.JSON(http.StatusBadRequest, gin.H{
// 		// 		"code":    1,
// 		// 		"message": "获取数据失败",
// 		// 	})
// 		// 	return
// 		// }
// 		// ctx.JSON(http.StatusOK, gin.H{
// 		// 	"code":    0,
// 		// 	"message": "success",
// 		// })
// 	}

// }
