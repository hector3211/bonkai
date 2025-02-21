package pkg

import (
	"encoding/json"
	"fmt"
	"my-chi/models"
	"os"
)

const UserDataPath = "/tmp/uer-data.json"

type DB struct {
	UserData []models.User
}

func PingDB() error {
	_, err := os.Stat(UserDataPath)
	if err != nil {
		fmt.Println("[File] ping failed, no such file exists")
		return err
	}
	return nil
}

func NewDB() (*DB, error) {
	err := PingDB()
	if err != nil {
		_, err := os.Create(UserDataPath)
		if err != nil {
			fmt.Println("[File] failed creating user-data file")
			return nil, err
		}
		contents, err := os.ReadFile(UserDataPath)
		if err != nil {
			fmt.Println("[File] no such file exists")
			os.Exit(0)
		}

		data := []models.User{}
		err = json.Unmarshal(contents, &data)

		fmt.Println("[File] successfully created file")
		return &DB{
			UserData: data,
		}, nil
	}
	return nil, err
}

func (d *DB) UpdateWithData(data []models.User) error {
	err := os.Remove(UserDataPath)
	if err != nil {
		fmt.Println("[DB] failed removing old data")
		return err
	}

	contents, err := json.Marshal(data)
	err = os.WriteFile(UserDataPath, contents, 0644)
	if err != nil {
		fmt.Println("[DB] failed updating DB")
		return err
	}
	fmt.Println("[File] successfully created file")
	return nil
}
