# PA CHECKER

## License

This project is using the `ISC` License.

## Contributors

* Asavoae Cosmin-Stefan
* Gatej Stefan-Alexandru
* Neamu Ciprian-Valentin
* Potop Horia-Ioan

## Dependencies

* `valgrind`
* `cppcheck`
* `git`

### Quick install

```bash
# Update the package repositories
sudo apt update

# Install valgrind
sudo apt install valgrind

# Install cppcheck
sudo apt install cppcheck

# Install git (if not already installed)
sudo apt install git
```

## Features

- [x] Parallel test running

- [x] Configuration
  - [x] Configurable tests
  - [x] Configurable modules
  - [x] User configuration
  - [x] Macros

- [x] Modules
  - [x] Module dependency checks
  - [x] Diff module
  - [x] Memory module _(valgrind backend)_
  - [x] Style module _(cppcheck backend)_
  - [x] Commit module _(git backend)_

- [x] Interface
  - [x] Basic - full module dump
  - [x] Interactive
    - [x] Live reload
    - [x] Module output visualization
      - [x] Side-by-side diff visualization
      - [x] Memory leak information

  
- [x] OS Compatibility
  - [x] `Linux / WSL` - full support
  - [ ] `OSX` - partial support _(no backend for the memory module)_
  - [ ] `Windows` - partial support _(no backend for the memory module)_

## Overview

### Running the checker

#### Basic
```bash
./checker
```

#### Interactive
```bash
./checker -i
```

### Navigating the interactive interface

* Use the `arrow keys` to navigate around
* Press `TAB` to switch between navigation and current section
* Press `ESC` to exit a fullscreen page
* Press `ESC` or `Ctrl+C` while on the main page to exit the program
* Press `~` to trigger a test run _(or modify the executable)_
* `Mouse` should be fully supported

### Configuration

Inside `config.json` or the `Options` tab you can modify the following:

* `Executable Path` - the executable that will be used to run the tests
* `Source Path` - the project root directory
* `Input Path` - the directory containing the input files
* `Output Path` - the directory where the test output will be stored
* `Ref Path` - the directory containing the reference files
* `Forward Path` - the directory where the `stdout` & `stderr` of each test will be stored
* `Valgrind` - whether to run the tests using valgrind or not _(disable for faster iteration)_
* `Tutorial` - display the tutorial again _(disabled afterward)_

### Interface screenshots
<div style="text-align: center;">

![refs example](./res/doc/ref_example.png)
<br>
`Refs tab example`

</div>

<br>


<div style="text-align: center;">

![diff example](./res/doc/diff_example.png)
<br>
`Diff visualization`

</div>

<br>

<div style="text-align: center;">

![options example](./res/doc/style_example.png)
<br>
`Style tab example`

</div>

<br>

<div style="text-align: center;">

![options example](./res/doc/options_example.png)
<br>
`Options tab example`

</div>


---

<br>

### FAQ

1. What does the live reload feature do?
> The live reload feature watches for any changes made to the provided executable and triggers a new test run when it's modified.

<br>

2. One or more modules went into panic! What do I do now?
> One or more modules might go into panic from various reasons. The common ones are:
> * The executable was deleted or an invalid path was provided
> * The checker doesn't have read / write access to one or more of the provided paths
> * The config was set up incorrectly
> 
> Simply look for any of these issues. After you're sure that the problem is fixed, just relaunch the checker or trigger a new run by recompiling your code. _(or by pressing `~`)_

<br>

3. I solved the whole assignment but my score is not 100! Where did my points go!?
> Please check that all modules are enabled first. _(no `DISABLED` status)_

<br>

4. All the modules are enabled but my score is still not 100!
> Probably there are still issues to be ironed out, make sure that each module page displays no errors.

<br>

5. I previously had a score of 100 on a test and although I didn't modify any of the code responsible for the test, my score is now lower.
> The memory and git modules run on your entire code. If you recently added code that produces leaks, for example, this will affect the score on all tests. Make sure to correct all memory and git issues before you submit your program for a specific task.

---

<br>

## Contributing

> This project is still a work-in-progress and any contributions are welcome!

### Project Structure

* root
  * `bin` - Windows & Linux compiled binaries 
  * `res` - project resources: config files
  * `src` - project source files
    * `checker-modules`
    * `display`
    * `menu`
    * `manager`
    * `utils`
  * `main.go` - project entrypoint
  * `Makefile` - use this to compile the project

### Building the checker

To build the checker simply run
```bash
make build-linux # ELF executable

make build-windows # Win32 executable

make build-macos # OSX executable
```

### Formatting

Before committing any changes, run
```bash
make vet

make lint
```

### Commit format

* Keep in mind that Andra recommended that the commits be in english.
* The commits must be signed (`git commit -s`)
* The commit messages should have the following structure

```
MODULE: <concise title>

*detailed description* (around 75 characters per line)
```
> Example commit message
> ```
> ref: Added order checks
> 
> Lorem ipsum odor amet, consectetuer adipiscing elit. Neque magna platea
> ornare a maecenas aptent tincidunt. Tellus dolor maecenas congue pharetra
> leo himenaeos dis curabitur. Accumsan venenatis eget ipsum enim montes
> volutpat quisque. Diam finibus leo mattis fames efficitur.
> ```
