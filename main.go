package main

import (
	"bufio"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string
	Completed   bool
	Description string
	UserID      uint
	User        User
}

func (t Todo) String() string {
	return fmt.Sprintf(
		"\n%s> Description: %s> Completed: %t\n> CreatedAt: %s\n> LastUpdated: %s",
		t.Title, t.Description, t.Completed, t.CreatedAt, t.UpdatedAt)
}

func (Todo) TableName() string {
	return "go_todo_todo"
}

type User struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Password string
	Todos    []Todo
}

func (User) TableName() string {
	return "go_todo_user"
}

func getAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func errorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func getTodosByUser(db *gorm.DB, user *User) ([]Todo, error) {
	var todos []Todo
	result := db.Where("user_id = ?", user.ID).Find(&todos)
	return todos, result.Error
}

func main() {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB")), &gorm.Config{})
	errorPanic(err)

	user := &User{}

	err = db.AutoMigrate(&User{}, &Todo{})
	errorPanic(err)

	reader := bufio.NewReader(os.Stdin)

main:
	for {
	inner:
		for {

			if user.Name == "" {
				fmt.Println(`
=============
1. Login
2. Register

0. Exit
=============`)

				var input string
				_, err = fmt.Scanln(&input)
				errorPanic(err)

				switch input {
				default: // Exit
					break main
				case "1": // Login
					var name string
					fmt.Println("\nEnter your name:")
					_, err = fmt.Scanln(&name)
					errorPanic(err)

					var password string
					fmt.Println("\nEnter your password:")
					_, err = fmt.Scanln(&password)
					errorPanic(err)

					var users []User
					result := db.Where("name = ? AND password = ?", name, password).Find(&users)
					if result.Error != nil {
						fmt.Println("\nError while logging in!")
						break inner
					}

					if len(users) > 0 {
						user = &users[0]
						fmt.Println("\nSuccessfully logged in!")
						break inner
					}

					fmt.Println("\nInvalid credentials! Try again!")
					break inner
				case "2": // Register

					fmt.Println("\nEnter your name (0 - exit):")
					var name string
					_, err = fmt.Scanln(&name)
					errorPanic(err)

					if name == "0" {
						break inner
					}

					var users []User
					result := db.Where("name = ?", name).Find(&users)
					if result.Error != nil {
						fmt.Println("\nError while registering!")
						break inner
					}
					if len(users) > 0 {
						fmt.Println("\nThis name already exists!")
						break inner
					}

					var password string
					var confirmPassword string
				password:
					for {
						fmt.Println("\nEnter your password (0 - exit):")
						_, err = fmt.Scanln(&password)
						errorPanic(err)

						if password == "0" {
							break inner
						}

						fmt.Println("\nConfirm your password:")
						_, err = fmt.Scanln(&confirmPassword)
						errorPanic(err)

						if password != confirmPassword {
							fmt.Println("\nPasswords don't match!")
						} else {
							break password
						}
					}
					user = &User{Name: name, Password: password}
					db.Create(user)
					fmt.Println("\nSuccessfully registered!")
					break inner
				}
			} else {
				fmt.Println(`
===============
1. Logout
2. List todos
3. Add todo
4. Delete todo
5. Update todo

6. Update user

0. Exit
===============`)
				var input string
				_, err := fmt.Scanln(&input)
				errorPanic(err)

				switch input {
				default: // Exit
					break main
				case "1": // Logout
					user = &User{}
					break inner
				case "2": // List todos
					todos, err := getTodosByUser(db, user)
					errorPanic(err)

					if len(todos) == 0 {
						fmt.Println("\nYou don't have any todos!")
					}

					for _, todo := range todos {
						fmt.Println(todo)
					}
				case "3": // Add tudu
					fmt.Println("\nEnter todo title:")
					title, err := reader.ReadString('\n')
					errorPanic(err)

					fmt.Println("\nEnter todo description:")
					description, err := reader.ReadString('\n')
					errorPanic(err)

					fmt.Println("\nTodo successfully added!")
					db.Create(
						&Todo{Title: title, Completed: false, Description: description, UserID: user.ID, User: *user})
					break inner
				}
			}
		}
	}
}
