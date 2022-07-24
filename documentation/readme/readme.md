



![logo](./documentation/logo.png)




----




# **z-run** -- lightweight Go-based scripts library tool


> Table of contents:
> * [Manual](#manual) and [Examples](#examples)
> * [Installation](#installation) and [FAQ](#faq)
> * [About](#about), [Status](#status), [Copyright and licensing](#notice-copyright-and-licensing) and [SBOM](#sbom-software-bill-of-materials)
> * [chat on Discord](https://discord.gg/WjBfs8rz), [discuss on GitHub](https://github.com/volution/z-run/discussions/categories/discussions), or [email author](mailto:ciprian.craciun@gmail.com)




----




## About


The best way to describe `z-run`, is to look at a simple example
(you know, "a snippet is worth a thousand man-pages"):
~~~~
<< hello world (in shell)
	echo "hello world!"
!!
<< hello world (in Python)
	#! <python3>
	print("hello world!")
!!
<< hello world (in PHP)
	#! <php>
	<?php print("hello world!") ?>
!!

:: ping / google :: z-run ':: ping / *' 8.8.8.8
:: ping / cloudflare :: z-run ':: ping / *' 1.1.1.1
--<< ping / *
	ping "${@}"
!!

<< ip / addr / show
	ip --json addr show \
	| z-run ':: ip / addr / show / jq'
!!
--<< ip / addr / show / jq
	#! <jq> -r
	.[]
	| select (.addr_info != [])
	| .ifname as $ifname
	| .addr_info[]
	| [$ifname, .family, .local]
	| @tsv
!!
~~~~


Provided that
one has saved the above into `_z-run` (in the current folder),
and has copied the `z-run` executable somewhere in the `$PATH`,
then one can just run `z-run` (in the current folder),
and be presented with the following menu,
from where one can select a scriptlet to be executed:
~~~~
  ping / google
  ping / cloudflare
  ip / addr / show
  hello world (in shell)
  hello world (in Python)
> hello world (in PHP)
  7/7
: _
~~~~


In a few sentences,
`z-run` is a lightweight and portable tool
that allows one to create and execute a library of scripts,
by mix-and-matching multiple languages,
from `bash`, Python, Ruby, PHP, up-to `jq`, `make` and `ninja`.


It works based on the observation that
the majority of interpreters can take as argument
a path that points to the file to be executed.
Based on this, when `z-run` is called to execute a scriptlet,
it creates a temporary file (and if small enough just a pipe),
and then calls the delegated interpreter.


Thus, `z-run` fulfills the following basic goals:
* allows one to easily author small scriptlets,
  independent of the interpreter,
  all in one file (or a small set of files);
* presents the user with a `fzf`-based menu,
  that allows one to easily select a scriptlet for execution;
* bootstraps the required script and environment for the actual interpreter;
* delegates the scriptlet execution to the actual interpreter;

Besides this, `z-run` also fulfills the following advanced goals:
* allows generating the library at compile-time;
  (sort of like macro expansion, but at a much lower level;)
* provides a built-in templating language;
  (especially useful in the previous feature;)
* provides remote execution of scriptlets over SSH,
  without any additional requirements than
  having `z-run` deployed on the remote machine;




----




## Status


Currently, `z-run` is still in a pre-release state,
perhaps somewhere between a beta and a release-candidate.

There is no promise of backward or forward compatibility,
there is little documentation (besides the examples),
there is no testing harness,
and there is no roadmap.

However, **I do use it personally for all my scripting tasks**,
from project development to automation.
Moreover, **I do also use it in production for various operational tasks**,
from driving Ansible and LetsEncrypt, to remote execution.

Also see the "[How is it tested?](#how-is-it-tested)",
"[How quick are issues fixed?](#how-quick-are-issues-fixed)",
"[How to ask for help?](#how-to-ask-for-help)",
and "[How to get commercial support?](#how-to-get-commercial-support)"
questions in the FAQ section.


Here are some examples of where it is used:

* [z-run](https://github.com/volution/z-run/tree/development/scripts)
  (this tool itself)
  -- used for all project development tasks (from building and testing to publishing);

* [z-scratchpad](https://github.com/volution/z-scratchpad/tree/master/scripts)
  (wiki-like and notes tool)
  -- used for all project development tasks;

* [kawipiko](https://github.com/volution/kawipiko/tree/development/scripts)
  (static HTTP server using CDB)
  -- used for all project development tasks;

* [vonuvoli](https://github.com/volution/vonuvoli-scheme/tree/development/scripts)
  (R7RS Scheme interpreter)
  -- used for all project development tasks, plus building remote building;

* [md5-tools](https://github.com/volution/volution-md5-tools/tree/master/scripts)
  (MD5/SHA*/Blake*/etc. recursive parallel hasher)
  -- used for all project development tasks;

* [covid19-datasets](https://github.com/cipriancraciun/covid19-datasets/tree/master/scripts)
  (derived and augmented COVID19 datasets based on JHU and NY-Times)
  -- used for all workflow aspects, from downloading, cleaning, merging, and publishing;
  (it drives `ninja`, `jq`, Python and Julia;  running everything remotely on a dedicated server;)

* [hyper-simple-server](https://github.com/console9/hyper-simple-server/tree/master/scripts) and
  [hyper-static-server](https://github.com/console9/hyper-static-server/tree/master/scripts)
  (low-level HTTP servers written in Rust based on `hyper`)
  -- used for all project development tasks;


Besides these few public examples, all my development projects (personal or professional)
have their `z-run` scripts.
I have even built a simple map-reduce framework to handle offline log processing.




----




## Manual

At the moment, the documentation is quite scarce...

One can consult the following resources:
* the `z-run` draft help in [documentation/help/z-run.txt](./documentation/help/z-run.txt);
* the `z-run` draft manual in [documentation/manual/z-run.1.ronn](./documentation/manual/z-run.1.ronn);
* `z-run --help`, `z-run --manual`, and `z-run --readme`;
  (the above files and this readme are embedded in the executable itself;)
  (there is also `--manual-man`, `--manual-html`, and `--readme-html` for `man` / HTML formats;)
* the various example files in [examples](./examples), which expose most of the basic and advanced `z-run` features;
* the scripts folders linked-at in the [status](#status) section;
* the simple snippets in the [examples](#examples) section;




----




## Examples


> **For the moment, this section shall also serve as a manual.**




### Introduction


Prerequisites:
* having `z-run` deployed somewhere on the `$PATH`;
  (it is recommended to place it in `/usr/local/bin`;)
* nothing else!


Terminology:

* scriptlet
  -- a small and simple script,
  executed by an interpreter (`bash`, Python, `jq`, etc.),
  focusing on one single well-defined task,
  delegating more lower-level tasks to other scriptlets;
  (you know, follow the old UNIX philosophy of KISS, "[keep-it-simple-stupid](https://en.wikipedia.org/wiki/KISS_principle)",
  or its more civilized variant "[do one thing and do it well](https://en.wikipedia.org/wiki/Unix_philosophy#Do_One_Thing_and_Do_It_Well)",
  and use a "[top-down structured programming](https://en.wikipedia.org/wiki/Top-down_and_bottom-up_design#Programming)" approach;)

* scriptlet label
  -- the unique identifier of the scriptlet,
  usually in a hierarchical form like `category / subcategory / task`;
  this is used when invoking a scriptlet as in `z-run ':: category / subcategory / task'`;

* scriptlet body
  -- the source code of the scriptlet,
  passed by `z-run` to the delegated interpreter;
  `z-run` doesn't look, change, or care what is inside the body;
  it is taken verbatim, without any expansions or changes;
  (it only must be a valid UTF-8 text;)

* library
  -- a set of related scriptlets,
  managed by `z-run`;

* workbench
  -- `z-run` allows one to have multiple self-contained scriptlet libraries,
  therefore, the folder where such library resides is called a workbench;
  however a workbench can contain other files,
  perhaps files the scriptlets require or operate on;
  for example in case of development projects,
  the library can be placed inside the `scripts` folder,
  but the workbench is the entire project folder;

* library source file
  -- the library is compiled from a set of source files;
  each can contain one or multiple scriptlets;
  each can include one or multiple other source files;


Behavior:

* `z-run` should be executed in the workbench folder;
  (thus the current folder is usually the workbench folder;)

* `z-run` searches the workbench folder for scriptlets based on the following rules:
  * first it tries to see if there is a `z-run`, `zrun`, `scriptlets`, `scripts`, or `bin` folder;
    (perhaps prefixed by `_` or `.`, as in `z-run`, `_z-run`, or `.z-run`;)
  * if none such folder exists, it assumes that the workbench folder is the basis for the library;
  * then given the folder previously identified, it tries to see if there is a `z-run` or `zrun` file;
    (perhaps prefixed by `_` or `.`, as in `z-run`, `_z-run`, or `.z-run`;)
  * if none such file exists, it then tries from the first step, but looking into `.git`, `.hg`, `.svn`, `.bzr`, `.darcs`;
    (this allows one to hide `z-run` scripts inside the VCS folder;)

* `z-run` only accepts valid UTF-8 files for the library source;

* `z-run` is quite "white-space" sensitive;
  (inside the scriptlet body it doesn't care;)

* `z-run` always handles source files relative to the workbench folder;

* `z-run` always executes scriptlets in the workbench folder;


Suggestions:

* although scriptlets can take arguments,
  if based on an argument the behavior changes radically,
  consider splitting that into two different scriptlets;

* have each scriptlet focus on a single well-defined task,
  consider delegating lower-level tasks to other scriptlets;

* instead of having complex shell pipe-lines,
  consider extracting part of a pipe-line into a dedicated scriptlet;

* if the same snippet of code repeats multiple times,
  consider extracting it into a dedicated scriptlet;

* if some scriptlets are not meant to be executed directly,
  consider hiding them from the main menu,
  by either catching them in sub-menus,
  or by completely hiding them by prefixing them with `--`,
  as in `--:: some-category / some-very-low-level-task :: ...`;

* many times the main purpose of a scriptlet is
  to just prepare the environment for another tool,
  and as the last action just executing that tool;
  in this case consider using `exec some-tool ...`,
  thus replacing the interpreter process with that tool's process;
  ( else one would end-up with a tree of processes that just wait for one-another;)




----




### Creating the library of scriptlets


`z-run` uses a simple syntax, one perhaps similar to `make`:

* each library source file contains a sequence of scriptlets (label and body),
  or a set of directives (including other source files, changing the environment, etc.);

* scriptlets can be written in one line
  by using the following syntax,
  where the `scriptlet label` can't contain `::`,
  and the interpreter is assumed to be `bash`:
~~~~
:: scriptlet label :: scriptlet body
~~~~

* scriptlets can be written on multiple lines
  by using the following syntax,
  where the `scriptlet label` can't contain `::`,
  the `scriptlet body` must be indented
  (tabs or spaces doesn't matter, just be consistent)
  (the first indentation level will be removed when sent to the interpreter),
  and the interpreter is assumed to be `bash`:
~~~~
<< scriptlet label
	scriptlet body line 1
	scriptlet body line 2
	...
!!
~~~~

* scriptlets can have other interpreters than `bash`
  by using the following syntax,
  just like the `#!` file header in normal scripts
  (mind the spaces between `#!` and the rest of that line):
~~~~
<< something in bash (A)
	...
!!
<< something in bash (B)
	#! <bash>
	...
!!
<< something in Python (A)
	#! <python2>
	...
!!
<< something in Python (B)
	#! <python3>
	...
!!
<< something in jq (with extra arguments)
	#! <jq> --slurp --raw-output
	...
!!
<< something in custom interpreter (A)
	#! my-custom-interpreter
	...
!!
<< something in custom interpreter (B)
	#! /usr/local/bin/my-custom-interpreter
!!
~~~~

* scriptlets can be excluded from the library
  (mind the two `##` just before `::` or `<<`):
~~~~
##:: this scriptlet is ignored (A) :: ...

##<< this scriptlet is ignored (B)
	...
!!
~~~~

* arbitrary multi-line comments can be written
  (mind the `##{{` and `##}}` that must start on the first column):
~~~~
##{{

some text that doesn't contain ##}}

##}}
~~~~

* scriptlets can be hidden from the main menu
  (mind the two `--` just before `::` or `<<`):
~~~~
--:: this scriptlet is hidden (A) :: ...

--<< this scriptlet is hidden (B)
	...
!!
~~~~

* scriptlets can be gathered in sub-menus
  (mind the `//` just after the `::`);
  the `...` suffix gathers all scriptlets having that prefix (except `...`)
  under that sub-menu, and hides them from the main menu;
  the `*` suffix also gathers all scriptlets having that prefix (except `*`),
  under that sub-menu, but doesn't hide them from the main menu;
  the `::// *` entry is just a sub-menu that contains all other scriptlets:
~~~~
::// category-a / ...
::// category-b / *
::// *

:: category-a / scriptlet-1 :: ...
:: category-a / scriptlet-2 :: ...
:: category-b / scriptlet-1 :: ...
:: category-b / scriptlet-2 :: ...
~~~~

* scriptlets can be forced in the main menu
  (mind the `++` just before `::` or `<<`):
~~~~
++:: this scriptlet is forced (A) :: ...

++<< this scriptlet is forced (B)
!!
~~~~

* environment variables can be defined to string values:
~~~~
&&== env SOME_ENV some-value
~~~~
* environment variables can be defined to absolute paths:
~~~~
&&== env SOME_PATH ./relative-path
~~~~
* the `$PATH` environment variable can be appended with absolute paths:
~~~~
&&== path ./bin
~~~~
* environment variables, similar to `$PATH`, can be appended with absolute paths:
~~~~
&&== env-path-append SOME_ENV_LIKE_PATH ./bin
~~~~
* environment variables can be removed:
~~~~
&&== env-exclude SOME_ENV_1 SOME_ENV_2
~~~~
* environment variables can be provided with defaults:
~~~~
&&== env-fallback SOME_ENV default-value
~~~~

* other library source files can be included:
~~~~
&& _/file-relative-to-the-current-source-file-parent
&& ./file-relative-to-the-current-workbench-folder
&& /file-with-absolute-path

&&?? _/file-that-might-not-exist
&&?? ./file-that-might-not-exist
~~~~


**TBD:**
* scriptlets bodies from files;
* scriptlets replacement bodies;
* single-line scriptlets with custom interpreter;
* complex sub-menus;
* scriptlets generators;
* depending on files or folders without including them (useful for generators);
* using print scriptlets;
* using template scriptlets;
* using the `bash+` interpreter extension;
* using the `python3+` interpreter extension;
* using the `go` interpreter;
* using the `go+` interpreter;
* using the `starlark` interpreter;
* environment variables values substitutions;
* included source file paths substitutions;



----




### Using the library of scriptlets from the console


* showing the main menu:
~~~~
z-run
~~~~

* showing a sub-menu
  (the label must exactly match how it was defined,
  that is with `...` or `*` suffix):
~~~~
z-run ':: category-a / ...'
z-run ':: category-b / *'
z-run ':: *'
~~~~

* executing a scriptlet
  (the `::` is mandatory, regardless of how the scriptlet was defined,
  with `::` or `<<`):
~~~~
z-run ':: category-a / scriptlet-1'
z-run ':: category-b / scriptlet-2' some-argument ...
~~~~

* listing all the scriptlet labels:
~~~~
z-run list
~~~~

* select a scriptlet from the menu,
  and print its label:
~~~~
z-run select-label
~~~~

* select a scriptlet from the menu,
  and print its body:
~~~~
z-run select-body
~~~~




**TBD:**
* compiling libraries and using them;
* executing remote scriptlets over SSH;
* using a library as a standalone tool with `z-run --exec`;
* executing standalone scriptlets with `z-run --scriptlet`;
* executing standalone scriptlets with `z-run --scriptlet-exec`;
* executing standalone templates with `z-run --template`;
* `z-run --input`, `z-run --select`, and `z-run --fzf`;
* `z-run --shell`;
* `z-run --export=shell-functions`;
* `z-run --version`;
* `z-run --help`;
* `z-run --manual`;
* `z-run --manual-man`;
* `z-run --manual-html`;
* `z-run --readme`;
* `z-run --readme-html`;
* `z-run --sbom`;
* `z-run --sbom-html`;
* `z-run --sbom-json`;
* `z-run --sources-md5`;
* `z-run --sources-cpio | gunzip | cpio -i -t`;




----




## Installation




### Download prebuilt executables


One should see the [releases page on GitHub](https://github.com/volution/z-run/releases),
where pre-built executables (only Intel 64bit architectures) are available for:
* Linux (the main development and testing environment) -- <br/>
  <https://github.com/volution/z-run/releases/download/v0.18.1/z-run--linux--v0.18.1>
* OSX (only Intel CPU's) (the second targeted environment) -- <br/>
  <https://github.com/volution/z-run/releases/download/v0.18.1/z-run--darwin--v0.18.1>
* OpenBSD (seldom tested) -- <br/>
  <https://github.com/volution/z-run/releases/download/v0.18.1/z-run--openbsd--v0.18.1>
* FreeBSD (seldom tested) -- <br/>
  <https://github.com/volution/z-run/releases/download/v0.18.1/z-run--freebsd--v0.18.1>

Also, each of these files are signed with my PGP key `5A974037A6FD8839`, thus do check the signature.


For example:

* import my PGP key:
~~~~
curl -s https://github.com/cipriancraciun.gpg | gpg2 --import
~~~~
~~~~
gpg: key 5A974037A6FD8839: public key "Ciprian Dorin Craciun <ciprian@volution.ro>" imported
gpg: Total number processed: 1
gpg:               imported: 1
~~~~

* download the executable and signature (replace the `linux` token with `darwin` (for OSX), `freebsd` or `openbsd`):
~~~~
curl \
        -s -S -f -L \
        -o /tmp/z-run \
        https://github.com/volution/z-run/releases/download/v0.18.1/z-run--linux--v0.18.1 \
#

curl \
        -s -S -f -L \
        -o /tmp/z-run.asc \
        https://github.com/volution/z-run/releases/download/v0.18.1/z-run--linux--v0.18.1.asc \
#
~~~~

* verify the executable:
~~~~
gpg2 --verify /tmp/z-run.asc /tmp/z-run
~~~~

* **check that the key is `58FC2194FCC2478399CB220C5A974037A6FD8839`**:
~~~~
gpg: assuming signed data in '/tmp/z-run'
gpg: Signature made Sun Oct  3 21:49:50 2021 EEST
gpg:                using DSA key 58FC2194FCC2478399CB220C5A974037A6FD8839
gpg: Good signature from "Ciprian Dorin Craciun <ciprian@volution.ro>" [unknown]
gpg:                 aka "Ciprian Dorin Craciun <ciprian.craciun@gmail.com>" [unknown]
gpg: WARNING: This key is not certified with a trusted signature!
gpg:          There is no indication that the signature belongs to the owner.
Primary key fingerprint: 58FC 2194 FCC2 4783 99CB  220C 5A97 4037 A6FD 8839
~~~~

* change the executable permissions:
~~~~
chmod a=rx /tmp/z-run
~~~~

* copy the executable on the `$PATH`:
~~~~
sudo cp /tmp/z-run /usr/local/bin/z-run
~~~~

* check that it works:
~~~~
z-run --version
~~~~
~~~~
* tool          : z-run
* version       : 0.19.0
* executable    : /usr/local/bin/z-run
* build target  : release, linux-amd64, go1.17.6, gc
* build number  : 12170, 2022-01-19-23-30-43
* code & issues : https://github.com/volution/z-run
* sources git   : a6e6ce1dbd63eb871ad0280d6dbafdbc9e7960d9
* sources hash  : d7ba8be84e85329ac187984750ecd242
* uname node    : some-workstation
* uname system  : Linux, 5.15.6-1-default, x86_64
* uname hash    : 7c8643e588e6fbb4063d2643de27c39b
~~~~




### Build from sources


Go is a prerequisite;
one can install it from any Linux or BSD package manager,
or OSX's `brew`,
or just downloading it from <https://golang.org/dl>.

The first step is preparing the environment:
~~~~
mkdir \
        /tmp/z-run \
        /tmp/z-run/bin \
        /tmp/z-run/src \
        /tmp/z-run/go \
#
~~~~

The second step is cloning the Git repository:
~~~~
git clone \
        --branch development \
        --depth 1 \
        http://github.com/volution/z-run.git \
        /tmp/z-run/src \
#
~~~~

Alternatively, one can just fetch the sources bundle:
~~~~
curl \
        -s -S -f -L \
        -o /tmp/z-run/src.tar.gz \
        https://github.com/volution/z-run/archive/refs/heads/development.tar.gz \
#

tar \
        -x -z \
        -f /tmp/z-run/src.tar.gz \
        -C /tmp/z-run/src \
        --strip-components 1 \
#
~~~~

Finally, compiling the `z-run` executable:
~~~~
cd /tmp/z-run/src/sources

env \
        GOPATH=/tmp/z-run/go \
go build \
        -tags 'netgo' \
        -gcflags 'all=-l=4' \
        -ldflags 'all=-s' \
        -trimpath \
        -o /tmp/z-run/bin/z-run \
        ./cmd/z-run.go \
#
~~~~

Optionally, deploying the `z-run` executable:
~~~~
sudo \
cp \
        /tmp/z-run/bin/z-run \
        /usr/local/bin/z-run \
#
~~~~

Optionally, check that `z-run` works:
~~~~
z-run --version
~~~~
~~~~
* tool          : z-run
* version       : 0.19.0
* executable    : /usr/local/bin/z-run
* build target  : release, linux-amd64, go1.17.6, gc
* build number  : 12170, 2022-01-19-23-30-43
* code & issues : https://github.com/volution/z-run
* sources git   : a6e6ce1dbd63eb871ad0280d6dbafdbc9e7960d9
* sources hash  : d7ba8be84e85329ac187984750ecd242
* uname node    : some-workstation
* uname system  : Linux, 5.15.6-1-default, x86_64
* uname hash    : 7c8643e588e6fbb4063d2643de27c39b
~~~~




----




## FAQ




### How to ask for help?

If you have encountered a bug,
just use the [GitHub issues](https://github.com/volution/z-run/issues).

If you are not sure about something,
want to give feedback,
or request new features,
just use the [GitHub discussions](https://github.com/volution/z-run/discussions/categories/discussions).

If you want to ask a quick question,
or just have a quick chat,
just head over to the [Discord channel](https://discord.gg/WjBfs8rz).


### How is it tested?

I use it myself every day,
in almost every task that involves my laptop,
from projects development,
to production operations,
application launching,
and even for auto-completion in the editor.

Most of my use-cases use many of the advanced features,
thus most corner-cases are well covered.


### How quick are issues fixed?

As stated in the previous answer,
because I use `z-run` myself for everything,
any major issue means I can't do my work;
thus I need to stop everything
and just fix the issue.

As a consequence, the turn-around-time is quite small.


### How to get commercial support?

If you want to use `z-run` in production,
and want to be sure that you are using it correctly,
or just want to be sure that issues (or features) that you need get prioritized,
I own a European limited-liability company
so just [email me](mailto:ciprian.craciun@gmail.com),
and we can discuss the options.




### Why `z-run` instead of classical scripts?


What sets `z-run` aside from a classical scripts is this:

* first, it allows one to easily mix multiple languages;
  for example one can easily write a shell script that chains together a Python, Ruby and another shell script,
  all within the same file, without having to resort to here-documents or quoting;
  (with classical scripts, one would have had to create multiple files,
  one for each language,
  and explicitly invoke them via the interpreter;)

* it allows one to easily bundle together unrelated scripts;
  for example say one uses `z-run` to aid in the development process,
  then one could bundle together scripts related to
  compiling, testing, deploying, and other single-shot shortcuts in the same file;
  (with classical scripts, one would usually have to put each of these in separate files,
  or by using `case` or `function` try to manage them inside a reduced set of files;)

* it allows one to easily invoke scriptlets from within the same library,
  without caring what language they are written in,
  or where their source code is stored;
  (with classical scripts, one has to take care of getting the paths to other scripts right,
  especially when running from a different folder than the original one;)


However, `z-run` takes it a step further by offering the following functionalities:

* hierarchical menus
  -- by simply calling `z-run` without any arguments,
  it presents one with a selectable menu of all the scriptlets in the library,
  and after one selects a scriptlet, it executes it;
  moreover one can group various scriptlets under other sub-menus,
  thus allowing arbitrary nesting and navigation;
  (when used from a terminal, `z-run` uses and embedded `fzf`-based UI;
  when used from outside a terminal, `z-run` tries to use
  `z-run--select` (any system),
  `rofi` or `dmenu` (under Linux and BSD's),
  or `choose` (under OSX);)

* environment variables management
  -- by allowing certain environment variables to be
  overridden, removed or appended (like in the case of `$PATH`);
  (with classical scripts, one has to manage these by themselves,
  resort to various "envdir" tools,
  or interact with shell magic;)

* SSH-based remote execution
  -- one can easily execute a given scriptlet on a remote server
  without having to previously copy anything there;
  moreover, once the scriptlet is executing on the remote server,
  it can invoke other scriptlets from the library that are also to be executed remotely;
  (the only requirement is having `z-run` installed on the remote machine;)

* scriptlets generation
  -- one can easily write a scriptlet that generates `z-run` compliant source,
  thus generating other scriptlets based on arbitrary criteria;
  (for example one could write a scriptlet that generates other scriptlets specific for each file in a given folder;)

* library compiling
  -- one can easily create a single file that contains the entire library,
  which can then be moved and used in another place;

* Go-based templates
  -- that are useful especially in generating other scriptlets;




### What is the performance?


`z-run` is designed to support thousands (and in extreme cases tens of thousands) of scriptlets,
especially when scriptlet generation is used.
For example:
* my own COVID19 workbench has ~20K scriptlets;
* my own photography workbench has ~47K scriptlets;
* and a production operations workbench has ~4K scriptlets;

Therefore, `z-run` has built-in optimizations
to cache the library contents (without any regeneration),
unless any of the following conditions are met:
* any of the files that comprise the source code of the library are changed;
  (based both on file timestamps and contents hashing;)
* any of the environment variables change;
* a different version of `z-run` is installed;
  (this covers both upgrading and downgrading;)

However, once a scriptlet is executed,
it and any other invoked scriptlets
will use exactly the same cached library contents.
Thus, it is safe to change the source code of the library
while a scriptlet is executing
(unlike shell scripts,
that because the shell reads one line at a time in a read-eval-loop,
it thus trips if the script is changed beyond the currently executed line).




### Why Go?


Because Go is highly portable, highly stable,
and especially because it can easily support
cross-compiling statically linked executables
to any platform it supports.




### Why not Rust?


Because Rust fails to easily support
cross-compiling (statically or dynamically linked) executables
to any platform it supports.

Because Rust is less portable than Go;
for example, Rust doesn't consider OpenBSD as a "tier-1" platform.




### What other open-source code does it depend on?

* [fzf](https://github.com/junegunn/fzf)
  (via my custom [fork](https://github.com/cipriancraciun/fzf));
* [cdb](https://github.com/colinmarc/cdb)
  (via my custom [fork](github.com/cipriancraciun/go-cdb-lib));
* [go-flags](github.com/jessevdk/go-flags);
* [liner](github.com/peterh/liner);
* also see the [go.mod](./sources/go.mod) file
  for other minor dependencies;




### Some related readings...

* [Development Environments at Slack](https://slack.engineering/development-environments-at-slack/)




----




## Notice (copyright and licensing)


### Authors


Ciprian Dorin Craciun:
* <ciprian@volution.ro> or <ciprian.craciun@gmail.com>
* <https://volution.ro/ciprian>
* <https://github.com/volution>
* <https://github.com/cipriancraciun>


Please also see the [SBOM (Software Bill of Materials)](./documentation/sbom/sbom.md)
for links this project's dependencies and their authors.




### Notice -- short version


The code is licensed under GPL 3 or later.

This is the same license as used by
Linux, Git, MySQL/MariaDB, WordPress, F-Droid,
and many other projects.

If you **change** the code within this repository
and use it for **non-personal** purposes,
you'll have to release the changed source code as per GPL.




### Notice -- long version


For details about the copyright and licensing,
please consult the [notice.txt](./documentation/licensing/notice.txt) file
in the [documentation/licensing](./documentation/licensing) folder.

If someone requires the sources and/or documentation to be released
under a different license, please email the authors,
stating the licensing requirements, accompanied by the reasons
and other details; then, depending on the situation, the authors might
release the sources and/or documentation under a different license.




### SBOM (Software Bill of Materials)


This project, like many other open-source projects,
incorporates code from other open-source projects
(besides other tools used to develop, build and test).

Strictly related to the project's dependencies (direct and transitive),
please see the [SBOM (Software Bill of Materials)](./documentation/sbom/sbom.md)
for links to these dependencies and their licenses.

