package cli

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NicolasFerreras/Gator/internal/config"
	"github.com/NicolasFerreras/Gator/internal/database"
	"github.com/google/uuid"
)

func sqlErrNoRows() error        { return sql.ErrNoRows }
func sqlNullStr(s string) sql.NullString { return sql.NullString{String: s, Valid: true} }
func sqlNullTime() sql.NullTime       { return sql.NullTime{Valid: false} }

// newMockDB devuelve un *database.Queries respaldado por sqlmock (sin Postgres).
func newMockDB(t *testing.T) (*database.Queries, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return database.New(db), mock
}

// isolateHome evita escribir ~/.gatorconfig.json real (necesario cuando el
// handler llama a Config.SetUser).
func isolateHome(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
}

func userRows(id uuid.UUID, username string) *sqlmock.Rows {
	now := time.Now()
	return sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username"}).
		AddRow(id, now, now, username)
}

func feedRows(id uuid.UUID, url string) *sqlmock.Rows {
	now := time.Now()
	return sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "url", "user_id", "last_fetched_at"}).
		AddRow(id, now, now, "Feed", url, uuid.New(), now)
}

// ---------- handlerLogin ----------

func TestHandlerLogin_NoArgs(t *testing.T) {
	q, _ := newMockDB(t)
	state := &State{Config: &config.Config{}, Db: q}
	err := handlerLogin(state, Command{Name: "login", args: []string{}})
	if err == nil {
		t.Fatal("se esperaba ErrNoUsername")
	}
}

func TestHandlerLogin_UserNotFound(t *testing.T) {
	q, mock := newMockDB(t)
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("ghost").WillReturnError(sqlErrNoRows())

	state := &State{Config: &config.Config{}, Db: q}
	err := handlerLogin(state, Command{Name: "login", args: []string{"ghost"}})
	if err == nil {
		t.Fatal("se esperaba error de usuario no encontrado")
	}
}

func TestHandlerLogin_Success(t *testing.T) {
	isolateHome(t)
	q, mock := newMockDB(t)
	id := uuid.New()
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").WillReturnRows(userRows(id, "nico"))

	state := &State{Config: &config.Config{}, Db: q}
	if err := handlerLogin(state, Command{Name: "login", args: []string{"nico"}}); err != nil {
		t.Fatalf("handlerLogin: %v", err)
	}
	if state.Config.CurrentUserName != "nico" {
		t.Fatalf("CurrentUserName = %q, want nico", state.Config.CurrentUserName)
	}
}

// ---------- handlerRegister ----------

func TestHandlerRegister_NoArgs(t *testing.T) {
	q, _ := newMockDB(t)
	state := &State{Config: &config.Config{}, Db: q}
	if err := handlerRegister(state, Command{Name: "register", args: []string{}}); err == nil {
		t.Fatal("se esperaba ErrNoUsername")
	}
}

func TestHandlerRegister_UserExists(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").WillReturnRows(userRows(id, "nico"))

	state := &State{Config: &config.Config{}, Db: q}
	if err := handlerRegister(state, Command{Name: "register", args: []string{"nico"}}); err == nil {
		t.Fatal("se esperaba ErrUserExists")
	}
}

func TestHandlerRegister_Success(t *testing.T) {
	isolateHome(t)
	q, mock := newMockDB(t)
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").WillReturnError(sqlErrNoRows())
	mock.ExpectQuery("INSERT INTO users").
		WillReturnRows(userRows(uuid.New(), "nico"))

	state := &State{Config: &config.Config{}, Db: q}
	if err := handlerRegister(state, Command{Name: "register", args: []string{"nico"}}); err != nil {
		t.Fatalf("handlerRegister: %v", err)
	}
	if state.Config.CurrentUserName != "nico" {
		t.Fatalf("CurrentUserName = %q, want nico", state.Config.CurrentUserName)
	}
}

// ---------- handlerUsers ----------

func TestHandlerUsers_ListsCurrent(t *testing.T) {
	q, mock := newMockDB(t)
	id := uuid.New()
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username"}).
			AddRow(uuid.New(), time.Now(), time.Now(), "nico").
			AddRow(uuid.New(), time.Now(), time.Now(), "otro"))
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").WillReturnRows(userRows(id, "nico"))

	state := &State{Config: &config.Config{CurrentUserName: "nico"}, Db: q}
	if err := handlerUsers(state, Command{Name: "users", args: []string{}}); err != nil {
		t.Fatalf("handlerUsers: %v", err)
	}
}

// ---------- handlerBrowse ----------

func TestHandlerBrowse_DefaultLimit(t *testing.T) {
	q, mock := newMockDB(t)
	uid := uuid.New()
	mock.ExpectQuery("SELECT id FROM users WHERE username").
		WithArgs("nico").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uid))
	mock.ExpectQuery("SELECT p.id, p.created_at").
		WithArgs(uid, int32(2), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "title", "url", "description", "published_at", "feed_id"}).
			AddRow(uuid.New(), time.Now(), time.Now(), "T", "https://x.com", sqlNullStr(""), sqlNullTime(), uuid.New()))

	state := &State{Config: &config.Config{CurrentUserName: "nico"}, Db: q}
	if err := handlerBrowse(state, Command{Name: "browse", args: []string{}}); err != nil {
		t.Fatalf("handlerBrowse default: %v", err)
	}
}

func TestHandlerBrowse_CustomLimit(t *testing.T) {
	q, mock := newMockDB(t)
	uid := uuid.New()
	mock.ExpectQuery("SELECT id FROM users WHERE username").
		WithArgs("nico").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uid))
	mock.ExpectQuery("SELECT p.id, p.created_at").
		WithArgs(uid, int32(5), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "title", "url", "description", "published_at", "feed_id"}).
			AddRow(uuid.New(), time.Now(), time.Now(), "T", "https://x.com", sqlNullStr(""), sqlNullTime(), uuid.New()))

	state := &State{Config: &config.Config{CurrentUserName: "nico"}, Db: q}
	if err := handlerBrowse(state, Command{Name: "browse", args: []string{"5"}}); err != nil {
		t.Fatalf("handlerBrowse custom: %v", err)
	}
}

func TestHandlerBrowse_TooManyArgs(t *testing.T) {
	q, _ := newMockDB(t)
	state := &State{Config: &config.Config{CurrentUserName: "nico"}, Db: q}
	if err := handlerBrowse(state, Command{Name: "browse", args: []string{"2", "x"}}); err == nil {
		t.Fatal("se esperaba ErrTooManyArguments")
	}
}

// ---------- handlerHelp ----------

func TestHandlerHelp(t *testing.T) {
	q, _ := newMockDB(t)
	state := &State{Config: &config.Config{}, Db: q}
	if err := handlerHelp(state, Command{Name: "help", args: []string{}}); err != nil {
		t.Fatalf("handlerHelp: %v", err)
	}
}

// ---------- middlewareLoggedIn + handlers con user ----------

func TestMiddlewareLoggedIn_NotLoggedIn(t *testing.T) {
	q, mock := newMockDB(t)
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("").WillReturnError(sqlErrNoRows())

	state := &State{Config: &config.Config{CurrentUserName: ""}, Db: q}
	called := false
	wrapped := middlewareLoggedIn(func(s *State, cmd Command, u database.User) error {
		called = true
		return nil
	})
	err := wrapped(state, Command{Name: "follow", args: []string{"https://x.com"}})
	if err == nil {
		t.Fatal("se esperaba error por no logueado")
	}
	if called {
		t.Fatal("el handler no debería ejecutarse sin usuario")
	}
}

func TestMiddlewareLoggedIn_Success(t *testing.T) {
	q, mock := newMockDB(t)
	uid := uuid.New()
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").WillReturnRows(userRows(uid, "nico"))

	state := &State{Config: &config.Config{CurrentUserName: "nico"}, Db: q}
	var gotUser database.User
	wrapped := middlewareLoggedIn(func(s *State, cmd Command, u database.User) error {
		gotUser = u
		return nil
	})
	if err := wrapped(state, Command{Name: "follow", args: []string{"https://x.com"}}); err != nil {
		t.Fatalf("middleware: %v", err)
	}
	if gotUser.ID != uid {
		t.Fatalf("usuario inyectado ID = %v, want %v", gotUser.ID, uid)
	}
}

func TestHandlerFollow_Success(t *testing.T) {
	q, mock := newMockDB(t)
	uid := uuid.New()
	feedID := uuid.New()
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").WillReturnRows(userRows(uid, "nico"))
	mock.ExpectQuery("SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds WHERE url").
		WithArgs("https://x.com").WillReturnRows(feedRows(feedID, "https://x.com"))
	mock.ExpectQuery("SELECT id FROM users WHERE username").
		WithArgs("nico").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uid))
	mock.ExpectQuery("WITH inserted_feed_follow").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), uid, feedID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "feed_id", "feed_name", "user_name"}).
			AddRow(uuid.New(), time.Now(), time.Now(), uid, feedID, "Feed", "nico"))

	state := &State{Config: &config.Config{CurrentUserName: "nico"}, Db: q}
	wrapped := middlewareLoggedIn(handlerFollow)
	if err := wrapped(state, Command{Name: "follow", args: []string{"https://x.com"}}); err != nil {
		t.Fatalf("handlerFollow: %v", err)
	}
}

func TestHandlerUnfollow_Success(t *testing.T) {
	q, mock := newMockDB(t)
	uid := uuid.New()
	feedID := uuid.New()
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").WillReturnRows(userRows(uid, "nico"))
	mock.ExpectQuery("SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds WHERE url").
		WithArgs("https://x.com").WillReturnRows(feedRows(feedID, "https://x.com"))
	mock.ExpectExec("DELETE FROM feed_follow").
		WithArgs(uid, feedID).WillReturnResult(sqlmock.NewResult(0, 1))

	state := &State{Config: &config.Config{CurrentUserName: "nico"}, Db: q}
	wrapped := middlewareLoggedIn(handlerUnfollowFeed)
	if err := wrapped(state, Command{Name: "unfollow", args: []string{"https://x.com"}}); err != nil {
		t.Fatalf("handlerUnfollow: %v", err)
	}
}

func TestHandlerFollowing_Success(t *testing.T) {
	q, mock := newMockDB(t)
	uid := uuid.New()
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").WillReturnRows(userRows(uid, "nico"))
	mock.ExpectQuery("SELECT").
		WithArgs(uid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "feed_id", "feed_name", "user_name"}).
			AddRow(uuid.New(), time.Now(), time.Now(), uid, uuid.New(), "Feed", "nico"))

	state := &State{Config: &config.Config{CurrentUserName: "nico"}, Db: q}
	wrapped := middlewareLoggedIn(handlerFollowin)
	if err := wrapped(state, Command{Name: "following", args: []string{}}); err != nil {
		t.Fatalf("handlerFollowing: %v", err)
	}
}

func TestHandlerAddFeed_Success(t *testing.T) {
	// FetchFeed hace HTTP real, pero a un server local (aislado, sin red externa).
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(`<?xml version="1.0"?><rss version="2.0"><channel><title>F</title><link>http://e.com</link><description>d</description><item><title>P</title><link>http://e.com/p</link><description>x</description><pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate></item></channel></rss>`))
	}))
	defer srv.Close()

	q, mock := newMockDB(t)
	uid := uuid.New()
	feedID := uuid.New()
	mock.ExpectQuery("SELECT id, created_at, updated_at, username FROM users WHERE username").
		WithArgs("nico").WillReturnRows(userRows(uid, "nico"))
	mock.ExpectQuery("INSERT INTO feeds").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "MiFeed", srv.URL, uid).
		WillReturnRows(feedRows(feedID, srv.URL))
	mock.ExpectQuery("SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds WHERE url").
		WithArgs(srv.URL).WillReturnRows(feedRows(feedID, srv.URL))
	mock.ExpectQuery("WITH inserted_feed_follow").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), uid, feedID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "feed_id", "feed_name", "user_name"}).
			AddRow(uuid.New(), time.Now(), time.Now(), uid, feedID, "MiFeed", "nico"))

	state := &State{Config: &config.Config{CurrentUserName: "nico"}, Db: q}
	wrapped := middlewareLoggedIn(handlerAddFeed)
	if err := wrapped(state, Command{Name: "addfeed", args: []string{"MiFeed", srv.URL}}); err != nil {
		t.Fatalf("handlerAddFeed: %v", err)
	}
}
