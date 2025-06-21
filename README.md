# qadam

Tools to localize Mise Quadam.

First, run `qdecomp.go` on the `TEXTS.FIL`:

`go run qdecomp.go src/TEXTS.FIL > quadam_texts.txt`

Then, edit the text file to localize (you can use `;` to comment out old text if you wish to keep it around.) `\n` is a newline character.

Then, run `qcompile.go` to write `TEXTS.FIL`:

`go run qcompile.go quadam_texts.txt dest/TEXTS.FIL`

If you didn't change `quadam_texts.txt` this should be bytewise identical to the original file!

Do the same steps for `RESOURCE.FIL` (replacing "TEXTS" with "RESOURCE" in above instructions) to localize section 11 of that file, (containing inventory item description strings.)
Note that there will be a lot of weird nonsense strings and such in the file in the other sections--leave them alone, and the file should compile back to identical data.

`go run qdecomp.go src/RESOURCE.FIL > quadam_resource.txt`

`go run qcompile.go quadam_resource.txt dest/RESOURCE.FIL`

Finally, if the size of `TEXTS.FIL` or `RESOURCE.FIL` has changed, the new size(s) will need to be patched in the `GAME.EXE`:

`go run fixgame.go src/GAME.EXE dest/GAME.EXE <SIZE OF dest/TEXTS.FIL> <SIZE OF dest/RESOURCE.FIL>`

