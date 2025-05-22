package usecase

import "post_service/internal/entity"

type PostUsecase struct {
	PostRepository entity.PostRepository
}

func NewPostUsecase(repo entity.PostRepository) *PostUsecase {
	return &PostUsecase{
		PostRepository: repo,
	}
}

func (u *PostUsecase) GetPostsUsecase() ([]entity.ReturnPost, error) {
	return u.PostRepository.GetPosts()
}

func (u *PostUsecase) CreatePostUsecase(post entity.Post) (entity.Post, error) {
	return u.PostRepository.CreatePost(post)
}

func (u *PostUsecase) UpdatePostUsecase(post entity.Post) (entity.Post, error) {
	return u.PostRepository.UpdatePost(post)
}

func (u *PostUsecase) DeletePostUsecase(post entity.Post) (entity.Post, error) {
	return u.PostRepository.DeletePost(post)
}

func (u *PostUsecase) GetPostUsecase(id string) (entity.ReturnPost, error) {
	return u.PostRepository.GetPost(id)
}
