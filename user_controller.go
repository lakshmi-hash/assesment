package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/user-management-app/backend/models"
)

// UserController handles CRUD operations
type UserController struct {
	db *sql.DB
	sq squirrel.StatementBuilderType
}

func NewUserController(db *sql.DB, sq squirrel.StatementBuilderType) *UserController {
	return &UserController{db: db, sq: sq}
}

// GetUsers retrieves all users
// @Summary Get all users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (uc *UserController) GetUsers(c echo.Context) error {
	rows, err := uc.sq.Select("id, name, email").From("users").RunWith(uc.db).Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan users"})
		}
		users = append(users, u)
	}
	return c.JSON(http.StatusOK, users)
}

// CreateUser adds a new user
// @Summary Create a user
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (uc *UserController) CreateUser(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Check if username exists
	var existingID int
	err := uc.sq.Select("id").From("users").Where(squirrel.Eq{"name": u.Name}).RunWith(uc.db).QueryRow().Scan(&existingID)
	if err == nil { // If no error, a row was found, meaning the username exists
		return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
	}
	if err != sql.ErrNoRows { // If error is not "no rows", something went wrong
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check username"})
	}

	// Insert new user
	res, err := uc.sq.Insert("users").Columns("name", "email").Values(u.Name, u.Email).RunWith(uc.db).Exec()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}
	id, _ := res.LastInsertId()
	u.ID = int(id)
	return c.JSON(http.StatusCreated, u)
}

// UpdateUser modifies an existing user
// @Summary Update a user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/{id} [put]
func (uc *UserController) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	u.ID = id

	_, err := uc.sq.Update("users").Set("name", u.Name).Set("email", u.Email).Where(squirrel.Eq{"id": id}).RunWith(uc.db).Exec()
	if err != nil {
		if err.Error() == "SQLITE_CONSTRAINT: UNIQUE constraint failed: users.name" {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}
	return c.JSON(http.StatusOK, u)
}

// DeleteUser removes a user
// @Summary Delete a user
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := uc.sq.Delete("users").Where(squirrel.Eq{"id": id}).RunWith(uc.db).Exec()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}
	return c.NoContent(http.StatusNoContent)
}
