# Whim
What is this repo? This is a Text editor called Whim.\
This text editor is based on [Snap Token Kilo Editor Guide](https://viewsourcecode.org/snaptoken/kilo/index.html)

## Why did I make this?
I decided that writing my own text editor would be fun.
One day I stumbled across the article and I decided to give creating a text editor a try.\
You might be asking why is it in go?
When I started this project I was learning go and I really started to like the language.
Also instead of copying the article entirely (I can't really learn that way) I decided that
trying it in go would be a better learning experience.

## What enviorments does this support?
Currently terminals that support ANSI escape codes & linux machines.

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
Currently its just a simple text editor. Check the features section below to see
what functionality Whim has.

## Credits
Once again this is based of the article [Snap Token Kilo Editor Guide](https://viewsourcecode.org/snaptoken/kilo/index.html)
Massive thanks to them for providing this amazing guide.

## Upcoming Features
* If a file is "Dirty" ask user to confirm to exit
* Soft indents
* Add support for other enviorments
* Copy & Paste
* Undo & Redo
* Add file highligting for languages:
    * Golang
    * Javascript
* Add custom Configuration
* Add Multiple buffers

## Current Defects
* Traversing a file with tabs the cursor moves inconsistently

## Tech Debt
* Fix/Remove Remaining TODO's

## Added Features
Below are the added feature after completing the article
* Line Numbers (Toggle, added relative line numbers as well)
* Modal editor (Normal, Insert, Command)

