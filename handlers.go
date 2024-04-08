package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

func (rr *Router) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error reading request body: %v\n", "UserService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Invalid payload"))
		return
	}
	defer r.Body.Close()

	user, err := Unmarshal[User](body)
	if err != nil {
		log.Printf("%-15s ==> 😕 Error unmarshal JSON: %v\n", "UserService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Invalid payload"))
		return
	}

	log.Printf("%-15s ==> 👀 Validating user payload: %v\n", "UserService", user)
	if err := validateUserPayload(user); err != nil {
		log.Printf("%-15s ==> ❌ Validation failed: %v\n", "UserService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse(err.Error()))
		return
	}

	log.Printf("%-15s ==> 🔑 Hashing password...", "UserService")
	hashedPwd, err := HashPwd(user.Password)
	if err != nil {
		log.Printf("%-15s ==> 🔒 Error hashing password: %v\n", "UserService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Invalid payload"))
		return
	}

	user.Password = hashedPwd

	params := &db.CreateUserParams{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	log.Printf("%-15s ==> 📝 Creating user in database...\n", "UserService")
	u, err := rr.store.CreateUser(ctx, *params)
	if err != nil {
		log.Printf("%-15s ==> 🛑 Error creating user: %v\n", "UserService", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating user"))
		return
	}

	log.Printf("%-15s ==> 🔐 Creating auth token...\n", "UserService")
	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		log.Printf("%-15s ==> ❌ Error creating auth token: %v\n", "UserService", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating auth token"))
		return
	}

	log.Printf("%-15s ==> ✅ User created successfully!\n", "UserService")
	WriteJson(w, http.StatusCreated, token)
}

func (rr *Router) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Invalid param",
		}
	}

	userID, err := GetAuthUserId(r)
	if err != nil {
		log.Printf("%-15s ==> ❌ No authenticated user found", "UserService")
		return &BasicError{
			Code:    http.StatusUnauthorized,
			Message: "access denied",
		}
	}

	if id != int(userID) {
		log.Printf("%-15s ==> ❌ User with ID: %d have no permissions to access account with ID: %d\n", "UserService", userID, id)
		return &BasicError{
			Code:    http.StatusUnauthorized,
			Message: "access denied",
		}

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
	secret := Envs.JWTSecret
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

func validateUserPayload(user User) error {
	log.Printf("%-15s ==> 📧 Checking if email is provided..", "UserService.")
	if user.Email == "" {
		log.Printf("%-15s ==> ❌ Email is required but not provided", "UserService")
		return errEmailRequired
	}

	log.Printf("%-15s ==> 📛 Checking if first name is provided..", "UserService.")
	if user.FirstName == "" {
		log.Printf("%-15s ==> ❌ First name is required but not provided", "UserService")
		return errFirstNameRequired
	}

	log.Printf("%-15s ==> 📛 Checking if last name is provided..", "UserService.")
	if user.LastName == "" {
		log.Printf("%-15s ==> ❌ Last name is required but not provided", "UserService")
		return errLastNameRequired
	}

	log.Printf("%-15s ==> 🔑 Checking if password is provided..", "UserService.")
	if user.Password == "" {
		log.Printf("%-15s ==> ❌ Password is required but not provided", "UserService")
		return errPasswordRequired
	}

	log.Printf("%-15s ==> ✅ User payload validation passed!", "UserService")
	return nil
}

func (rr *Router) handleCreatePost(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	postRequest, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error reading post request: %v\n", "PostController", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error reading post request",
		}
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
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating post"))
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
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error reading post request",
		}
	}

	posts, err := rr.store.ListUserPosts(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting user posts from store %v\n", "PostController", err)
		return err
	}

	postResponses := []PostResponse{}
	for _, p := range posts {
		likes, err := rr.store.ListPostLikes(ctx, p.ID)
		if err != nil {
			return err
		}

		postResponses = append(postResponses, PostResponse{
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
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error reading post request",
		}
	}

	post, err := rr.store.GetPost(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting post by Id from stor:%v\n", "PostController", err)
		return err
	}

	postResponse := &PostResponse{
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
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Invalid param",
		}
	}

	postRequest, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading post request %v\n", "PostController", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Invalid content",
		}
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
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Invalid param",
		}
	}

	if err := rr.store.DeletePost(ctx, id); err != nil {
		log.Printf("%-15s ==> 😫 Error deleting post by Id from store %v\n", "PostController", err)
		return err
	}

	log.Printf("%-15s ==> 🗑️ Successfully deleted post by Id\n", "PostController")

	return WriteJson(w, http.StatusNoContent, nil)
}

func readPostReqType(r *http.Request) (*PostRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	p, err := Unmarshal[PostRequest](body)
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
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Not valid ID param",
		}
	}

	cReq, err := readCommentReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error reading comment request",
		}
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
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error parsing Id param",
		}
	}

	c, err := rr.store.ListPostComments(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting comment by Id from stor %v\n", "PostService ", err)
		return err
	}

	commentResponses := []CommentResponse{}
	for _, comment := range c {
		likes, err := rr.store.ListCommentLikes(ctx, comment.ID)
		if err != nil {
			return err
		}

		commentResponses = append(commentResponses, CommentResponse{
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
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Not valid ID param",
		}

	}

	c, err := readCommentReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Not valid ID param",
		}

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
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Not valid ID param",
		}

	}

	err = rr.store.DeleteComment(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error deleting comment by Id from stor\n ", "PostService")
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully deleted comment by Id\n", "PostService")

	return WriteJson(w, http.StatusNoContent, nil)
}

func readCommentReqType(r *http.Request) (*CommentRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c, err := Unmarshal[CommentRequest](body)
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
		return fmt.Errorf("user with ID: %d can not modify post of author with ID: %d\n", userId, post.AuthorID)
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
		return fmt.Errorf("user with ID: %d can not modify post of author with ID: %d\n", userId, comment.ID)
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
