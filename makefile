build:
	go build -o taskcli

package: build
	mkdir -p packaging/deb/usr/local/bin
	mkdir -p packaging/deb/etc/systemd/system
	cp taskcli packaging/deb/usr/local/bin/
	cp packaging/linux/taskcli.service packaging/deb/etc/systemd/system/
	dpkg-deb --build packaging/deb taskcli_1.0.0_amd64.deb
