package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/malyshEvhen/meow_mingle/config"
	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
	"github.com/malyshEvhen/meow_mingle/errors"
	"github.com/malyshEvhen/meow_mingle/types"
)

func (rr *Router) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error reading request body: %v\n", "UserService", err)
		return err
	}
	defer r.Body.Close()

	createUserParams, err := Unmarshal[db.CreateUserParams](body)
	if err != nil {
		log.Printf("%-15s ==> 😕 Error unmarshal JSON: %v\n", "UserService", err)
		return err
	}

	user := types.UserFromParams(createUserParams)

	log.Printf("%-15s ==> 👀 Validating user payload: %s\n", "UserService", user.String())
	if err := Validate(user); err != nil {
		return err
	}

	log.Printf("%-15s ==> 🔑 Hashing password...", "UserService")
	hashedPwd, err := HashPwd(createUserParams.Password)
	if err != nil {
		log.Printf("%-15s ==> 🔒 Error hashing password: %v\n", "UserService", err)
		return err
	}

	createUserParams.Password = hashedPwd

	log.Printf("%-15s ==> 📝 Creating user in database...\n", "UserService")
	u, err := rr.store.CreateUserTx(ctx, createUserParams)
	if err != nil {
		log.Printf("%-15s ==> 🛑 Error creating user: %v\n", "UserService", err)
		return err
	}

	log.Printf("%-15s ==> 🔐 Creating auth token...\n", "UserService")
	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		log.Printf("%-15s ==> ❌ Error creating auth token: %v\n", "UserService", err)
		return err
	}

	log.Printf("%-15s ==> ✅ User created successfully!\n", "UserService")
	return WriteJson(w, http.StatusCreated, token)
}

func (rr *Router) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return errors.NewValidationError("ID parameter is invalid")
	}

	userID, err := GetAuthUserId(r)
	if err != nil {
		log.Printf("%-15s ==> ❌ No authenticated user found", "UserService")
		return err
	}

	if id != int(userID) {
		log.Printf("%-15s ==> ❌ User with ID: %d have no permissions to access account with ID: %d\n", "UserService", userID, id)
		return errors.NewForbiddenError()
	}

	log.Printf("%-15s ==> 🕵️ Searching for user with Id:%s\n", "UserService", strId)

	u, err := rr.store.GetUser(ctx, int64(id))
	if err != nil {
		log.Printf("%-15s ==> 😕 User not found for Id:%d\n", "UserService", id)
		return err
	}

	log.Printf("%-15s ==> 👍 Found user: %d\n", "UserService", u.ID)

	return WriteJson(w, http.StatusOK, u)
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	log.Printf("%-15s ==> 🔑 Generating JWT token..\n", "UserService.")
	secret := config.Envs.JWTSecret
	token, err := CreateJwt([]byte(secret), id)
	if err != nil {
		log.Printf("%-15s ==> ❌ Error generating JWT token: %s\n", "UserService", err)
		return "", err
	}

	log.Printf("%-15s ==> 🍪 Setting auth cookie..\n", "UserService.")
	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	log.Printf("%-15s ==> ✅ Auth cookie set successfully!\n", "UserService")
	return token, nil
}

func (rr *Router) handleCreatePost(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	postRequest, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error reading post request: %v\n", "PostController", err)
		return err
	}

	userId, err := GetAuthUserId(r)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error getting user Id from token %v\n", "PostController ", err)
		return err
	}

	createPostParams := &db.CreatePostParams{
		AuthorID: userId,
		Content:  postRequest.Content,
	}

	postResponse, err := rr.store.CreatePost(ctx, *createPostParams)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating post in store %v\n", "PostController", err)
		WriteJson(w, http.StatusInternalServerError, types.NewErrorResponse("Error creating post"))
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully created new post\n", "PostController")

	return WriteJson(w, http.StatusCreated, postResponse)
}

func (rr *Router) handleGetUserPosts(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id param %v\n", "PostController", err)
		return err
	}

	posts, err := rr.store.ListUserPosts(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting user posts from store %v\n", "PostController", err)
		return err
	}

	postResponses := []types.PostResponse{}
	for _, p := range posts {
		likes, err := rr.store.ListPostLikes(ctx, p.ID)
		if err != nil {
			return err
		}

		postResponses = append(postResponses, types.PostResponse{
			Id:       p.ID,
			Content:  p.Content,
			AuthorId: p.AuthorID,
			Likes:    len(likes),
			Created:  p.CreatedAt,
			Updated:  p.UpdatedAt,
		})

	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved user posts\n", "PostController")

	return WriteJson(w, http.StatusOK, postResponses)
}

func (rr *Router) handleGetPostsById(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para:%v\n", "PostController", err)
		return err
	}

	post, err := rr.store.GetPost(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting post by Id from stor:%v\n", "PostController", err)
		return err
	}

	postResponse := &types.PostResponse{
		Id:       post.ID,
		Content:  post.Content,
		AuthorId: post.AuthorID,
		Created:  post.CreatedAt,
		Updated:  post.UpdatedAt,
	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved post by Id\n", "PostController")

	return WriteJson(w, http.StatusOK, postResponse)
}

func (rr *Router) handleUpdatePostsById(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostController", err)
		return err
	}

	postRequest, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading post request %v\n", "PostController", err)
		return err
	}

	params := &db.UpdatePostParams{
		ID:      id,
		Content: postRequest.Content,
	}

	postResponse, err := rr.store.UpdatePost(ctx, *params)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error updating post by Id in store %v\n", "PostController", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully updated post by Id\n", "PostController")

	return WriteJson(w, http.StatusOK, postResponse)
}

func (rr *Router) handleDeletePostsById(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id param %v\n", "PostController", err)
		return err
	}

	if err := rr.store.DeletePost(ctx, id); err != nil {
		log.Printf("%-15s ==> 😫 Error deleting post by Id from store %v\n", "PostController", err)
		return err
	}

	log.Printf("%-15s ==> 🗑️ Successfully deleted post by Id\n", "PostController")

	return WriteJson(w, http.StatusNoContent, nil)
}

func readPostReqType(r *http.Request) (*types.PostRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.NewValidationError("parameter ID is not valid")
	}
	defer r.Body.Close()

	p, err := Unmarshal[types.PostRequest](body)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func parseIdParam(r *http.Request) (int64, error) {
	id := r.PathValue("id")

	numId, err := strconv.Atoi(id)
	if err != nil {
		return 0, nil
	}

	return int64(numId), nil
}

func (rr *Router) handleCreateComment(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	postId, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing post Id param %v\n", "PostService ", err)
		return err
	}

	cReq, err := readCommentReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		return err
	}

	userId, err := GetAuthUserId(r)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error getting authenticated user Id %v\n", "PostService ", err)
		return err
	}

	params := &db.CreateCommentParams{
		AuthorID: userId,
		PostID:   postId,
		Content:  cReq.Content,
	}

	cResp, err := rr.store.CreateComment(ctx, *params)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating comment in store %v\n", "PostService ", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully created comment\n", "PostService")

	return WriteJson(w, http.StatusCreated, cResp)
}

func (rr *Router) handleGetComments(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostService ", err)
		return err
	}

	c, err := rr.store.ListPostComments(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting comment by Id from stor %v\n", "PostService ", err)
		return err
	}

	commentResponses := []types.CommentResponse{}
	for _, comment := range c {
		likes, err := rr.store.ListCommentLikes(ctx, comment.ID)
		if err != nil {
			return err
		}

		commentResponses = append(commentResponses, types.CommentResponse{
			Id:       comment.ID,
			Content:  comment.Content,
			AuthorId: comment.AuthorID,
			PostId:   comment.PostID,
			Likes:    len(likes),
			Created:  comment.CreatedAt,
			Updated:  comment.UpdatedAt,
		})
	}

	log.Printf("%-15s ==> 🎉 Successfully got comment by Id\n", "PostService!")

	return WriteJson(w, http.StatusOK, commentResponses)
}

func (rr *Router) handleUpdateComments(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostService ", err)
		return err

	}

	c, err := readCommentReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		return err

	}

	params := &db.UpdateCommentParams{
		ID:      id,
		Content: c.Content,
	}

	cr, err := rr.store.UpdateComment(ctx, *params)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error updating comment by Id in stor %v\n", "PostService ", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully updated comment by Id\n", "PostService")

	return WriteJson(w, http.StatusOK, cr)
}

func (rr *Router) handleDeleteComments(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para\n ", "PostService")
		return err

	}

	err = rr.store.DeleteComment(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error deleting comment by Id from stor\n ", "PostService")
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully deleted comment by Id\n", "PostService")

	return WriteJson(w, http.StatusNoContent, nil)
}

func readCommentReqType(r *http.Request) (*types.CommentRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c, err := Unmarshal[types.CommentRequest](body)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (rr *Router) handleLikePost(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		return err
	}

	userId, err := GetAuthUserId(r)
	if err != nil {
		return err
	}

	params := db.CreatePostLikeParams{
		UserID: userId,
		PostID: id,
	}

	if err := rr.store.CreatePostLike(context.Background(), params); err != nil {
		return err
	}

	return WriteJson(w, http.StatusNoContent, nil)
}

func (rr *Router) handleRemoveLikeFromPost(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		return err
	}

	userId, err := GetAuthUserId(r)
	if err != nil {
		return err
	}

	post, err := rr.store.GetPost(context.Background(), id)
	if err != nil {
		return err
	}

	if post.AuthorID != userId {
		return fmt.Errorf("user with ID: %d can not modify post of author with ID: %d", userId, post.AuthorID)
	}

	params := db.DeletePostLikeParams{
		PostID: id,
		UserID: userId,
	}

	if err := rr.store.DeletePostLike(context.Background(), params); err != nil {
		return err
	}

	return WriteJson(w, http.StatusNoContent, nil)
}

func (rr *Router) handleLikeComment(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		return err
	}

	userId, err := GetAuthUserId(r)
	if err != nil {
		return err
	}

	params := db.CreateCommentLikeParams{
		UserID:    userId,
		CommentID: id,
	}

	if err := rr.store.CreateCommentLike(context.Background(), params); err != nil {
		return err
	}

	return WriteJson(w, http.StatusNoContent, nil)
}

func (rr *Router) handleRemoveLikeFromComment(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		return err
	}

	userId, err := GetAuthUserId(r)
	if err != nil {
		return err
	}

	comment, err := rr.store.GetComment(context.Background(), id)
	if err != nil {
		return err
	}

	if comment.AuthorID != userId {
		return fmt.Errorf("user with ID: %d can not modify post of author with ID: %d", userId, comment.ID)
	}

	params := db.DeleteCommentLikeParams{
		CommentID: id,
		UserID:    userId,
	}

	if err := rr.store.DeleteCommentLike(context.Background(), params); err != nil {
		return err
	}

	return WriteJson(w, http.StatusNoContent, nil)
}
