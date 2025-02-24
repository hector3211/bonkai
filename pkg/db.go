package pkg

import (
	"encoding/json"
	"fmt"
	"my-chi/models"
	"my-chi/utils"
	"os"
)

type DB struct {
	UserData []models.User
}

func PingDB() error {
	_, err := os.Stat(utils.UserDataPath)
	if err != nil {
		fmt.Println("[DB] ping failed, no such file exists")
		return err
	}
	return nil
}

func NewDB() (*DB, error) {
	err := PingDB()
	if err != nil {
		file, err := os.Create(utils.UserDataPath)
		if err != nil {
			fmt.Println("[DB] Failed creating user-data file:", err)
			return nil, err
		}
		defer file.Close()

		emptyData := []models.User{}
		data, _ := json.Marshal(emptyData)
		file.Write(data)

		fmt.Println("[DB] Successfully created user-data file")
		return &DB{UserData: emptyData}, nil
	}

	contents, err := os.ReadFile(utils.UserDataPath)
	if err != nil {
		fmt.Println("[DB] Error reading user-data file:", err)
		return nil, err
	}

	var users []models.User
	err = json.Unmarshal(contents, &users)
	if err != nil {
		fmt.Println("[DB] Error parsing JSON:", err)
		return nil, err
	}

	fmt.Println("[DB] Successfully loaded user-data file")
	return &DB{UserData: users}, nil
}

func (d *DB) UpdateWithData(data models.User) error {
	contents, err := os.ReadFile(utils.UserDataPath)
	if err != nil {
		fmt.Println("[DB] failed querying old data")
		contents = []byte("[]")
	}

	var users []models.User
	err = json.Unmarshal(contents, &users)
	if err != nil {
		fmt.Println("[DB] failed Unmarshaling when updating DB")
		users = make([]models.User, 0)
	}
	users = append(users, data)

	updatedContents, err := json.Marshal(users)
	if err != nil {
		fmt.Println("[DB] failed marshaling when updating DB")
		return err
	}

	if err := os.WriteFile(utils.UserDataPath, updatedContents, 0644); err != nil {
		fmt.Println("[DB] failed updating DB")
		return err
	}

	d.UserData = users

	fmt.Println("[DB] successfully created file")
	return nil
}
