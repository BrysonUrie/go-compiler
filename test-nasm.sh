#!/bin/bash

nasm -f elf64 output.asm -o output.o
gcc output.o -o program    # Links against libc
./program
