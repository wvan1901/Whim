# Whim
What is this repo? This is a Linux(I think... It works on my linux laptop) Text editor called Whim.\
This text editor is based on [Snap Token Kilo Editor Guide](https://viewsourcecode.org/snaptoken/kilo/index.html)
I tried to use little to no external Libaries

## Why did I make this?
Why not. I decided that writing my own text editor would be fun.
I tried to make a text editor with out any guidence but I ended with many bugs implementing
a rope data structure so I dropped it. One day I stumbled across the article and 
I decided to give a text editor a try again. You might be asking why is it in go?
I wrote it in go because I've been learning go and I really started to like the language.
Also instead if copying the article entirely (I can't really learn that way) I decided that
trying it in go would make me learn better.

## What enviorments does this support?
Currently only Linux, it works on my ubuntu Linux laptop. Hopefully I can add support to other enviorments

## How install and run Whim?
Currently this project is still in development phase. So to run this project 
Clone the repo and enter one of the commands:
```bash
# To edit a file
go run main.go aFile.txt
# To open whim
go run main.go
```

## How to use Whim?
Currently its just a simple text editor. As of writing this its not a super useful
editor since some keys are reserved as actions so you wont be able to edit a file
properly but hopely I will fix this soon when I implement a modes!

## Credits
Once again this is based of the article [Snap Token Kilo Editor Guide](https://viewsourcecode.org/snaptoken/kilo/index.html)
So massive thanks to them for providing this amazing guide.

## Upcoming Feature
* Add Line Numbers
* Add modal editing
* If a file is "Dirty" as user to confirm to exit
* Soft indents
* Add support for other enviorments
* Copy & Paste
* Add file highligting for languages:
    * Golang
    * Javascript
* Add custom Configuration
* Add Multiple buffers

## Current Defects
* Traversing a file with tabs the cursor moves inconsistently
