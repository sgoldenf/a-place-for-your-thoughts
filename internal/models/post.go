package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	Id    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Text  string `json:"text" db:"text"`
}

type PostModel struct {
	Pool *pgxpool.Pool
}

func (m *PostModel) CreatePost(title, text string) (int, error) {
	var id int
	err := m.Pool.QueryRow(
		context.Background(),
		`insert into posts (title, text) values ($1, $2) returning id;`,
		title, text,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *PostModel) ReadPost(id int) (*Post, error) {
	post := &Post{}
	err := m.Pool.QueryRow(
		context.Background(),
		`select * from posts where id=$1;`, id,
	).Scan(&post.Id, &post.Title, &post.Text)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return post, nil
}

func (m *PostModel) UpdatePost(id int, title, text string) (*Post, error) {
	updated := &Post{Id: id}
	err := m.Pool.QueryRow(
		context.Background(),
		`update posts set title = $1, text = $2 where id = $3 returning title, text;`,
		title, text, id,
	).Scan(&updated.Title, &updated.Text)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (m *PostModel) DeletePost(id int) (int, error) {
	var deleted int
	err := m.Pool.QueryRow(
		context.Background(),
		`delete from posts where id=$1 returning id;`, id,
	).Scan(&deleted)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (m *PostModel) GetPostsList(page int) ([]*Post, error) {
	rows, err := m.Pool.Query(context.Background(), `select * from posts order by id desc limit 10 offset $1;`, (page-1)*3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Text)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (m *PostModel) GetPostsCount() (count int, err error) {
	err = m.Pool.QueryRow(
		context.Background(),
		`select count (*) from posts;`,
	).Scan(&count)
	return
}
