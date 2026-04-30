package users

type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UpdateUserInput struct {
    Username      *string `json:"username,omitempty"`
    NoTelp        *string `json:"no_telp,omitempty"`
    JenisKelamin *bool   `json:"jenis_kelamin"`
    ProfileImage  *string `json:"profile_image,omitempty"`
}