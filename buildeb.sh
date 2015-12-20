#!/bin/sh
#colors
export red='\033[0;31m'
export green='\033[0;32m'
export NC='\033[0m' # No Color

set -x 
die () {
    echo -e >&2 "$@"
    echo -e "${NC}..."
    exit 1
}

if [[ $# -eq 0 ]];
 then
 die "${red}build number argument required as \$1, none provided"
fi

export BUILD=$1
export REFRESH=$2
export DEPLOY=$3
if [ "$2" != "" ]; then
##refresh dependencies#
echo -e "${green}updating go get libs $BUILD${NC}"
go get github.com/op/go-logging
go get github.com/vaughan0/go-ini
go get golang.org/x/net/icmp
go get golang.org/x/net/ipv4
go get golang.org/x/net/ipv6
go get github.com/andelf/go-curl
go get github.com/asaskevich/govalidator
go get github.com/cactus/go-statsd-client/statsd
go get github.com/gorilla/pat
go get github.com/gorilla/context
go get github.com/miekg/dns
go get github.com/gorhill/cronexpr
go get github.com/garyburd/redigo/redis
fi

#do libs first
cd src\/libs
#go test

test_status=$?

if test $test_status -eq 0
then
	echo -e "${green}Build is a go! libs ver $BUILD${NC}"
else
	echo -e "${red}Build is a NO go! dhcs ver $BUILD${NC}"
	exit
fi

cd -

#now onto executables

for src in nodder
do
	#what .deb we want to build
	export BASE=$src

	#test
	cd src\/$BASE
#	go test 

	test_status=$?

	if test $test_status -eq 0
	then
		echo -e "${green}Build is a go! $BASE ver $BUILD${NC}"
	else
		echo -e "${red}Build is a NO go! $BASE ver $BUILD${NC}"
		exit
	fi

	cd -
	
	#now we build
	go build -o bin\/$BASE -ldflags "-X main.Build $BUILD" src\/$BASE\/$BASE.go 

	status=$?

	if test $status -eq 0
	then
		echo -e "${green}Build is a go! $BASE ver $BUILD${NC}"
	else
		echo -e "${red}Build is a NO go! $BASE ver $BUILD${NC}"
		exit
	fi

exit

	#now we build .deb package
	rm -rf builds/$BUILD\/$BASE
	mkdir -p builds/$BUILD\/$BASE
	mkdir -p builds/$BUILD\/$BASE/var/dhcs
	mkdir -p builds/$BUILD\/$BASE/etc/dhcs
	mkdir -p builds/$BUILD\/$BASE/etc/init
	mkdir -p builds/$BUILD\/$BASE/var/log/dhcs
	mkdir -p builds/$BUILD\/$BASE/etc/monit/conf.d
	mkdir -p builds/$BUILD\/$BASE/DEBIAN/

	echo "Package: $BASE" >> builds/$BUILD\/$BASE/DEBIAN/control
	echo "Architecture: amd64" >> builds/$BUILD\/$BASE/DEBIAN/control
	echo "Maintainer: Andrew Yasinsky" >> builds/$BUILD\/$BASE/DEBIAN/control
	echo "Depends: debconf (>= 0.5.00)" >> builds/$BUILD\/$BASE/DEBIAN/control
	echo "Priority: optional" >> builds/$BUILD\/$BASE/DEBIAN/control
	echo "Version: $BUILD" >> builds/$BUILD\/$BASE/DEBIAN/control
	echo "Description: $BASE" >> builds/$BUILD\/$BASE/DEBIAN/control
	 
	echo "/etc/nodder/$BASE.cfg" >> builds/$BUILD\/$BASE/DEBIAN/conffiles
	echo "/etc/init/$BASE.conf" >> builds/$BUILD\/$BASE/DEBIAN/conffiles
	echo "/etc/monit/conf.d/$BASE" >> builds/$BUILD\/$BASE/DEBIAN/conffiles

	echo "#!/bin/sh" >> builds/$BUILD\/$BASE/DEBIAN/preinst
	echo "set -e" >> builds/$BUILD\/$BASE/DEBIAN/preinst

	#copy files where they need to be
	cp bin\/$BASE  builds/$BUILD\/$BASE/var/dhcs\/$BASE
	cp etc/nodder/$BASE.cfg  builds/$BUILD\/$BASE/etc/nodder/$BASE.cfg
	cp etc/init\/$BASE.conf  builds/$BUILD\/$BASE/etc/init\/$BASE.conf
	cp etc/monit/conf.d\/$BASE  builds/$BUILD\/$BASE/etc/monit/conf.d\/$BASE

	chmod 0775 builds/$BUILD\/$BASE/DEBIAN/preinst
	
	dpkg-deb --build builds/$BUILD\/$BASE

	status=$?

	if test $status -eq 0
	then
		echo -e "${green}DEB Build is a go! $BASE ver $BUILD${NC}"
	else
		echo -e "${red}DEB Build is a NO go! $BASE ver $BUILD${NC}"
		exit
	fi
done

rm -rf latest	
ln -s builds/$BUILD latest

if [ "$3" = "deploy" ]; then
	echo -e "${red}WARNING: pushing to debian repo"
	echo -e "${NC}..."
	git add .
	git commit -m "build $BUILD"
	git push origin master

	#push to mirror
	for src in nodder
	do
		
		export BASEDEB=$src
		export FINALDEB=$src-$BUILD-$( date +"%Y%m%d%H%M%S" )_amd64.deb

		#scp latest/$BASEDEB.deb root@SERVER:/root/BUILDARCH/$FINALDEB
	done	

else
	echo -e "${green}Not pushing to git or deploy add deploy as \$2 flag"
	echo -e "${NC}..."
fi