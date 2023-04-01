package routers

// func GetAllUserPosts(ctx *gin.Context) {
// 	userID := ctx.Query("userID")
// 	oid, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	posts, err := services.GetPostsService().GetAllUserPosts(oid)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	postOutlines := make([]models.PostOutline, 0)
// 	for i := 0; i < len(posts); i++ {
// 		postOutline, err := services.GetPostsService().PostToOutline(&posts[i])
// 		if err != nil {
// 			ctx.String(http.StatusBadRequest, err.Error())
// 			return
// 		}
// 		postOutlines = append(postOutlines, postOutline)
// 	}
// 	ctx.JSON(http.StatusOK, postOutlines)
// }
