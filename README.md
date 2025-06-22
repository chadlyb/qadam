# qadam

Tools to localize Mise Quadam.

First, run `qdecomp.go` on the `TEXTS.FIL`:

`go run qdecomp.go src/TEXTS.FIL > quadam_texts.txt`

Then, edit the text file to localize (you can use `;` to comment out old text if you wish to keep it around.) `\n` is a newline character. You may increase or decrease the length of strings (as long as overall you still fit in memory.)

Then, run `qcompile.go` to write `TEXTS.FIL`:

`go run qcompile.go quadam_texts.txt dest/TEXTS.FIL`

If you didn't change `quadam_texts.txt` this should be bytewise identical to the original file!

Do the same steps for `RESOURCE.FIL` (replacing "TEXTS" with "RESOURCE" in above instructions) to localize section 11 of that file, (containing inventory item description strings.)
Note that there will be a lot of weird nonsense strings and such in the file in the other sections--leave them alone, and the file should compile back to identical data.

`go run qdecomp.go src/RESOURCE.FIL > quadam_resource.txt`

`go run qcompile.go quadam_resource.txt dest/RESOURCE.FIL`

Finally, if the size of `TEXTS.FIL` or `RESOURCE.FIL` has changed, the new size(s) will need to be patched in the `GAME.EXE`:

`go run fixgame.go src/GAME.EXE dest/GAME.EXE <SIZE OF dest/TEXTS.FIL> <SIZE OF dest/RESOURCE.FIL>`


To patch strings in the EXE, first:

`go run qgetstrings.go src/GAME.EXE > exe_strings.txt`

Then, manually remove all the non-human-readable lines from `exe_strings.txt`.
You can edit the text to localize BUT unlike the FIL files, each new string should be the same size or smaller than the original!
(also if there are weird characters at the beginning of the string, keep them intact, since it's probably important stuff from before the actual string)

Then:

`go run qpatchstrings.go dest/GAME.EXE dest/GAME.EXE exe_strings.txt`

To patch the strings (the above command will change `dest/GAME.EXE` in place, I'm assuming that is the output of `fixgame` above--NOTE that you will have to run `fixgame` and then also `qpatchstrings` on its output to get the final EXE, or the other way around.)