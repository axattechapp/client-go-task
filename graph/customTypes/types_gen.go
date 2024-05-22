// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package customTypes

import (
	"fmt"
	"io"
	"strconv"
)

type Career struct {
	ID          int     `json:"id"`
	UserID      int     `json:"userId"`
	Title       string  `json:"title"`
	Company     *string `json:"company,omitempty"`
	Description *string `json:"description,omitempty"`
	Skill       *Skill  `json:"skill"`
	User        *User   `json:"user"`
}

type CreateCareerInput struct {
	Title       string  `json:"title"`
	Company     *string `json:"company,omitempty"`
	Description *string `json:"description,omitempty"`
	SkillID     int     `json:"skillId"`
}

type CreateProfileInput struct {
	Bio         *string  `json:"bio,omitempty"`
	Company     *string  `json:"company,omitempty"`
	JobRole     string   `json:"jobRole"`
	Description *string  `json:"description,omitempty"`
	Skills      []string `json:"skills"`
}

type CreateUserInput struct {
	FullName  string   `json:"fullName"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Role      UserRole `json:"role"`
	AvatarURL *string  `json:"avatarUrl,omitempty"`
}

type Mutation struct {
}

type Profile struct {
	ID          int      `json:"id"`
	UserID      int      `json:"userId"`
	Bio         *string  `json:"bio,omitempty"`
	Company     *string  `json:"company,omitempty"`
	JobRole     string   `json:"jobRole"`
	Description *string  `json:"description,omitempty"`
	User        *User    `json:"user"`
	Skills      []*Skill `json:"skills"`
}

type ProfileSkill struct {
	ID      int      `json:"id"`
	Profile *Profile `json:"profile"`
	Skill   *Skill   `json:"skill"`
}

type Query struct {
}

type Skill struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SkillInput struct {
	Name string `json:"name"`
}

type UpdateCareerInput struct {
	Title       *string `json:"title,omitempty"`
	Company     *string `json:"company,omitempty"`
	Description *string `json:"description,omitempty"`
	SkillID     *int    `json:"skillId,omitempty"`
}

type UpdateProfileInput struct {
	Bio         *string  `json:"bio,omitempty"`
	Company     *string  `json:"company,omitempty"`
	JobRole     *string  `json:"jobRole,omitempty"`
	Description *string  `json:"description,omitempty"`
	Skills      []string `json:"skills"`
}

type UpdateUserInput struct {
	FullName  *string   `json:"fullName,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Password  *string   `json:"password,omitempty"`
	Role      *UserRole `json:"role,omitempty"`
	AvatarURL *string   `json:"avatarUrl,omitempty"`
}

type User struct {
	ID        int      `json:"id"`
	FullName  string   `json:"fullName"`
	Email     string   `json:"email"`
	Role      UserRole `json:"role"`
	AvatarURL string   `json:"avatarUrl"`
	CreatedAt string   `json:"createdAt"`
}

type UserFilter struct {
	Email *string   `json:"email,omitempty"`
	Role  *UserRole `json:"role,omitempty"`
}

type UserRole string

const (
	UserRoleJobseeker UserRole = "JOBSEEKER"
	UserRoleAdmin     UserRole = "ADMIN"
	UserRoleRecruiter UserRole = "RECRUITER"
)

var AllUserRole = []UserRole{
	UserRoleJobseeker,
	UserRoleAdmin,
	UserRoleRecruiter,
}

func (e UserRole) IsValid() bool {
	switch e {
	case UserRoleJobseeker, UserRoleAdmin, UserRoleRecruiter:
		return true
	}
	return false
}

func (e UserRole) String() string {
	return string(e)
}

func (e *UserRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserRole", str)
	}
	return nil
}

func (e UserRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}