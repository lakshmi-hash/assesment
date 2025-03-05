package controllers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestUserController(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "User Controller Suite")
}

var _ = ginkgo.Describe("UserController", func() {
	var uc *UserController
	var db *sql.DB

	ginkgo.BeforeEach(func() {
		var err error
		db, err = sql.Open("sqlite3", ":memory:")
		gomega.Expect(err).To(gomega.BeNil())
		_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL UNIQUE, email TEXT NOT NULL)")
		gomega.Expect(err).To(gomega.BeNil())
		sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
		uc = NewUserController(db, sq)
	})

	ginkgo.AfterEach(func() {
		db.Close()
	})

	ginkgo.It("should successfully create a user and count it", func() {
		_, err := uc.sq.Insert("users").Columns("name", "email").Values("test", "test@example.com").RunWith(db).Exec()
		gomega.Expect(err).To(gomega.BeNil())
		var count int
		err = uc.sq.Select("COUNT(*)").From("users").RunWith(db).QueryRow().Scan(&count)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(count).To(gomega.Equal(1))
	})

	ginkgo.It("should retrieve users correctly", func() {
		_, err := uc.sq.Insert("users").Columns("name", "email").Values("test", "test@example.com").RunWith(db).Exec()
		gomega.Expect(err).To(gomega.BeNil())

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = uc.GetUsers(c)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"name":"test"`))
	})

	ginkgo.It("should update a user successfully", func() {
		res, err := uc.sq.Insert("users").Columns("name", "email").Values("test", "test@example.com").RunWith(db).Exec()
		gomega.Expect(err).To(gomega.BeNil())
		id, _ := res.LastInsertId()

		_, err = uc.sq.Update("users").Set("name", "updated").Where(squirrel.Eq{"id": id}).RunWith(db).Exec()
		gomega.Expect(err).To(gomega.BeNil())

		var name string
		err = uc.sq.Select("name").From("users").Where(squirrel.Eq{"id": id}).RunWith(db).QueryRow().Scan(&name)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(name).To(gomega.Equal("updated"))
	})

	ginkgo.It("should delete a user successfully", func() {
		res, err := uc.sq.Insert("users").Columns("name", "email").Values("test", "test@example.com").RunWith(db).Exec()
		gomega.Expect(err).To(gomega.BeNil())
		id, _ := res.LastInsertId()

		_, err = uc.sq.Delete("users").Where(squirrel.Eq{"id": id}).RunWith(db).Exec()
		gomega.Expect(err).To(gomega.BeNil())

		var count int
		err = uc.sq.Select("COUNT(*)").From("users").RunWith(db).QueryRow().Scan(&count)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(count).To(gomega.Equal(0))
	})

	ginkgo.It("should handle duplicate username", func() {
		_, err := uc.sq.Insert("users").Columns("name", "email").Values("test", "test@example.com").RunWith(db).Exec()
		gomega.Expect(err).To(gomega.BeNil())

		_, err = uc.sq.Insert("users").Columns("name", "email").Values("test", "test2@example.com").RunWith(db).Exec()
		gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
		gomega.Expect(err.Error()).To(gomega.ContainSubstring("UNIQUE constraint failed"))
	})
})
