package util

import (
	"math/rand"
	"strconv"
)

var firstPart = []string{"Red", "Green", "Blue", "Orange", "Black", "White", "Yellow", "Purple", "Pink", "Brown", "Grey", "Silver", "Gold", "Cyan", "Magenta", "Maroon", "Navy", "Olive", "Teal", "Violet", "Lime", "Coral", "Aqua", "Azure", "Beige", "Bisque"}

var secondPart = []string{"Apple", "Pear", "Banana", "Kiwi", "Grape", "Grapefruit", "Lemon", "Lime", "Mango", "Melon", "Orange", "Peach", "Pear", "Pineapple", "Plum", "Raspberry", "Strawberry", "Tomato"}

func GenerateRandomName() string {
	var rand1, rand2 = rand.Intn(len(firstPart)), rand.Intn(len(secondPart))
	return firstPart[rand1] + secondPart[rand2] + strconv.Itoa(rand.Intn(100))
}
