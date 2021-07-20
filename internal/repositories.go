package internal

type ImageRepository struct {
	db FakeDB
}

func NewImageRepository(db FakeDB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) InsertImage(image *Image, user string) int {
	// Since we are using 0 for invalid imageIds we have to start from 1.
	counter := len(r.db[user]) + 1
	if _, ok := r.db[user]; !ok {
		r.db[user] = make(Images)
	}
	r.db[user][counter] = image
	return counter
}

func (r *ImageRepository) SelectImage(imageId int, user string) *Image {
	return r.db[user][imageId]
}

func (r *ImageRepository) SelectImages(user string) Images {
	return r.db[user]
}
