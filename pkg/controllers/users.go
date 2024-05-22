package controllers

import (
	db "client_task/pkg/common/db/sqlc"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"client_task/pkg/payloads"

	"database/sql"

	"github.com/gin-gonic/gin"
	// "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
)

type UsersController struct {
	DB  *db.Queries
	CTX context.Context
}

func NewUsersController(db *db.Queries, ctx context.Context) *UsersController {
	return &UsersController{db, ctx}
}

func (cc *UsersController) CreateUser(ctx *gin.Context) {
	var payload *payloads.CreateUser
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(payload.PasswordHash)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	args := &db.CreateUserParams{
		FullName:     payload.FullName,
		Email:        payload.Email,
		PasswordHash: hashedPassword,
		UserType:     payload.UserType,
	}

	user, err := cc.DB.CreateUser(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving user", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully created user", "user": user})
}

func (cc *UsersController) UpdateUser(ctx *gin.Context) {
	var payload *payloads.UpdateUser
	UserId := ctx.Param("UserId")
	user_id, _ := strconv.Atoi(UserId)

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(payload.PasswordHash)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	args := &db.UpdateUserParams{
		ID:           int32(user_id),
		FullName:     sql.NullString{String: payload.FullName, Valid: payload.FullName != ""}.String,
		PasswordHash: hashedPassword,
	}

	User, err := cc.DB.UpdateUser(ctx, *args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve User with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving User", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully updated User", "User": User})
}

func (cc *UsersController) GetUserById(ctx *gin.Context) {
	UserId := ctx.Param("UserId")
	user_id, _ := strconv.Atoi(UserId)
	println(int32(user_id))
	User, err := cc.DB.GetUserByID(ctx, int32(user_id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve User with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving User", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "User": User})
}

func (cc *UsersController) GetAllUsers(ctx *gin.Context) {

	users, err := cc.DB.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed to retrieve Users", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrieved all Users", "size": len(users), "users": users})
}

func (cc *UsersController) DeleteUserById(ctx *gin.Context) {
	UserId := ctx.Param("UserId")
	user_id, _ := strconv.Atoi(UserId)

	_, err := cc.DB.GetUserByID(ctx, int32(user_id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve User with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving User", "error": err.Error()})
		return
	}

	err = cc.DB.DeleteUserByID(ctx, int32(user_id))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfuly deleted"})

}

// profile

type ProfilesController struct {
	DB  *db.Queries
	CTX context.Context
}

func NewProfilesController(db *db.Queries, ctx context.Context) *ProfilesController {
	return &ProfilesController{db, ctx}
}

func profileSkillsCreation(ctx context.Context, pc_db *db.Queries, skills []string, userID int32) error {

	fmt.Println("skills", skills)
	for _, skill := range skills {

		skillObj, err := pc_db.GetSkillByName(ctx, skill)
		// fmt.Println("skill objection")
		// fmt.Println(skillObj, err.Error())
		if err != nil {
			skillObj, _ = pc_db.CreateSkill(ctx, skill)
			// fmt.Println(skillObj)
		}

		params := &db.GetUserSkillsByUserIDAndSkillIdParams{
			UserID:  userID,
			SkillID: skillObj.ID,
		}

		profile_skills, err := pc_db.GetUserSkillsByUserIDAndSkillId(ctx, *params)
		// fmt.Println("profile_skills")
		// fmt.Println(profile_skills, err.Error())
		if err != nil {

			Cparams := &db.AddSkillToUserParams{
				UserID:  userID,
				SkillID: skillObj.ID,
			}
			pc_db.AddSkillToUser(ctx, *Cparams)
		}
		fmt.Println(profile_skills)

	}
	return nil
}

func (pc *ProfilesController) CreateProfile(ctx *gin.Context) {
	var payload *payloads.CreateProfile
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	profile_check, err := pc.DB.GetProfileByUserID(ctx, payload.UserID)
	if err == nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed creating profile", "error": "Profile Already exist"})
		fmt.Println(profile_check)
		return
	}

	args := &db.CreateProfileParams{
		UserID:      payload.UserID,
		Bio:         pgtype.Text{String: payload.Bio.String, Valid: payload.Bio.String != ""},
		Company:     pgtype.Text{String: payload.Company, Valid: payload.Company != ""},
		JobRole:     payload.JobRole,
		Description: pgtype.Text{String: payload.Description.String, Valid: payload.Description.String != ""},
	}

	profile, err := pc.DB.CreateProfile(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed creating profile", "error": err.Error()})
		return
	}

	profileSkillsCreation(ctx, pc.DB, payload.Skills, payload.UserID)

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully created profile", "profile": profile})
}

func (pc *ProfilesController) UpdateProfile(ctx *gin.Context) {
	var payload *payloads.UpdateProfile
	profileID := ctx.Param("ProfileID")
	id, err := strconv.Atoi(profileID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Invalid profile ID", "error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	args := &db.UpdateProfileParams{
		ID:          int32(id),
		Bio:         pgtype.Text{String: payload.Bio.String, Valid: payload.Bio.String != ""},
		Company:     pgtype.Text{String: payload.Company, Valid: payload.Company != ""},
		JobRole:     payload.JobRole,
		Description: pgtype.Text{String: payload.Description.String, Valid: payload.Description.String != ""},
	}

	profile, err := pc.DB.UpdateProfile(ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Profile not found"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed updating profile", "error": err.Error()})
		return
	}
	profileSkillsCreation(ctx, pc.DB, payload.Skills, profile.UserID)

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully updated profile", "profile": profile})
}

func (pc *ProfilesController) GetProfileByID(ctx *gin.Context) {
	profileID := ctx.Param("ProfileID")
	id, err := strconv.Atoi(profileID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Invalid profile ID", "error": err.Error()})
		return
	}

	profile, err := pc.DB.GetProfileByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Profile not found"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving profile", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived profile", "profile": profile})
}

func (pc *ProfilesController) DeleteProfileByID(ctx *gin.Context) {
	ProfileID := ctx.Param("ProfileID")
	profile_id, _ := strconv.Atoi(ProfileID)

	_, err := pc.DB.GetProfileByID(ctx, int32(profile_id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve profile with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving profile", "error": err.Error()})
		return
	}

	err = pc.DB.DeleteProfileByID(ctx, int32(profile_id))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfuly deleted"})

}

func (pc *ProfilesController) GetAllProfiles(ctx *gin.Context) {

	profiles, err := pc.DB.GetAllProfile(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed to retrieve profile", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrieved all profile", "size": len(profiles), "profiles": profiles})
}

type CareersController struct {
	DB  *db.Queries
	CTX context.Context
}

func NewCareersController(db *db.Queries, ctx context.Context) *CareersController {
	return &CareersController{db, ctx}
}

type SkillsController struct {
	DB  *db.Queries
	CTX context.Context
}

func NewSkillsController(db *db.Queries, ctx context.Context) *SkillsController {
	return &SkillsController{db, ctx}
}

func (cc *SkillsController) GetUserSkills(ctx *gin.Context) {
	userID := ctx.Param("UserID")
	id, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Invalid user ID", "error": err.Error()})
		return
	}

	userSkills, err := cc.DB.GetAllUserSkillsByUserID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "No skills found for this user"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed to retrieve user skills", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrieved user skills", "size": len(userSkills), "skills": userSkills})
}

func (cc *SkillsController) CreateSkill(ctx *gin.Context) {
	var payload payloads.CreateSkill
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	skill, err := cc.DB.CreateSkill(ctx, payload.Name)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed to create skill", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully created skill", "skill": skill})
}

func (cc *SkillsController) GetAllSkills(ctx *gin.Context) {
	skills, err := cc.DB.GetAllskills(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed to retrieve skills", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrieved all skills", "size": len(skills), "skills": skills})
}

func (cc *CareersController) CreateCareer(ctx *gin.Context) {
	var payload *payloads.CreateCareer
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	args := &db.CreateCareerParams{
		UserID:      payload.UserID,
		Title:       payload.Title,
		Company:     pgtype.Text{String: payload.Company, Valid: payload.Company != ""},
		Description: pgtype.Text{String: payload.Description, Valid: payload.Description != ""},
		SkillID:     payload.SkillID,
	}

	career, err := cc.DB.CreateCareer(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed creating career", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully created career", "career": career})
}

func (cc *CareersController) UpdateCareer(ctx *gin.Context) {
	var payload *payloads.UpdateCareer
	careerID := ctx.Param("CareerID")
	id, err := strconv.Atoi(careerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Invalid career ID", "error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	args := &db.UpdateCareerParams{
		ID:          int32(id),
		Title:       payload.Title,
		Company:     pgtype.Text{String: payload.Company, Valid: payload.Company != ""},
		Description: pgtype.Text{String: payload.Description, Valid: payload.Description != ""},
		SkillID:     payload.SkillID,
	}

	career, err := cc.DB.UpdateCareer(ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Career not found"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed updating career", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully updated career", "career": career})
}

func (cc *CareersController) GetCareerByID(ctx *gin.Context) {
	careerID := ctx.Param("CareerID")
	id, err := strconv.Atoi(careerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Invalid career ID", "error": err.Error()})
		return
	}

	career, err := cc.DB.GetCareerByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Career not found"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving career", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrieved career", "career": career})
}

func (cc *CareersController) GetAllCareers(ctx *gin.Context) {

	careers, err := cc.DB.GetAllCareers(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed to retrieve careers", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrieved all careers", "size": len(careers), "careers": careers})
}

func (cc *CareersController) DeleteCareerByID(ctx *gin.Context) {
	careerID := ctx.Param("CareerID")
	id, err := strconv.Atoi(careerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Invalid career ID", "error": err.Error()})
		return
	}

	_, err = cc.DB.GetCareerByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Career not found"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving career", "error": err.Error()})
		return
	}

	err = cc.DB.DeleteCareerByID(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully deleted career"})
}

func (cc *CareersController) GetCareersByUserID(ctx *gin.Context) {
	userID := ctx.Param("UserID")
	id, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Invalid user ID", "error": err.Error()})
		return
	}

	careers, err := cc.DB.GetCareersByUserID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "No careers found for this user"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving careers", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrieved careers", "size": len(careers), "careers": careers})
}
