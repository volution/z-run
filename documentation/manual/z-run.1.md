

# z-run -- lightweight scripts library tool




## DESCRIPTION

`z-run` is a lightweight and portable tool that allows one to create and execute a library of scripts.

A "script", or "scriptlet" how it's named in the case of `z-run`,
is nothing more than an usual Bash, Python, Ruby, NodeJS, or any other interpreted language, script.
Basically one can just take the contents of a normal script file, and use it as the body of a scriptlet.

A "library" is just a collection of many such "scriptlets" that are bundled together in (usually) one,
or (sometimes for more complex scenarios) multiple, files.




## SYNOPSIS

### Basic modes

`z-run` <br/>
`z-run` `<scriptlet>` [ `<argument>` ... ] <br/>
`z-run` [ `<flag>` ... ] `<command>` [ `<scriptlet>` [ `<argument>` ... ] ] <br/>
`z-run` `--version` <br/>
`z-run` `--help`

### Advanced modes

`z-run` `--exec` `<library>` [ `<scriptlet>` [ `<argument>` ... ] ] <br/>
`z-run` `--ssh` [ `<flag>` ... ] `<scriptlet>` [ `<argument>` ... ] <br/>
`z-run` `--invoke` `<invoke-payload>`

### Input mode

`z-run` `--select` <br/>
`z-run` `--input` [ ... ] <br/>
`z-run` `--fzf` [ ... ]




## ARGUMENTS

`<command>`

  the name of a built-in `z-run` command;
  it is always a word (made of letters or digits), or multiple words joined by hyphen (`'-'`);
  (see bellow for details;)

`<scriptlet>`

  the label of the scriptlet within the library;
  it always starts with two colons `'::'`, then optionally followed by spaces, while the rest should match the label of an existing scriptlet;
  they never contain control characters (e.g. `'\0'`, `'\t'`, `'\n'`, etc.), but they might contain any other UTF-8 characters;
  (obviously shell quoting or escaping is required if spaces are used;)
  (e.g. `':: scriptlet'`;)

`<argument>`

  in case of commands that execute the scriptlet, the arguments passed to the scriptlet, they are never inspected or interpreted by `z-run`;
  in case of other commands, they might have a different meaning;
  (regardless the case, any `z-run` flag should come before the scriptlet;)

`<flag>`

  the name of a built-in `z-run` flag;
  it always starts with `--` and is followed by a word (made of letters or digits), or multiple words joined by a hyphen;
  if it takes a value it should be immediately followed by an `=` and the given value;
  (i.e. flag values are part of the flag argument, like `--flag=value`, and not like `--flag value`;)
  any flags must appear before the `<command>` and `<scriptlet>`;

`<rpc-target>`

  a TCP or UNIX domain socket that will be used to export or use the library over the `z-run` specific RPC;
  it is either `tcp:<ip-address>:<port>`, `tcp:<dns>:<port>`, `unix:<path>`, or `unix:@<token>`;
  (the last form is only supported on Linux, and represents an abstract UNIX domain socket;)

`<ssh-target>`

  usually `user@machine` or just `machine` that identifies the SSH target;

`<argument-0>`

  (i.e. the first argument in `main (argv)`)
  must always be the actual executable path (either absolute or relative), or alternatively `z-run` (or a few other aliases);
  other values are internally used by `z-run` itself;




## COMMANDS


### Scriptlet execution commands


`execute-scriptlet`, `execute` <br/>
`z-run` [ `<flag>` ... ] `execute-scriptlet` `<scriptlet>` [ `<argument>` ... ]

  expects a mandatory `<scriptlet>`, plus optional scriptlet `<arguments>`;
  executes the specified scriptlet, on the local machine;
  passes any specified arguments to the scriptlet;


`execute-scriptlet-ssh`, `execute-ssh`, `ssh` <br/>
`z-run` [ `<flag>` ... ] `execute-scriptlet-ssh` `<ssh-target>` `<scriptlet>` [ `<argument>` ... ]

  expects a `<ssh-target>`, a mandatory `<scriptlet>`, plus optional scriptlet `<arguments>`;
  executes the specified scriptlet, on the remote machine designated by the SSH target (usually `user@machine`);
  passes any specified arguments to the scriptlet;


`select-execute-scriptlet`, `select-execute` <br/>
`z-run` [ `<flag>` ... ] `select-execute-scriptlet` <br/>
`z-run` [ `<flag>` ... ] `select-execute-scriptlet` `<scriptlet>`

  expects an optional `<scriptlet>`, but no other arguments;
  if no scriptlet is specified, it presents a menu of all the scriptlets in the library, and allows the user to choose one;
  executes the scriptlet (either specified or selected), on the local machine;
  no arguments are passed to the scriptlet;
  usually if scriptlet is specified, it is a sub-menu that is presented instead of the entire library;


`select-execute-scriptlet-loop`, `select-execute-loop`, `loop` <br/>
`z-run` [ `<flag>` ... ] `select-execute-scriptlet-loop` <br/>
`z-run` [ `<flag>` ... ] `select-execute-scriptlet-loop` `<scriptlet>`

  expects an optional `<scriptlet>`, but no other arguments;
  similar to the `select-execute-scriptlet` command, but it will execute in a loop, presenting a menu, executing the scriptlet, pausing after the execution, and looping until the user exits the main menu;


### Scriptlet related commands


`export-scriptlet-labels`, `export-labels`, `list` <br/>
`z-run` [ `<flag>` ... ] `export-scriptlet-labels`

  expects no arguments or scriptlet;
  writes to `/dev/stdout` the labels of all the scriptlets in the library, one item per line, without the `'::'` prefix;
  (as mentioned bellow, scriptlet labels can't contain control characters, including `'\n'`;)


`export-scriptlet-body`, `export-body` <br/>
`z-run` [ `<flag>` ... ] `export-scriptlet-body` `<scriptlet>`

  expects a mandatory `<scriptlet>`, but no other arguments;
  writes to `/dev/stdout` the body of the specified scriptlet;
  (as mentioned bellow, scriptlet bodies should be UTF-8 compliant, but might not be;)


`select-export-scriptlet-label`, `select-label`, `select` <br/>
`z-run` [ `<flag>` ... ] `select-export-scriptlet-label` <br/>
`z-run` [ `<flag>` ... ] `select-export-scriptlet-label` `<scriptlet>`

  expects an optional `<scriptlet>`, but no other arguments;
  similar to the `select-execute-scriptlet` command, but it will write to `/dev/stdout` the label of the specified scriptlet, without the `'::'` prefix, followed by `'\n`';


`select-export-scriptlet-body`, `select-body` <br/>
`z-run` [ `<flag>` ... ] `select-export-scriptlet-body` <br/>
`z-run` [ `<flag>` ... ] `select-export-scriptlet-body` `<scriptlet>`

  expects an optional `<scriptlet>`, but no other arguments;
  similar to the `select-execute-scriptlet` command, but it will write to `/dev/stdout` the body of the specified scriptlet, followed by `'\n'`;


`select-export-scriptlet-label-and-body` <br>
`z-run` [ `<flag>` ... ] `select-export-scriptlet-label-and-body` <br/>
`z-run` [ `<flag>` ... ] `select-export-scriptlet-label-and-body` `<scriptlet>`

  expects an optional `<scriptlet>`, but no other arguments;
  similar to the `select-execute-scriptlet` command, but it will write to `/dev/stdout` the label of the specified scriptlet, with the `'::'` prefix, followed by `'\n'`, and then followed by the body of the specified scriptlet, followed by `'\n'`;


### Library related commands


`export-library-json` <br/>
`z-run` [ `<flag>` ... ] `export-library-json`

  expects no arguments;
  writes to `/dev/stdout` a series of JSON objects that represents the key-value store that backs the library;
  it uses pretty-printing, thus one JSON object will span over multiple lines;
  the outer serialization format (i.e. `namespace`, `key` and `value`) is unlikely to change in the future;
  the inner serialization format (i.e. `namespace` values, `key` and `value` contents) might change in the future;


`export-library-cdb` <br/>
`z-run` [ `<flag>` ... ] `export-library-cdb` `<cdb-path>`

  expects a single `<cdb-path>`, but no other arguments or scriptlet;
  writes to the specified file path the CDB database that represents the key-value store that backs the library;


`export-library-rpc` <br/>
`z-run` [ `<flag>` ... ] `export-library-cdb` `<rpc-target>`

  expects a single `<rpc-target>`, but no other arguments or scriptlet;
  listens to specified target for `z-run` specific RPC, that exports the library to remote clients;


`export-library-url` <br/>
`z-run` [ `<flag>` ... ] `export-library-url`

  expects no arguments or scriptlet;
  writes to `/dev/stdout` a line suitable for using it as value for the `--library-url` flag;
  currently it is either a CDB database file `<path>` or `<rpc-target>`;
  **however it should always be treated as an opaque value**, containing any ASCII character, except control characters, as it might change in future versions;


`export-library-fingerprint` <br/>
`z-run` [ `<flag>` ... ] `export-library-fingerprint`

  expects no arguments or scriptlet;
  writes to `/dev/stdout` a line containing the fingerprint of the current library version;
  currently it is an hex-encoded hash;
  **however it should always be treated as an opaque value**, containing any ASCII character, except control characters, as it might change in future versions;


`parse-library` <br/>
`z-run` [ `<flag>` ... ] `parse-library`

  expects no arguments or scriptlet;
  writes to `/dev/stdout` a single JSON object that represents the internal serialization of the library object;
  it uses pretty-printing, thus the JSON object will span over multiple lines;
  the serialization format is likely to change in the future;


### Advanced modes


`z-run` `--exec` `<source-path>` [ `<scriptlet>` [ `<argument>` ... ] ]

  expects a library `<source-path>`, an optional `<scriptlet>` and scriptlet `<arguments>`;
  it behaves similarly with the `execute-scriptlet` command;
  it enables one to write executable `z-run` scripts by using the `#!/usr/bin/env -S z-run --exec` header;


`z-run` `--ssh` [ `<ssh-flag>` | `<flag>` ... ] `<scriptlet>` [ `<argument>` ... ]

  expects a mandatory `<scriptlet>`, plus optional scriptlet `<arguments>`;
  similar to the `execute-scriptlet-ssh` command, however it allows certain SSH specific arguments as discussed bellow;


`z-run` `--invoke` `<invoke-payload>`

  expects a mandatory `<invoke-payload>` argument, no other flags, arguments or scriptlet;
  the payload contains encoded all the necessary flags, scriptlet and scriptlet arguments;
  it behaves `execute-scriptlet` command, however it allows one to easily execute `z-run` over SSH without bothering with `ssh` and `sh` command quoting and escaping;
  the serialization format is likely to change in the future;


### Input modes


`z-run` `--select`

  expects no arguments;
  reads from `/dev/stdin` a list of strings (mandatory compliant with UTF-8), presents a menu to the user, and if anything is selected it writes it to `/dev/stdout` followed by `'\n'`;
  it expects (and checks) that both `/dev/stdin` and `/dev/stdout` are non-TTY;  (i.e. they must be redirected to a file, pipe, or socket;)
  it expects (and checks) that `/dev/stderr` is a TTY, and thus requires an usual `TERM` value;
  it uses an embedded variant of the `fzf(1)` tool, disregardin any `fzf` specific flags or environment variables;  (but this should be treated as an implementation detail, and not relied upon;)

`z-run` `--input` [ `--message=<message>` ] [ `--prompt=<prompt>` ] [ `--sensitive` ]

  optionally allows any of the flags above;
  writes to `/dev/stderr` the `<message>` followed by `'\n'`;
  writes to `/dev/stderr` the `<prompt>`, or by default `'>> '`;
  if `--sensitive` is specified, it disables input echo;
  reads from `/dev/stderr` a single line (up to the first `'\n'`), that it then writes to `/dev/stdout`;
  all values (message, prompt, and input) must be compliant with UTF-8;
  it expects (and checks) that `/dev/stdout` is non-TTY;  (i.e. it must be redirected to a file, pipe or socket;)
  it expects (and checks) that `/dev/stderr` is a TTY, and thus requires an usual `TERM` value;


`z-run` `--fzf` [ ... ]

  optionally allows any of the flags accepted by `fzf(1)`;
  similar to `--select`, however it allows customizing `fzf(1)` through `fzf` specific flags and environment variables;
  it expects (and checks) that both `/dev/stdin` and `/dev/stdout` are non-TTY;  (i.e. they must be redirected to a file, pipe, or socket;)
  it expects (and checks) that `/dev/stderr` is a TTY, and thus requires an usual `TERM` value;




### Miscellaneous


`z-run` `--version`

  writes to `/dev/stdout` a series of lines describing the version, executable, build related, and other miscellaneous information;
  the output format is likely to change in the future;

`z-run` `--version`

  writes to `/dev/stdout` a copy of this manual;




## FLAGS


`--untainted`

  if `z-run` is invoked within the context of `z-run` execution, disregard the context, and treat this invocation as a new separate context;
  **must appear as the first flag;**
  (the same applies in `--exec` mode, which implies `--untainted`;)


`--exec`, `--ssh`, and `--invoke`

  these trigger the advanced execution modes described in the sections above;
  **must appear as the first flag;**


`--select`, `--input`, and `--fzf`

  these trigger the input execution modes described in the sections above;
  **must appear as the first flag;**


`--library-source=<source-path>`

  specifies a library `<source-path>`, that overrides the default library source detection mecanism;


`--library-url=<cache-url>`

  specifies a library `<cache-url>`, either a `<cdb-path>` or a `<rpc-target>`;
  specifying both `--library-source=...` and `--library-url=...` is not allowed;


`--workspace=<path>`

  specifies a folder that `z-run` switches to before executing;
  (if no `--library-source=...` or `--library-url=...` is specified, the default library source detection mecanism uses this folder as the root;)


`--ssh-target=<ssh-target>`

  **only in SSH mode;**
  specifies the SSH target;


`--ssh-workspace=<path>`

  **only in SSH mode;**
  specifies a path on the remote machine that `z-run` switches to before executing;


`--ssh-export=<name>`

  **only in SSH mode;**
  specifies an environment variable name that is exported on the remote machine;


`--ssh-path=<path>`

  **only in SSH mode;**
  specifies a value that is appended to the `PATH` environment variable on the remote machine;


`--ssh-terminal=<terminal>`

  **only in SSH mode;**
  specifies a value that overrides the `TERM` environment variable on the remote machine;




## ENVIRONMENT


`ZRUN_LIBRARY_SOURCE`

  an alternative to the `--library-source=...` flag;
  never exported inside scriptlet execution;

`ZRUN_LIBRARY_URL`

  an alternative to the `--library-url=...` flag;
  always exported inside the scriptlet execution environment;  (never unset it explicitly;)

`ZRUN_WORKSPACE`

  an alternative to the `--workspace=...` flag;
  always exported inside the scriptlet execution environment;  (never unset it explicitly;)

`ZRUN_EXECUTABLE`

  always exported inside the scriptlet execution environment;  (never unset it explicitly;)

`ZRUN_LIBRARY_FINGERPRINT`

  always exported inside the scriptlet execution environment;  (never unset it explicitly;)

`ZRUN_CACHE`

  an alternative folder to the default `$HOME/.cache/z-run`, where various files (and pipes, sockets, etc.) are created;
  (if explicitly specified, it is exported in the scriptlet execution environment;)




## LIBRARY




### Library source resolution

**TBD**




### Library source syntax

**TBD**




### Library directives

**TBD**




### Scriptlet bodies

**TBD**



### Scriptlet interpreters

**TBD**




## NOTES




### Why `z-run` vs classical scripts?

What sets `z-run` aside from a "classical" scripts is:

* first it allows one to easily mix multiple languages;
for example one can easily write a Bash script that chains together a Python, Ruby and another Bash script,
all within the same file;
(with classical scripts one would have had to create multiple files, one for each language, and explicitly invoke them via the interpreter;)

* it allows one to easily bundle together unrelated scripts;
for example say one uses `z-run` to aid in the development process, then one could bundle together scripts related
to compiling, testing, deploying, and other single-shot shortcuts in the same file;
(again, with classical scripts, one would have usually put each of these in separate files,
or by using `case` or `function` try to manage them inside a reduced set of files;)

* it allows one to easily invoke scriptlets from within the same library, without caring what language they are written in, or where their source code is stored;
(with classical scripts one has to take care of getting the paths to other scripts right, especially when running from a different folder than the original one;)

However `z-run` takes it a step further by offering the following functionalities:

* environment variables management -- by allowing certain environment variables to be overriden, removed or appended (like in the case of `$PATH`);
(with classical scripts one has to manage these by themselves, resort to various "envdir" tools, or interact with shell magic;)

* hierarchical menus -- by simply calling `z-run` without any arguments, it presents one with a selectable menu of all the scriptlets in the library,
and after one selects a scriptlet it executes it;  moreover one can group various scriptlets under other sub-menus, thus allowing arbitrary nesting and navigation;
(when used from a terminal, `z-run` uses and embedded `fzf`-based UI;  when used from outside a terminal, `z-run` tries to use `z-run--select` (any system), `rofi` or `dmenu` (under Linux and BSD's) or `choose` (under OSX);)

* SSH-based remote execution -- one can easily execute a given scriptlet on a remote server without having to previously copy anything there;
moreover once the scriptlet is executing on the remote server, it can invoke other scriptlets from the library that are also to be executed remotely;
(the only requirement is having `z-run` installed on the remote machine;)

* scriptlets generation -- one can easily write a scriptlet that generates `z-run` compliant source, thus generating other scriptlets based on arbitrary criteria;
(for example one could write a scriptlet, that generates other scriptlets specific for each file in a given folder;)

* library compiling -- one can easily create a single file that contains the entire library, which can then be moved and used in another place;

* Go-based templates -- that are useful especially in generating other scriptlets;




### Caching

`z-run` is designed to support thousands (and in extreme cases tens of thousands) of scriptlets, especially when scriptlet generation is used.

Therefore `z-run` has built-in optimizations to cache the library contents (without any regeneration) unless any of the following conditions are met:

* any of the files that comprise the source code of the library are changed;  (based both on file timestamps and contents hashing;)
* any of the environment variables change;
* a different version of `z-run` is installed;  (this covers both upgrading and downgrading;)

However once a scriptlet is executed, it and any other invoked scriptlets will use exactly the same cached library contents.
Thus it is safe to change the source code of the library, while a scriptlet is executing.




## EXAMPLES

**TBD**

