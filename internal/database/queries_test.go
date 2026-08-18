package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

// newMockDB crea un *Queries respaldado por sqlmock, sin Postgres real.
func newMockDB(t *testing.T) (*Queries, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return New(db), mock
}

func TestCreateUser(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	now := time.Now()

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(id, now, now, "nico").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username"}).
			AddRow(id, now, now, "nico"))

	got, err := q.CreateUser(context.Background(), CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  "nico",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if got.Username != "nico" || got.ID != id {
		t.Fatalf("CreateUser devolvió %+v", got)
	}
}

func TestGetUserByUsername_Found(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username"}).
			AddRow(id, now, now, "nico"))

	u, err := q.GetUserByUsername(context.Background(), "nico")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	if u.Username != "nico" {
		t.Fatalf("username = %q, want nico", u.Username)
	}
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	q, mock := newMockDB(t)

	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("ghost").
		WillReturnError(sql.ErrNoRows)

	_, err := q.GetUserByUsername(context.Background(), "ghost")
	if err == nil {
		t.Fatal("se esperaba error sql.ErrNoRows")
	}
}

func TestGetUsers(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username"}).
			AddRow(id, now, now, "nico").
			AddRow(uuid.New(), now, now, "otro"))

	users, err := q.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("GetUsers: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("len(users) = %d, want 2", len(users))
	}
}

func TestGetUserID(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()

	mock.ExpectQuery("SELECT id FROM users WHERE username").
		WithArgs("nico").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

	got, err := q.GetUserID(context.Background(), "nico")
	if err != nil {
		t.Fatalf("GetUserID: %v", err)
	}
	if got != id {
		t.Fatalf("GetUserID = %v, want %v", got, id)
	}
}

func TestCreateFeed(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("INSERT INTO feeds").
		WithArgs(id, now, now, "TechCrunch", "https://tc.com", userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "url", "user_id", "last_fetched_at"}).
			AddRow(id, now, now, "TechCrunch", "https://tc.com", userID, now))

	got, err := q.CreateFeed(context.Background(), CreateFeedParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "TechCrunch",
		Url:       "https://tc.com",
		UserID:    userID,
	})
	if err != nil {
		t.Fatalf("CreateFeed: %v", err)
	}
	if got.Name != "TechCrunch" {
		t.Fatalf("Name = %q", got.Name)
	}
}

func TestGetFeedByURL(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds WHERE url").
		WithArgs("https://tc.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "url", "user_id", "last_fetched_at"}).
			AddRow(id, now, now, "TechCrunch", "https://tc.com", userID, now))

	got, err := q.GetFeedByURL(context.Background(), "https://tc.com")
	if err != nil {
		t.Fatalf("GetFeedByURL: %v", err)
	}
	if got.Url != "https://tc.com" {
		t.Fatalf("Url = %q", got.Url)
	}
}

func TestGetNextFeedToFetch(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()

	mock.ExpectQuery("SELECT id, url FROM feeds").
		WillReturnRows(sqlmock.NewRows([]string{"id", "url"}).AddRow(id, "https://tc.com"))

	got, err := q.GetNextFeedToFetch(context.Background())
	if err != nil {
		t.Fatalf("GetNextFeedToFetch: %v", err)
	}
	if got.ID != id || got.Url != "https://tc.com" {
		t.Fatalf("GetNextFeedToFetch = %+v", got)
	}
}

func TestMarkFeedFetched(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	now := time.Now()

	mock.ExpectExec("UPDATE feeds").
		WithArgs(now, now, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := q.MarkFeedFetched(context.Background(), MarkFeedFetchedParams{
		UpdatedAt:     now,
		LastFetchedAt: now,
		ID:            id,
	}); err != nil {
		t.Fatalf("MarkFeedFetched: %v", err)
	}
}

func TestCreateFeedFollow(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	userID := uuid.New()
	feedID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("WITH inserted_feed_follow").
		WithArgs(id, now, now, userID, feedID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "feed_id", "feed_name", "user_name"}).
			AddRow(id, now, now, userID, feedID, "TechCrunch", "nico"))

	got, err := q.CreateFeedFollow(context.Background(), CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		FeedID:    feedID,
	})
	if err != nil {
		t.Fatalf("CreateFeedFollow: %v", err)
	}
	if got.FeedName != "TechCrunch" || got.UserName != "nico" {
		t.Fatalf("CreateFeedFollow = %+v", got)
	}
}

func TestDeleteFeedFollowByUserID(t *testing.T) {
	q, mock := newMockDB(t)
	userID := uuid.New()
	feedID := uuid.New()

	mock.ExpectExec("DELETE FROM feed_follow").
		WithArgs(userID, feedID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := q.DeleteFeedFollowByUserID(context.Background(), DeleteFeedFollowByUserIDParams{
		UserID: userID,
		FeedID: feedID,
	}); err != nil {
		t.Fatalf("DeleteFeedFollowByUserID: %v", err)
	}
}

func TestCreatePost(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	feedID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("INSERT INTO posts").
		WithArgs(id, now, now, "Title", "https://p.com", sql.NullString{String: "desc", Valid: true}, sql.NullTime{Valid: false}, feedID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "title", "url", "description", "published_at", "feed_id"}).
			AddRow(id, now, now, "Title", "https://p.com", sql.NullString{String: "desc", Valid: true}, sql.NullTime{Valid: false}, feedID))

	got, err := q.CreatePost(context.Background(), CreatePostParams{
		ID:          id,
		CreatedAt:   now,
		UpdatedAt:   now,
		Title:       "Title",
		Url:         "https://p.com",
		Description: sql.NullString{String: "desc", Valid: true},
		PublishedAt: sql.NullTime{Valid: false},
		FeedID:      feedID,
	})
	if err != nil {
		t.Fatalf("CreatePost: %v", err)
	}
	if got.Title != "Title" {
		t.Fatalf("Title = %q", got.Title)
	}
}

func TestGetPostsByUserId(t *testing.T) {
	q, mock := newMockDB(t)
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT p.id, p.created_at").
		WithArgs(userID, int32(2), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "title", "url", "description", "published_at", "feed_id"}).
			AddRow(uuid.New(), now, now, "T1", "https://p1.com", sql.NullString{Valid: false}, sql.NullTime{Valid: false}, uuid.New()).
			AddRow(uuid.New(), now, now, "T2", "https://p2.com", sql.NullString{Valid: false}, sql.NullTime{Valid: false}, uuid.New()))

	posts, err := q.GetPostsByUserId(context.Background(), GetPostsByUserIdParams{
		UserID: userID,
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("GetPostsByUserId: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("len(posts) = %d, want 2", len(posts))
	}
}
