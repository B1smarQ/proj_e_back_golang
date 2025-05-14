package usecase

import "post_service/internal/entity"

type ThreadUsecase struct {
	ThreadRepository entity.ThreadRepository
}

func NewThreadUsecase(repo entity.ThreadRepository) *ThreadUsecase {
	return &ThreadUsecase{
		ThreadRepository: repo,
	}
}

func (u *ThreadUsecase) GetThreadsUsecase() ([]entity.ReturnThread, error) {
	return u.ThreadRepository.GetThreads()
}

func (u *ThreadUsecase) CreateThreadUsecase(thread entity.Thread) (entity.Thread, error) {
	return u.ThreadRepository.CreateThread(thread)
}

func (u *ThreadUsecase) UpdateThreadUsecase(thread entity.Thread) (entity.Thread, error) {
	return u.ThreadRepository.UpdateThread(thread)
}

func (u *ThreadUsecase) DeleteThreadUsecase(thread entity.Thread) (entity.Thread, error) {
	return u.ThreadRepository.DeleteThread(thread)
}

func (u *ThreadUsecase) GetThreadUsecase(id int) (entity.ReturnThread, error) {
	return u.ThreadRepository.GetThread(id)
}

func (u *ThreadUsecase) GetThreadPostsUsecase(threadID int) ([]entity.ReturnPost, error) {
	return u.ThreadRepository.GetThreadPosts(threadID)
}
