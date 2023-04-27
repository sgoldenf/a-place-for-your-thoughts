package mocks

import (
	"errors"

	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models"
)

var mockPost = &models.Post{
	Id:    1,
	Title: "Title 1",
	Text:  "Text 1",
}

type PostModel struct{}

func (m *PostModel) CreatePost(title, text string) (int, error) {
	return 2, nil
}

func (m *PostModel) ReadPost(id int) (*models.Post, error) {
	switch id {
	case 1:
		return mockPost, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *PostModel) UpdatePost(id int, title, text string) (*models.Post, error) {
	switch id {
	case 1:
		return &models.Post{Id: id, Title: title, Text: text}, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *PostModel) DeletePost(id int) (int, error) {
	switch id {
	case 1:
		return id, nil
	default:
		return 0, models.ErrNoRecord
	}
}

func (m *PostModel) GetPostsList(page int) ([]*models.Post, error) {
	switch page {
	case 1:
		return []*models.Post{mockPost}, nil
	default:
		return nil, errors.New("not found")
	}
}

func (m *PostModel) GetPostsCount() (int, error) {
	return 1, nil
}
