How I Set Up Xv6 On macOS
August 1, 2024

To learn more about operating systems I've set up [Xv6](https://en.wikipedia.org/wiki/Xv6), a simple Unix-like operating system designed for teaching. It's loosely based on the v6 version of Unix.

## Setup

I found some basic setup instructions at the bottom of the [README](https://github.com/mit-pdos/xv6-riscv):

> You will need a RISC-V "newlib" tool chain from
> https://github.com/riscv/riscv-gnu-toolchain, and qemu compiled for
> riscv64-softmmu. Once they are installed, and in your shell
> search path, you can run "make qemu".

Looking at [the toolchain installation](https://github.com/riscv-collab/riscv-gnu-toolchain?tab=readme-ov-file#troubleshooting-build-problems), I noticed the following warning about case-insensitive file systems:

> If building a linux toolchain on a MacOS system, or on a Windows system using the Linux subsystem or cygwin, you must ensure that the filesystem is case-sensitive. A build on a case-insensitive filesystem will fail when building glibc because _.os and _.oS files will clobber each other during the build eventually resulting in confusing link errors.

I'm using macOS Sonoma for this, so I'm using a case-insensitive file system (APFS). To avoid a world of pain here I instead decided to go with the prebuilt packages at https://github.com/riscv-software-src/homebrew-riscv.

```
brew tap riscv-software-src/riscv
brew install riscv-tools
brew install riscv64-elf-gdb
```

These commands initially failed for me. I needed to update to the latest Xcode command line tools to be able to install the homebrew packages. I ran `xcode-select --install` and followed the prompts.

I already had a compatible version of QEMU, I imagine that'd just be `brew install qemu`.

Then I cloned the [Xv6 repo](https://github.com/mit-pdos/xv6-riscv) and ran `make qemu`, and it worked first time! You get a shell prompt and can run some base commands like `ls` and `echo`. To exit, I found that I can use `Ctrl-A X`.

## Debugger

I'm starting the debugger by running `make CPUS=1 qemu-gdb` in one terminal, and then in another terminal run `riscv64-elf-gdb`. By running `riscv64-elf-gdb` from the Xv6 directory it will automatically load the `.gdbinit` file which sets up the debugger to connect to the QEMU instance.

## Book

There's accompanying book for which the source can be found at https://github.com/mit-pdos/xv6-riscv-book. I've built a recent PDF of the book using the following commands:

```
brew install --cask mactex
make
```

This outputs a `book.pdf` file in the repository directory.

## Course

Xv6 is used in an MIT course on operating systems. The currently latest version of this course is available at https://pdos.csail.mit.edu/6.828/2023/index.html. I also found a [playlist with lecture recordings from 2020](https://www.youtube.com/playlist?list=PLTsf9UeqkReZHXWY9yJvTwLJWYYPcKEqK).
