package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type MockRepository struct {
	FindByEmailFunc func(email string) (User, error)
	FindByIDFunc    func(ID int) (User, error)
	SaveFunc        func(user User) (User, error)
	FindAllFunc     func() ([]User, error)
	UpdateFunc      func(user User) (User, error)
}

func (m *MockRepository) Save(user User) (User, error) {
	if m.SaveFunc != nil {
			return m.SaveFunc(user)
	}
	return user, nil
}

func (m *MockRepository) FindByEmail(email string) (User, error) {
	if m.FindByEmailFunc != nil {
			return m.FindByEmailFunc(email)
	}
	return User{}, errors.New("FindByEmailFunc is not implemented")
}

func (m *MockRepository) FindByID(ID int) (User, error) {
	if ID == 1 {
			return User{ID: 1, Name: "John", Email: "existing@example.com", PasswordHash: "hashed_password", Role: "user"}, nil
	}
	return User{}, nil
}

func (m *MockRepository) Update(user User) (User, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return user, nil
}

func (m *MockRepository) FindAll() ([]User, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return []User{}, nil
}

func (m *MockRepository) Delete(ID int) error {
	return nil
}

// TestRegisterUser
func TestRegisterUser(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	input := RegisterUserInput{Name: "John", Email: "john@example.com", Password: "password"}
	user, err := service.RegisterUser(input)

	assert.NoError(t, err)
	assert.Equal(t, "John", user.Name)
}

func TestRegisterUser_RepositoryError(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	// Mock user input
	input := RegisterUserInput{
			Name:     "John",
			Email:    "john@example.com",
			Password: "password",
	}

	// Mock repository's Save method to return an error
	repo.SaveFunc = func(user User) (User, error) {
			return User{}, errors.New("repository error")
	}

	// Perform user registration
	_, err := service.RegisterUser(input)

	// Assert error occurred
	assert.Error(t, err)
	// Assert the error message is correct
	assert.EqualError(t, err, "repository error")
}

// TestLogin
func TestLogin_Success(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	// Mock user data
	email := "existing@example.com"
	password := "password"

	// Mock hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	// Mock user returned by repository
	mockUser := User{
			ID:           1,
			Name:         "John",
			Email:        email,
			PasswordHash: string(hashedPassword),
			Role:         "user",
	}

	// Mock repository's FindByEmail method to return the mock user
	repo.FindByEmailFunc = func(email string) (User, error) {
			if email == mockUser.Email {
					return mockUser, nil
			}
			return User{}, errors.New("user not found")
	}

	// Input for login
	input := LoginInput{
			Email:    email,
			Password: password,
	}

	// Perform login
	loggedInUser, err := service.Login(input)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the correct user is returned
	assert.Equal(t, mockUser, loggedInUser)
}

func TestLogin_IncorrectPassword(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	// Mock user data
	email := "existing@example.com"
	password := "password"

	// Mock hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	// Mock user returned by repository
	mockUser := User{
			ID:           1,
			Name:         "John",
			Email:        email,
			PasswordHash: string(hashedPassword),
			Role:         "user",
	}

	// Mock repository's FindByEmail method to return the mock user
	repo.FindByEmailFunc = func(email string) (User, error) {
			if email == mockUser.Email {
					return mockUser, nil
			}
			return User{}, errors.New("user not found")
	}

	// Input for login with incorrect password
	input := LoginInput{
			Email:    email,
			Password: "wrongpassword",
	}

	// Perform login
	_, err := service.Login(input)

	// Assert error occurred
	assert.Error(t, err)
	// Assert error message is correct
	assert.EqualError(t, err, "crypto/bcrypt: hashedPassword is not the hash of the given password")
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	// Mock repository's FindByEmail method to return no user found error
	repo.FindByEmailFunc = func(email string) (User, error) {
			return User{}, errors.New("No user found")
	}

	// Input for login with non-existing email
	input := LoginInput{
			Email:    "nonexisting@example.com",
			Password: "password",
	}

	// Perform login
	_, err := service.Login(input)

	// Assert error occurred
	assert.Error(t, err)
	// Assert error message is correct
	assert.EqualError(t, err, "No user found")
}

// Get User By ID
func TestGetUserByID(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	user, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "John", user.Name)
}

func TestGetUserByID_UserNotFoundError(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	// Test getting non-existing user by ID
	user, err := service.GetUserByID(2)

	assert.Error(t, err, "expected error when getting non-existing user by ID")
	assert.Equal(t, User{}, user)
}

// Update User
func TestUpdateUser(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	input := FormUpdateUserInput{ID: 1, Name: "Updated John", Email: "updated@example.com"}
	user, err := service.UpdateUser(input)

	assert.NoError(t, err)
	assert.Equal(t, "Updated John", user.Name)
}

func TestUpdateUser_InvalidID(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	initialUser, _ := service.GetUserByID(0)

	input := FormUpdateUserInput{ID: 0, Name: "Updated John", Email: "updated@example.com"}
	_, err := service.UpdateUser(input)

	assert.NoError(t, err, "unexpected error when updating user with invalid ID")

	unchangedUser, _ := service.GetUserByID(0)
	assert.Equal(t, initialUser, unchangedUser, "expected user to remain unchanged")
}

func TestUpdateUser_RepositoryError(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	// Mock user input
	input := FormUpdateUserInput{ID: 1, Name: "Updated John", Email: "updated@example.com"}

	// Mock repository's FindByID method to return a user
	repo.FindByIDFunc = func(ID int) (User, error) {
		return User{ID: 1, Name: "John", Email: "john@example.com", PasswordHash: "hashed_password", Role: "user"}, nil
	}

	// Mock repository's Update method to return an error
	repo.UpdateFunc = func(user User) (User, error) {
		return User{}, errors.New("repository error")
	}

	// Perform user update
	_, err := service.UpdateUser(input)

	// Assert error occurred
	assert.Error(t, err)
	// Assert the error message is correct
	assert.EqualError(t, err, "repository error")
}

// Cek Email
func TestIsEmailAvailable_EmailAvailable(t *testing.T) {
	// Simulate an available email
	repo := &MockRepository{}
	repo.FindByEmailFunc = func(email string) (User, error) {
			// Return nil error to simulate email not found
			return User{}, nil
	}
	service := NewService(repo)

	input := CheckEmailInput{Email: "new@example.com"}
	available, err := service.IsEmailAvailable(input)

	assert.NoError(t, err)
	assert.True(t, available, "expected email to be available")
}

func TestIsEmailAvailable_EmailNotAvailable(t *testing.T) {
	// Simulate an unavailable email
	repo := &MockRepository{}
	repo.FindByEmailFunc = func(email string) (User, error) {
			// Return a user to simulate email found
			return User{ID: 1}, nil
	}
	service := NewService(repo)

	input := CheckEmailInput{Email: "existing@example.com"}
	available, err := service.IsEmailAvailable(input)

	assert.NoError(t, err)
	assert.False(t, available, "expected email to be unavailable")
}

func TestIsEmailAvailable_FindByEmailError(t *testing.T) {
	repo := &MockRepository{}
	repo.FindByEmailFunc = func(email string) (User, error) {
			return User{}, errors.New("find by email error")
	}
	service := NewService(repo)

	input := CheckEmailInput{Email: "new@example.com"}
	available, err := service.IsEmailAvailable(input)

	assert.Error(t, err, "expected error when finding by email")
	assert.False(t, available, "expected email to be unavailable")
}

// Get All Users
func TestGetAllUsers(t *testing.T) {
	// Define mock users data
	mockUsers := []User{
			{ID: 1, Name: "John", Email: "john@example.com", PasswordHash: "hashed_password", Role: "user"},
			{ID: 2, Name: "Alice", Email: "alice@example.com", PasswordHash: "hashed_password", Role: "user"},
	}

	// Create a mock repository with a FindAll function that returns mock users
	repo := &MockRepository{}
	service := NewService(repo)

	// Mock repository's FindAll method to return mock users
	repo.FindAllFunc = func() ([]User, error) {
			return mockUsers, nil
	}

	// Perform the GetAllUsers operation
	users, err := service.GetAllUsers()

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the number of returned users matches the number of mock users
	assert.Equal(t, len(mockUsers), len(users))
}

func TestGetAllUsers_Error(t *testing.T) {
	// Create a mock repository with a FindAll function that returns an error
	repo := &MockRepository{}
	service := NewService(repo)

	// Mock repository's FindAll method to return an error
	repo.FindAllFunc = func() ([]User, error) {
		return nil, errors.New("error retrieving users")
	}

	// Perform the GetAllUsers operation
	_, err := service.GetAllUsers()

	// Assert that an error occurred
	assert.Error(t, err)
	assert.EqualError(t, err, "error retrieving users")
}

// Avatar
func TestSaveAvatar(t *testing.T) {
	// Create a mock repository
	repo := &MockRepository{}
	service := NewService(repo)

	// Mock user data
	mockUser := User{
		ID:             1,
		Name:           "John",
		Email:          "john@example.com",
		PasswordHash:   "hashed_password",
		Role:           "user",
		AvatarFileName: "",
	}

	// Mock repository's FindByID method to return the mock user
	repo.FindByIDFunc = func(ID int) (User, error) {
		if ID == mockUser.ID {
			return mockUser, nil
		}
		return User{}, errors.New("user not found")
	}

	// Define file location for avatar
	fileLocation := "../../images/first_test.jpg"

	// Mock repository's Update method to return updated user
	repo.UpdateFunc = func(user User) (User, error) {
		mockUser.AvatarFileName = fileLocation
		return mockUser, nil
	}

	// Perform SaveAvatar operation
	userWithAvatar, err := service.SaveAvatar(mockUser.ID, fileLocation)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert that the avatar file location is updated
	assert.Equal(t, fileLocation, userWithAvatar.AvatarFileName)
}

func TestSaveAvatar_Error(t *testing.T) {
	// Create a mock repository
	repo := &MockRepository{}
	service := NewService(repo)

	// Mock user data
	mockUser := User{
		ID:             1,
		Name:           "John",
		Email:          "john@example.com",
		PasswordHash:   "hashed_password",
		Role:           "user",
		AvatarFileName: "",
	}

	// Mock repository's FindByID method to return the mock user
	repo.FindByIDFunc = func(ID int) (User, error) {
		if ID == mockUser.ID {
			return mockUser, nil
		}
		return User{}, errors.New("user not found")
	}

	// Define file location for avatar
	fileLocation := "../../images/first_test.jpg"

	// Mock repository's Update method to return an error
	repo.UpdateFunc = func(user User) (User, error) {
		return User{}, errors.New("error updating user")
	}

	// Perform SaveAvatar operation
	_, err := service.SaveAvatar(mockUser.ID, fileLocation)

	// Assert an error occurred
	assert.Error(t, err)
	assert.EqualError(t, err, "error updating user")
}




