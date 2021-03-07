package main

import (
	"context"
	"errors"
	socketio "github.com/googollee/go-socket.io"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type client struct {
	messages chan message // msg // talking to players
	commands chan message // cmd // commands prefixed with /
	priority chan message // pri // priority (respond to server without / prefix)
	session  session
	conn     socketio.Conn
}

func (c *client) sendPriority(msg string) {
	server.BroadcastToRoom("/", c.conn.ID(), "pri", newMessage(serverPlayer, msg))
}

func handleClient(c *client) error {

	isHandled := false

	for isHandled {
		c.sendPriority("Welcome to Perpetuan\n\nAre you here to <sign up> or to <log in>\n(Type Response)")

		reply, ok := <- c.priority
		if !ok {
			return errors.New("error reading from priority channel")
		}

		choice := reply.Text

		var user userData
		var playerData player

		if choice == "sign up" || choice == "signin" {
			valid := false
			for !valid {
				c.sendPriority("Enter a username:")
				usernameReply, ok := <- c.priority
				if !ok {
					return errors.New("error reading from priority channel")
				}
				if UsernameValid(usernameReply.Text) && UsernameUnique(usernameReply.Text) {
					user.Username = usernameReply.Text
					playerData.Username = usernameReply.Text
					c.sendPriority("Enter your email:")
					emailReply, ok := <- c.priority

					if !ok {
						return errors.New("error reading from priority channel")
					}

					if EmailValid(emailReply.Text) && EmailUnique(emailReply.Text) {
						user.Email = emailReply.Text
						c.sendPriority("Enter your password:")
						passwordReply, ok := <- c.priority

						if !ok {
							return errors.New("error reading from priority channel")
						}

						c.sendPriority("Re-enter your password:")
						passwordReReply, ok := <- c.priority

						if !ok {
							return errors.New("error reading from priority channel")
						}

						if PasswordValid(passwordReply.Text) && passwordReply.Text == passwordReReply.Text {
							user.Password = passwordReply.Text

							settingList, ok := updateSettingList([]setting{})
							settings := ""

							if !ok {
								return errors.New("error updating setting list")
							}

							for i := range settingList {
								settings += "[" + settingList[i].key + "] "
							}

							c.sendPriority("Choose your starting region from the list: " + settings)
							settingReply, ok := <- c.priority

							if !ok {
								return errors.New("error reading from priority channel")
							}
							if SettingValid(settingReply.Text) {
								playerData.Room = settingReply.Text
								playerData.Pos = position{
									Pos:     [2]int{0,0},
									Setting: settingReply.Text,
								}
								user.Player = playerData

								err := registerUser(user)
								if err != nil {
									return errors.New("error reading from priority channel")
								}

								syncPlayer(c)
								valid = true
								isHandled = true
								break
							} else {
								c.sendPriority("Not a valid region choice")
							}

						} else {
							c.sendPriority("Passwords do not match")
						}

					} else {
						c.sendPriority("Email invalid, its either already registered or the syntax is wrong or the domain doesn't exist")
					}

				} else {
					c.sendPriority("Username invalid, the only symbols you can use in a username and - and _, \nUsername might already be taken if that's not it")
				}
				c.sendPriority("Type exit at any time to return to options")
			}



		} else if choice == "log in" || choice == "login" {

			exit := false

			for !exit {
				c.sendPriority("Enter username or email")

				username, ok := <- c.priority
				if !ok {
					return errors.New("error reading from priority channel")
				}

				c.sendPriority("Enter password")

				password, ok := <- c.priority
				if !ok {
					return errors.New("error reading from priority channel")
				}

				ok = loginUser(username.Text, password.Text)

				if !ok {
					c.sendPriority("Invalid Credentials")
				} else {
					exit = true
					c.sendPriority("User logged in")
					isHandled = true
					break
				}
				c.sendPriority("Type exit at any time to return to options")
			}

		} else if choice == "exit" {
			continue
		} else {
			c.sendPriority("Type exit at any time to return to options")
		}
	}
	return nil
}

func registerUser(user userData) error {
	hashedPass, err := hashPassword(user.Password)

	if err != nil {
		log.Println("Hash error: ", err)
		return err
	}

	_, err = connect.userCollection.InsertOne(context.Background(), bson.M{
		"_id": user.Player.ID,
		"Username": user.Username,
		"Email": user.Email,
		"Password": hashedPass,
		"Player": user.Player,
		"Online": true,
	})

	if err != nil {
		log.Println("Insert error: ", err)
		return err
	} else {
		log.Println("User: ", user.Username, " created")
		return nil
	}
}

func loginUser(username string, password string) bool {
	isFound := false

	userRes := connect.userCollection.FindOne(context.Background(), bson.M{"Username": username})
	if userRes.Err() != nil {
		isFound = false
		emailRes := connect.userCollection.FindOne(context.Background(), bson.M{"Email": username})
		if emailRes.Err() != nil {
			isFound = false
		} else {
			isFound = true
			BSONData := struct {
				ID    string  `bson:"_id"`
				Username string `bson:"Username"`
				Email string `bson:"Email"`
				Password []byte `bson:"Password"`
				Player player `bson:"Player"`
				Online bool `bson:"Online"`
			}{}

			decodeError := emailRes.Decode(&BSONData)

			if decodeError != nil {
				log.Println("Decode error: ", decodeError)
				return false
			}

			if !compareHash(BSONData.Password, password) {
				return false
			}
		}
	} else {
		isFound = true
		BSONData := struct {
			ID    string  `bson:"_id"`
			Username string `bson:"Username"`
			Email string `bson:"Email"`
			Password []byte `bson:"Password"`
			Player player `bson:"Player"`
			Online bool `bson:"Online"`
		}{}

		decodeError := userRes.Decode(&BSONData)

		if decodeError != nil {
			log.Println("Decode error: ", decodeError)
			return false
		}

		if !compareHash(BSONData.Password, password) {
			return false
		}
	}

	return isFound

}

func hashPassword(passwordRaw string) ([]byte, error) {
	password := []byte(passwordRaw)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func compareHash(hash []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return false
	}
	return true
}

func syncPlayer(c *client) error {
	res := connect.userCollection.FindOne(context.Background(), bson.M{"_id": c.session.Player.ID})
	if res.Err() != nil {
		return res.Err()
	}

	BSONData := struct {
		ID    string  `bson:"_id"`
		Username string `bson:"Username"`
		Email string `bson:"Email"`
		Password []byte `bson:"Password"`
		Player player `bson:"Player"`
	}{}

	decodeError := res.Decode(&BSONData)

	if decodeError != nil {
		log.Println("Decode error: ", decodeError)
		return decodeError
	}

	c.session.Player = BSONData.Player
	return nil
}