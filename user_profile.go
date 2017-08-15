package main

import (
	"github.com/JormungandrK/microservice-user-profile/app"
	"github.com/JormungandrK/microservice-user-profile/db"
	"github.com/goadesign/goa"

	"gopkg.in/mgo.v2/bson"
)

// UserProfileController implements the userProfile resource.
type UserProfileController struct {
	*goa.Controller
	Repository db.UserProfileRepository
}

// NewUserProfileController creates a userProfile controller.
func NewUserProfileController(service *goa.Service, Repository db.UserProfileRepository) *UserProfileController {
	return &UserProfileController{
		Controller: service.NewController("UserProfileController"),
		Repository: Repository,
	}
}

// GetUserProfile runs the GetUserProfile action.
func (c *UserProfileController) GetUserProfile(ctx *app.GetUserProfileUserProfileContext) error {
	// Build the resource using the generated data structure.
	res := &app.UserProfile{}

	// Return whether ctx.UserID is a valid hex representation of an ObjectId.
	if bson.IsObjectIdHex(ctx.UserID) != true {
		return ctx.NotFound(goa.ErrNotFound("Invalid User Id"))
	}

	// Return an ObjectId from the provided hex representation.
	userID := bson.ObjectIdHex(ctx.UserID)

	// Return true if userID is valid. A valid userID must contain exactly 12 bytes.
	if userID.Valid() != true {
		return ctx.NotFound(goa.ErrNotFound("Invalid User Id"))
	}

	// Return one user profile by id.
	if err := c.Repository.GetUserProfile(userID, res); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if res.CreatedOn == 0 {
		return ctx.NotFound(goa.ErrNotFound("User Profile not found"))
	}

	res.UserID = ctx.UserID

	return ctx.OK(res)
}
