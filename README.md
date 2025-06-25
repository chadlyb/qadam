# qadam

Tools to localize Mise Quadam.

To use, first build the exes (assuming you've installed Go from e.g., https://go.dev/doc/install)
`go build ./cmd/build`
`go build ./cmd/extract`

Then you can run `extract` (or `extract.exe` in Windows) and pass it in a folder containing the original game as an argument
(in Windows, you can do this easily by dragging the folder onto the `extract.exe`.)

This will create an `extracted` folder alongside the original folder. It will contain:
 - `og` This subdirectory contains a copy of the original files
 - `texts.txt` This is the main file to localize. 
   - (you can use `;` to comment out old text if you wish to keep it around.)
   - `\n` is a newline character. 
   - You MAY increase or decrease the length of strings (as long as overall you still fit in memory.)
   - Leave any strings you don't edit alone. 
   - I believe the hex characters before the strings are an identifier, a color, a screen location, and 32 (or 1 for he first item in a group for some reason.) 
 - `resource.txt` Similar to just above, but contains inventory item strings in section 11.
   - Note that there will be a LOT of weird nonsense strings and such in the file in the other sections--*leave them alone*, and the file should compile back to identical data.
 - `game_exe.txt` Contains strings from the EXE. 
   - This is a bit different than the first two--you may (and probably should for your sanity) delete any strings in this file which aren't human readable before editing. 
   - You should NOT get rid of any weird characters before the strings you edit, since those probably contain important non-string data. 
   - Also, unlike the above two, you may NOT make these strings longer than the original (the builder will complain if you do.)
 - `install_exe.txt` Similar to just above, but for the installer.

Then, when you want to test your changes to the files in `extracted`, you can run `build` (or `build.exe` in windows) and pass it in the `extracted` folder as an argument (or, again, in Windows, you can drag the `extracted` folder onto `build.exe`.)

This will generate the patched copy of the game in `built`.

If you run `extract` and then `build` immediately, the output in the `built` folder should exactly match the original files.

Good luck!
