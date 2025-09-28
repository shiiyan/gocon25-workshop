package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string
	Age  int
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	return json.Unmarshal(data, &aux)
}

func Unmarshal[T any, PT Unmarshaller[T]](data []byte) (T, error) {
	var v T
	err := PT(&v).UnmarshalJSON(data)
	return v, err
}

type Unmarshaller[T any] interface {
	*T
	json.Unmarshaler
}

func main() {
	data := []byte(`{"Name": "Alice", "Age": 30}`)
	user, err := Unmarshal[User](data)
	if err != nil {
		panic(err)
	}
	fmt.Println(user.Name, user.Age) // Output: Alice 30
}
