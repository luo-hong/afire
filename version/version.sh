# !/bin/sh
version=`git log --date=iso --pretty=format:"%cd @%H" -1`
if [ $? -ne 0 ]; then
    version="unknown version"
fi

branch=`git symbolic-ref --short -q HEAD 2>/dev/null`
if [ $? -ne 0 ]; then
    branch="unknown branch"
fi

compile=`date +"%F %T %z"`" by "`go version`
if [ $? -ne 0 ]; then
    compile="unknown datetime"
fi

describe=`git describe --tags 2>/dev/null`
if [ $? -eq 0 ]; then
    version="${version} @${describe}"
fi

SCRIPTDIR=$(cd $(dirname "${BASH_SOURCE[0]}") >/dev/null && pwd)
echo -E "package version

import \"fmt\"

const (
	githash = \"${version}\"
    branch = \"${branch}\"
	buildat = \"${compile}\"
    host = \"${HOSTNAME}\"
)

func Show() string {
	return fmt.Sprintf(\"git: %v\\nbranch: %v\nbuild: %v\\nhost: %v\", githash, branch, buildat, host)
}

" > ${SCRIPTDIR}/version.go