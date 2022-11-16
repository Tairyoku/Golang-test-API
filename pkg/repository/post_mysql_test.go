package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"test"
	"testing"
)

type args struct {
	userId int
	post   test.Post
}

type ddd struct {
	db   *gorm.DB
	mock sqlmock.Sqlmock
	args
}

func TestPostMySql_Create(t *testing.T) {
	s := &ddd{}
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	if db == nil {
		t.Error("mock db is null")
	}

	if s.mock == nil {
		t.Error("sqlmock is null")
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to open gorm v2 db, got error: %v", err)
	}

	if s.db == nil {
		t.Error("gorm db is null")
	}

	s.args = args{
		userId: 1,
		post: test.Post{
			Title: "title",
			Anons: "anons",
		},
	}
	defer db.Close()

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"postId"}).AddRow(2)
	s.mock.ExpectQuery("INSERT INTO posts").
		WithArgs(s.args.post.Title, s.args.post.Anons).WillReturnRows(rows)

	s.mock.ExpectExec("INSERT INTO posts_list").
		WithArgs(s.args.userId, 2).WillReturnResult(sqlmock.NewResult(1, 2))
	s.mock.ExpectCommit()

	if err = s.db.Create(&s.args.post).Error; err != nil {
		t.Errorf("Failed to insert to gorm db, got error: %v", err)
	}

	err = s.mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

//	//var db *sql.DB
//	//var err error
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer db.Close()
//
//	//db, mock, err := sqlmock.New()
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//cnf := Config{
//	//	Username: "root",
//	//	Password: "",
//	//	Host:     "tcp",
//	//	Url:      "127.0.0.1:3306",
//	//	DBName:   "myFirstDB",
//	//}
//	var (
//		db  *sql.DB
//		err error
//	)
//
//	db, s.mock, err = sqlmock.New()
//	if err != nil {
//		t.Errorf("Failed to open mock sql db, got error: %v", err)
//	}
//
//	if db == nil {
//		t.Error("mock db is null")
//	}
//
//	if s.mock == nil {
//		t.Error("sqlmock is null")
//	}
//
//	dialector := postgres.New(postgres.Config{
//		DSN:                  "sqlmock_db_0",
//		DriverName:           "postgres",
//		Conn:                 db,
//		PreferSimpleProtocol: true,
//	})
//	s.db, err = gorm.Open(dialector, &gorm.Config{})
//	if err != nil {
//		t.Errorf("Failed to open gorm v2 db, got error: %v", err)
//	}
//
//	if s.db == nil {
//		t.Error("gorm db is null")
//	}
//	type args struct {
//		userId int
//		post   test.Post
//	}
//	type mockBehavior func(postId int, args args)
//
//	testTable := []struct {
//		name         string
//		mockBehavior mockBehavior
//		args         args
//		postId       int
//		wantErr      bool
//	}{
//		{
//			name: "ok",
//			args: args{
//				userId: 1,
//				post: test.Post{
//					Title: "title",
//					Anons: "anons",
//				},
//			},
//			postId: 2,
//			mockBehavior: func(postId int, args args) {
//				mock.ExpectBegin()
//				rows := sqlmock.NewRows([]string{"postId"}).AddRow(postId)
//				mock.ExpectQuery("INSERT INTO posts").
//					WithArgs(args.post.Title, args.post.Anons).WillReturnRows(rows)
//
//				mock.ExpectExec("INSERT INTO posts_list").
//					WithArgs(args.userId, postId).WillReturnResult(sqlmock.NewResult(1, 1))
//				mock.ExpectCommit()
//
//			},
//		},
//	}
//
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mockBehavior(testCase.postId, testCase.args)
//			got, err := r.Create(testCase.args.userId, testCase.args.post)
//			if testCase.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoErrorf(t, err, "")
//				assert.Equal(t, testCase.postId, got)
//			}
//		})
//	}
//}
//
