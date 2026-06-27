package profile

type Service interface {
	GetProfile(userID uint) (*ProfileResponse, error)
	UpdateProfile(userID uint, request UpdateProfileRequest) (*ProfileResponse, error)
	UpdateProfileImage(userID uint, imageURL string) (*ProfileResponse, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) GetProfile(userID uint) (*ProfileResponse, error) {
	return s.repository.FindByUserID(userID)
}

func (s *service) UpdateProfile(
	userID uint,
	request UpdateProfileRequest,
) (*ProfileResponse, error) {
	return s.repository.Update(userID, request)
}

func (s *service) UpdateProfileImage(
	userID uint,
	imageURL string,
) (*ProfileResponse, error) {
	return s.repository.UpdateProfileImage(userID, imageURL)
}
