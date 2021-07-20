package internal

type ImageService struct {
	repository *ImageRepository
}

func NewImageService(repository *ImageRepository) *ImageService {
	return &ImageService{repository: repository}
}

func (s *ImageService) InsertImage(image *Image, user string) int {
	return s.repository.InsertImage(image, user)
}

func (s *ImageService) SelectImage(imageId int, user string) *Image {
	return s.repository.SelectImage(imageId, user)
}

func (s *ImageService) SelectImages(user string) Images {
	return s.repository.SelectImages(user)
}