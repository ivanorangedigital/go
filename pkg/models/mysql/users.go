package mysql

import (
	"database/sql"
	"digitalcorporation/pkg/cache"
	"digitalcorporation/pkg/models"
	"encoding/json"
	"log"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert() (int, error) {
	return 0, nil
}

func (m *UserModel) Get() {}

func (m *UserModel) List() ([]*models.User, error) {
	// initialize key and user slice
	key := "sql_query:users"
	users := []*models.User{}

	item, err := cache.Get(key)

	// founded err
	if err != nil {
		return nil, nil
	}

	// item founded
	if item != nil {
		value := item.Value

		if err = json.Unmarshal(value, &users); err != nil {
			return nil, err
		}

		log.Println("users founds")

		return users, nil
	}

	log.Println("users not founds")

	// statment
	stmt := "SELECT ID, Email, Password, CreatedAt, UpdatedAt FROM USERS"

	rows, err := m.DB.Query(stmt)

	defer rows.Close()

	for rows.Next() {
		user := new(models.User)
		rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		users = append(users, user)
	}

	bytes, err := json.Marshal(users)

	if err != nil {
		return nil, err
	}

	cache.Set(key, bytes, 0)

	return users, nil
}
