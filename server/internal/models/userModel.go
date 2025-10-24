package models

type UserRegister struct {
	User_id      string `bson:"user_id" json:"user_id"`
	Name         string `bson:"name" json:"name" binding:"required"`
	PhomeNumber  string `bson:"phone_number" json:"phone_number" omitempty:"true"`
	Role         string `bson:"role" json:"role" default:"user"`
	Email        string `bson:"email" json:"email" binding:"required"`
	Password     string `bson:"password" json:"password" binding:"required,min=6"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

type UserLogin struct {
	Email        string `bson:"email" json:"email" binding:"requied"`
	Password     string `bson:"password" json:"password" binding:"required,min=6"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}


