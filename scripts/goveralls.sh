#!/bin/bash
echo "mode: count" > profile.cov
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d); do
	if ls $dir/*.go &> /dev/null; then
		go test -covermode=count -coverprofile=$dir/profile.tmp $dir
		if [ -f $dir/profile.tmp ]; then
			cat $dir/profile.tmp | tail -n +2 >> profile.cov
			rm $dir/profile.tmp
		fi
	fi
done
goveralls -coverprofile=profile.cov -service=travis-ci -repotoken $COVERALLS_TOKEN
