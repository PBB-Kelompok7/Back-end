package user

type UserFormatter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	ImageURL string `json:"image_url"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Token:    token,
		ImageURL: user.AvatarFileName,
	}

	return formatter
}

type GetUserFormatter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageURL string `json:"image_url"`
}

func GetFormatUser(user User) GetUserFormatter {
	userFormatter := GetUserFormatter{}
	userFormatter.ID = user.ID
	userFormatter.Name = user.Name
	userFormatter.Email = user.Email
	userFormatter.ImageURL = user.AvatarFileName

	return userFormatter
}

func GetFormatUsers(users []User) []GetUserFormatter {
	usersFormatter := []GetUserFormatter{}

	for _, user := range users {
		userFormatter := GetFormatUser(user)
		usersFormatter = append(usersFormatter, userFormatter)
	}

	return usersFormatter
}