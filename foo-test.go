package main

import "./foo"

func main() {
  var a foo.SpammerStruct //init to zero
  //does not work: "implicit assignment"
  //a.Spam()
  
  //does not work: "implicit assignment"
  //ptrA := &a
  //ptrA.Spam()
  
  //does not work either: "implicit assignment"
  //var interfaceA foo.SpammerInterface = a
  //interfaceA.Spam()
  
  //works:
  var interfaceFromPtrA foo.SpammerInterface = &a
  interfaceFromPtrA.Spam()
  
  //so this is the short version:
  foo.SpammerInterface(&a).Spam()
  
  //sidenote: this doesn't work either
  //"cannot take the address of foo.SpammerInterface(a)"
  //var interfacePtrA *foo.SpammerInterface = &(foo.SpammerInterface(a))
  //so how would I assign a *foo.SpammerInterface
}

/*
func main() {
  var a foo.Spammer = foo.InitSpammer(1)
  a.Spam() //prints "1"
  
  if b, ok := a.(foo.Foo); ok {
    b.SetHidden(2)
    //works:
    var bAsSpammer foo.Spammer = &b
    bAsSpammer.Spam() //prints "2"
    //also works:
    foo.Spammer(&b).Spam() //prints "2"
    //why do I need to take the pointer here?
    //does not work: (&b).Spam()
    
    var pb = &b
    pb.SetHidden(3)
    //does not work: "implicit assignment"
    //pb.Spam()
    //works:
    foo.Spammer(pb).Spam() //prints "3"
    foo.Spammer(&b).Spam() //prints "3"
    a.Spam() //still prints "1"

    //how do I make a variable of type HiddenSetter?
    //does not work:
    var hs foo.SpammerSetter
    hs = &b
    hs.SetHidden(4)
    hs.Spam()
    bAsSpammer.Spam()
  }
}*/