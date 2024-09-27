# Whim
What is this repo? This is a Text editor called Whim.\
This text editor is based on [Snap Token Kilo Editor Guide](https://viewsourcecode.org/snaptoken/kilo/index.html)

## Why did I make this?
I decided that writing my own text editor would be fun.
One day I stumbled across the article, and I decided to give creating a text editor a try.\
You might be asking, Why is it in Go?
When I started this project, I was learning Go, and I really started to like the language.
Instead of copying the article entirely (I can't really learn that way) I decided that
trying it in Go would be a better learning experience.

## What environments does this support?
Currently, Linux based terminals

## How to install and run Whim?
Currently, this project is still in development. To run this project
Clone the repo and enter one of the commands:
```bash
# To edit a file
go run main.go aFile.txt
# To open whim
go run main.go
```

## How to use Whim?
Currently, it's just a simple text editor. Check out the features section below to see
what functionality Whim has.

## Credits
Once again, this is based on the article [Snap Token Kilo Editor Guide](https://viewsourcecode.org/snaptoken/kilo/index.html)
Massive thanks to them for providing this amazing guide.

## Upcoming Features
* Copy & Paste
* Soft indents
* Add support for other environments
* Undo & Redo
* Add file highlighting for languages:
    * Golang
    * JavaScript
* Add custom Configuration
* Add Multiple buffers

## Current Defects
* When Traversing a file with tabs, the cursor moves inconsistently

## Tech Debt
* Fix/Remove Remaining TODO's

## Added Features
Below are the added features after completing the article
* Line Numbers (Toggle, added relative line numbers as well)
* Modal editor (Normal, Insert, Command)
* If a file is "Dirty" ask the user to confirm to exit (Normal Mode Only)

