#!/bin/sh

# Run Test Suite
go test -race -v ./...
 
# Run test coverage on each subdirectories and merge the coverage profile.
echo "mode: count" > profile.cov
 
# Standard go tooling behavior is to ignore dirs with leading underscors
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d); do
  if ls $dir/*.go &> /dev/null; then
    go test -covermode=count -coverprofile=$dir/profile.tmp $dir
    if [ -f $dir/profile.tmp ]; then
    	cat $dir/profile.tmp | tail -n +2 >> profile.cov
    	rm $dir/profile.tmp
    fi
  fi
done
 
# Create the coverage report for submission to Coveralls
go tool cover -func profile.cov
