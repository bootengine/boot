# Boot

Boot is the new generation project bootstrapper.

### 🙏 It's easy to use

It relies on "configuration files" (using usual configuration language like yaml) to describe how to generate your new project.
Those files a reusable and we did our best to reduce the complexity.


### 🔋 It's battery included

Boot comes with some default pre-installed module giving you almost all you need to start right now.

It comes with :
- a [git plugin] so that you can `init`, `add`, `commit` and `push` your code as soon as possible.
- a [filer plugin](https://github.com/bootengine/boot-filer-plugin) developped in zig so that you can quickly create a predefined folder structure.
- a [template_engine plugin] based on [jinja](https://jinja.palletsprojects.com/en/stable/), a simple yet powerful python template engine.

Boot defines 4 types of modules :
- vcs -> will perform Version Source Control actions (ex: `git init`, ....)
- filer -> will perform filesystem related action (ex: create files, create folders, ...)
- template_engine -> will render a template using a template engine (ex: ejs, pug, jinja, ...)
- cmd -> will execute language/framework action (ex: `npm install -s express`, `go mod init my-super-project`, ...)

It comes with a default module for 3 of the 4 differents types.

The only missing part is the one where we couldn't guess what would be the perfect fit : the module for you beloved language/framework.

### 🛠️ It's highly customizable

Boot module's system should support any language/framework/technology as it is based on [extism](https://extism.org/).
This will give you the possibility to create modules using almost any language (as long as you can compile it to wasm/wasi).

Moreover, the 3 default modules we provide are totally removable. You can rewrite them if you want. 

## Installation

## Usage

## License

[GPL V3.0](https://choosealicense.com/licenses/gpl-3.0/)
