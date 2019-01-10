package main

import "fmt"


type Name interface {
	Name() string
}

type A struct {
}

func (self A) say() {
	println(self.Name())
}

func (self A) sayReal(child Name) {
	fmt.Println(child.Name())
}

func (self A) Name() string {
	return "I'm A"
}

type B struct {
	A
}

func (self B) Name() string {
	return "I'm B"
}




type C struct {
	A
}

type Eatable interface {
	Eat()
}
type Animal struct {
}
func (a *Animal) Eat() {
	println("Animal eat")
}
type Cat struct {
	Animal
}
func (c *Cat) Eat() {
	println("Cat eat")
}

func main() {
	//b := B{}
	//b.say()         //I'm A
	//b.sayReal(b)    //I'm B
	//
	//c := C{}
	//c.say()         //I'm A
	//b.sayReal(b)    //I'm A

	var eatable Eatable
	eatable = &Animal{}
	eatable.Eat()
	eatable = &Cat{}
	eatable.Eat()
}