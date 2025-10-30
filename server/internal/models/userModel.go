package models

type UserRegister struct {
	User_id     string `bson:"user_id" json:"user_id"`
	Name        string `bson:"name" json:"name" binding:"required"`
	PhoneNumber string `bson:"phone_number" json:"phone_number" omitempty:"true"`
	Role        string `bson:"role" json:"role" default:"user"`
	Email       string `bson:"email" json:"email" binding:"required"`
	Password    string `bson:"password" json:"password" binding:"required,min=6"`
}
type UserRegisterResponse struct {
	User_id      string `bson:"user_id" json:"user_id"`
	Name         string `bson:"name" json:"name" `
	PhoneNumber  string `bson:"phone_number" json:"phone_number" omitempty:"true"`
	Role         string `bson:"role" json:"role" default:"user"`
	Email        string `bson:"email" json:"email" `
	Password     string `bson:"password" json:"password"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

type UserLogin struct {
	Email    string `bson:"email" json:"email" binding:"required"`
	Password string `bson:"password" json:"password" binding:"required,min=6"`
}
type UserLoginResponse struct {
	User_id      string `bson:"user_id" json:"user_id"`
	Name         string `bson:"name" json:"name"`
	Email        string `bson:"email" json:"email"`
	PhoneNumber  string `bson:"phone_number" json:"phone_number" omitempty:"true"`
	Role         string `bson:"role" json:"role" default:"user"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

type UserSearchResponse struct {
	User_id    string `bson:"user_id" json:"user_id"`
	Name       string `bson:"name" json:"name"`
	ProfileIMG string `bson:"profile_pic" json:"profile_pic"`
}
