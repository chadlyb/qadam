# qadam

Tools to localize Mise Qadam.

First, run `qdecomp.go` on the `TEXTS.FIL`:

`go run qdecomp.go > qadam.txt`

Then, edit the text file to localize (you can use `;` to comment out old text if you wish to keep it around.) `\n` is a newline character.

Then, run `qcompile.go` to write `TEXTS_NEW.FIL`:

`go run qcompile.go qadam.txt`

If you didn't change `qadam.txt` this should be bytewise identical to the original file!

Copy this file atop the old `TEXTS.FIL` to use it.

Finally, if the size of `TEXTS.FIL` has changed, a particular value will need to be patched in the `GAME.EXE`

`go run fixgame.go <SIZE OF TEXTS.FIL>`

This will generate `GAME2.EXE` (copy atop `GAME.EXE` to use) and you are good to go!
