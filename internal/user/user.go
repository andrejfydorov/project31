package user

import (
	"fmt"
)

type User struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Friends []*User `json:"friends"`
}

func (u User) ToString() string {
	res := fmt.Sprintf("Id is %d, name is %s and age is %d\n", u.Id, u.Name, u.Age)

	res += fmt.Sprintf("friends:\n")
	for _, friend := range u.Friends {
		res += fmt.Sprintf("friend id is %d, name is %s and age is %d\n", friend.Id, friend.Name, friend.Age)
	}

	return res
}

func (u User) ToStringShort() string {
	res := fmt.Sprintf("Id is %d, name is %s and age is %d\n", u.Id, u.Name, u.Age)
	return res
}

func (u User) GetFriends() string {
	var res string
	for _, this := range u.Friends {
		res += this.ToString() + "\n"
	}
	return res
}

func Remove(u *User, i int) {
	u.Friends[i] = u.Friends[len(u.Friends)-1]
	u.Friends = u.Friends[:len(u.Friends)-1]
}
