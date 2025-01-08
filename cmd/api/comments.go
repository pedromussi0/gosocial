package main

import (
	"net/http"

	"github.com/pedromussi0/gosocial.git/internal/store"
)

type CreateUpdateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUpdateCommentPayload
	post := getPostFromCtx(r)

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  103,
		Content: payload.Content,
		User:    store.User{ID: 103, Username: "Alice0"},
	}

	ctx := r.Context()

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
