package post

import "github.com/sgoldenf/a-place-for-your-thoughts/internal/models"

type PostModelInterface interface {
	CreatePost(title, text string) (int, error)
	ReadPost(id int) (*models.Post, error)
	UpdatePost(id int, title, text string) (*models.Post, error)
	DeletePost(id int) (int, error)
	GetPostsList(page int) ([]*models.Post, error)
	GetPostsCount() (count int, err error)
}
