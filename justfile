# Copyright 2024 Shaolong Chen. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

alias t := test
alias c := check

default:
  just --list

check: fmt lint test
  license-eye header check

fmt:
  golines --max-len=99 --base-formatter="gofumpt -extra" -w .
  gosimports -local github.com/maolonglong -w .

lint:
  go vet ./...

test:
  go test -shuffle=on -v -race -count=1 ./...
