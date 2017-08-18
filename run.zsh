#!/usr/bin/env zsh

# see https://unix.stackexchange.com/questions/32407/zsh-excluding-files-from-a-pattern
setopt extendedglob
go run (^*_test).go $@


