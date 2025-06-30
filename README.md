# QADAM - Mise Quadam Localization Tools

Tools to extract, localize, and rebuild Mise Quadam game files. This project provides a complete workflow for translating the Czech game into other languages.

## Building (if not using precompiled package)

### Prerequisites

- [Go](https://go.dev/doc/install) 1.19 or later
- Make, if using make.

### Building

The easiest way to build the tools is using the provided Makefile:

```bash
# Build for your current platform
make build

# Build for all platforms (Windows, macOS, Linux)
make build-all

# Clean build artifacts
make clean
```

Alternatively, you can build manually:

```bash
# Build extract tool
go build -o extract ./cmd/extract

# Build build tool  
go build -o build ./cmd/build
```

## Usage

1. **Extract strings from original game files:**
   ```bash
   ./extract <path-to-original-game-folder>
   ```

   - In Windows Explorer, this can be achieved by dragging the original game folder onto the extract EXE.

   This creates an `extracted` folder containing:
   - `og/` - Copy of original files
   - `texts.txt` - Main game text (section 0)
   - `resource.txt` - Inventory item strings (section 11)
   - `game_exe.txt` - Strings from GAME.EXE
   - `install_exe.txt` - Strings from INSTALL.EXE

2. **Edit the extracted text files:**
   - `texts.txt` - Main localization file
     - Use `;` for comments
     - `\n` for newlines
     - You can change string lengths (within memory constraints)
     - Leave unedited strings unchanged
   - `resource.txt` - Inventory items
     - Only edit section 11 strings
     - Leave other sections unchanged
   - `game_exe.txt` - Executable strings
     - Delete lines containing non-human-readable strings for clarity--they will be unchanged if you do this.
     - Don't modify "garbage characters" at beginning of strings! This is probably important non-string data.
     - Don't make strings longer than original
   - `install_exe.txt` - Installer strings
     - Same rules as game_exe.txt

3. **Build localized game:**
   ```bash
   ./build <path-to-extracted-folder>
   ```

   - In Windows, this can be achieved by dragging the (possibly edited) extracted folder onto the build EXE.

   This creates a `built` folder with the localized game files.

## Testing

### Round-trip Test

To verify that the extraction and build process works correctly, you can run a round-trip test:

```bash
make roundtrip-test
```

This test:
1. Extracts strings from the `src/` directory to `temp_extract/`
2. Builds the game from the extracted files to `temp_build/`
3. Compares the rebuilt files with the original `src/` files
4. Cleans up the temporary directories

If the test passes, it confirms that the extraction and build process preserves all data correctly.

## File Format Details

### texts.txt and resource.txt (from .FIL files)
- Format: `[hex-data] "string content"`
- Example: `[01 02 03 04 05] "Hello world"`
- Hex data (probably?) contains: identifier, color, screen location, and flags
- Use `;` for comments
- This file is used to reconstruct the FIL file from scratch, so don't delete anything (other than editing inside strings)!

### game_exe.txt and install_exe.txt (from executables)
- Format: `<beginoffset>-<endoffset>:"string content"`
- Example: `00001236-0000123f: "New Game"`
- The string may begin with "garbage" characters -- **Leave these characters completely unchanged** - they contain important game data
- Use ';' for comments
- This is a patch file, so you can delete lines that you don't want to patch and the underlying EXE won't be changed.

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run specific test suites
make test-language-model
make test-extraction

# Run tests with coverage
make test-coverage

# Discussed above--assumes game is in `src/` subfolder.
make roundtrip-test
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run benchmarks
make benchmark
```

### Creating Releases

```bash
# Create release packages for all platforms
make release
```

This creates:
- `qadam-<version>-windows-amd64.zip`
- `qadam-<version>-linux-amd64.tar.gz`  
- `qadam-<version>-darwin-amd64.tar.gz`

### Debug Mode

Enable verbose output for troubleshooting:

```bash
./extract -v <source-directory>
```

