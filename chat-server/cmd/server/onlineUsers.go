package main

import "fmt"

func (app *app) addToOnlineUsers(userId string) {
	err := app.redisConnection.SAdd(ctx, onlineUsersKey, userId).Err()
	if err != nil {
		fmt.Println("Error adding user to online users")
	}
	fmt.Println("User added to online users")
}

func (app *app) countOnlineUsers() (int64, error) {
	count, err := app.redisConnection.SCard(ctx, onlineUsersKey).Result()
	if err != nil {
		fmt.Println("Error getting online users count")
		return 0, err
	}
	return count, nil
}

func (app *app) removeFromOnlineUsers(userId string) {
	err := app.redisConnection.SRem(ctx, onlineUsersKey, userId).Err()
	if err != nil {
		fmt.Println("Error removing user from online users")
	}
	fmt.Println("User removed from online users")
}
