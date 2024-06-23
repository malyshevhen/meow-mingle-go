package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/types"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
)

type PostHandler struct {
	postRepo db.IPostRepository
}

func NewPostHandler(postRepo db.IPostRepository) *PostHandler {
	return &PostHandler{
		postRepo: postRepo,
	}
}

func (ph *PostHandler) HandleCreatePost() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		postContent, err := ReadReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading post request: %v\n", "Post Handler", err)
			return err
		}

		params := utils.Map(postContent, func(s ContentForm) db.CreatePostParams {
			return db.CreatePostParams{Content: postContent.Content}
		})

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting user Id from token %v\n", "Post Handler ", err)
			return err
		}

		params.AuthorID = userId

		savedPost, err := ph.postRepo.CreatePost(ctx, params)
		if err != nil {
			log.Printf("%-15s ==> Error creating post in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created new post\n", "Post Handler")

		return utils.WriteJson(w, http.StatusCreated, savedPost)
	}
}

func (ph *PostHandler) HandleGetUserPosts() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id param %v\n", "Post Handler", err)
			return err
		}

		postResponses, err := ph.postRepo.ListUserPosts(ctx, id)
		if err != nil {
			return err
		}

		log.Printf("%-15s ==> Successfully retrieved user posts\n", "Post Handler")

		return utils.WriteJson(w, http.StatusOK, postResponses)
	}
}

func (ph *PostHandler) HandleGetPostsById() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para:%v\n", "Post Handler", err)
			return err
		}

		post, err := ph.postRepo.GetPost(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> Error getting post by Id from stor:%v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully retrieved post by Id\n", "Post Handler")

		return utils.WriteJson(w, http.StatusOK, post)
	}
}

func (ph *PostHandler) HandleUpdatePostsById() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			return err
		}

		postContent, err := ReadReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading update request: %v\n", "Post Handler", err)
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := utils.Map(postContent, func(content ContentForm) db.UpdatePostParams {
			return db.UpdatePostParams{
				ID:       id,
				Content:  content.Content,
				AuthorId: userId,
			}
		})

		postResponse, err := ph.postRepo.UpdatePost(ctx, params)
		if err != nil {
			log.Printf("%-15s ==> Error updating post by Id in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully updated post by Id\n", "Post Handler")

		return utils.WriteJson(w, http.StatusOK, postResponse)
	}
}

func (ph *PostHandler) HandleDeletePostsById() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id param %v\n", "Post Handler", err)
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		if err := ph.postRepo.DeletePost(ctx, userId, id); err != nil {
			log.Printf("%-15s ==> Error deleting post by Id from store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully deleted post by Id\n", "Post Handler")

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func (ph *PostHandler) HandleLikePost() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := utils.ParseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.CreatePostLikeParams{
			UserID: userId,
			PostID: id,
		}

		if err := utils.Validate(params); err != nil {
			return err
		}

		if err := ph.postRepo.CreatePostLike(context.Background(), params); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func (ph *PostHandler) HandleOwnersFeed() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		authUserID, err := utils.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		feed, err := ph.postRepo.GetFeed(ctx, authUserID)
		if err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusOK, feed)
	}
}

func (ph *PostHandler) HandleUsersFeed() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			return errors.NewValidationError("ID parameter is invalid")
		}

		feed, err := ph.postRepo.GetFeed(ctx, id)
		if err != nil {
			return err
		}
		return utils.WriteJson(w, http.StatusOK, feed)
	}
}

func (ph *PostHandler) HandleRemoveLikeFromPost() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.DeletePostLikeParams{
			PostID: id,
			UserID: userId,
		}

		if err := utils.Validate(params); err != nil {
			return err
		}

		if err := ph.postRepo.DeletePostLike(ctx, params); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}
