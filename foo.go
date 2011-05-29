package foo;

import "fmt"

type SpammerInterface interface {
  Spam()
}

type SpammerStruct struct {
  hidden int
}

func (s SpammerStruct) Spam() {
  fmt.Println("Spam", s.hidden)
}

/*
type SpammerSetter interface {
  SpammerInterface
  SetHidden(b int)
}

func (f *Foo) SetHidden(x int) {
  f.hidden = x
}

func InitSpammer(x int) Spammer {
  return Foo{x}
}*/