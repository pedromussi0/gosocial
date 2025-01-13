package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pedromussi0/gosocial.git/internal/store"
	"golang.org/x/net/context"
)

type userKey string

const userCtx userKey = "user"

// getUserHandler godoc
//	@Summary		Get user by ID
//	@Description	Get user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/{ID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

// followUserHandler godoc
//	@Summary		Follow a user
//	@Description	Follow a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	FollowUser	true	"Follow User Payload"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/follow [post]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	//TODO: REVERT BACK TO AUTH USERID FROM CTX
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followerUser.ID, payload.UserID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// unfollowUserHandler godoc
//	@Summary		Unfollow a user
//	@Description	Unfollow a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	FollowUser	true	"Unfollow User Payload"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/unfollow [post]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowedUser := getUserFromContext(r)

	//TODO: REVERT BACK TO AUTH USERID FROM CTX
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
	}

	ctx := r.Context()

	if err := app.store.Followers.UnFollow(ctx, unFollowedUser.ID, payload.UserID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userID)
		//switch errors
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, ok := r.Context().Value(userCtx).(*store.User)
	if !ok {
		return nil
	}

	return user
}
