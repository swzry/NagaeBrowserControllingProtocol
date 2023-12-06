#!/usr/bin/env just --justfile

jfdir := replace(justfile_directory(), "\\", "/")

goos := `go env GOOS`

pkgtest_src := jfdir / "pkgtest"
pkgtest_dist := jfdir / "pkgtest" / "dist"
executable_suffix := if goos == "windows" { ".exe" } else { "" }

export GOOS := goos

default:
  just --list

tidy name:
	cd {{jfdir / name}} && go mod tidy -e

update name:
	cd {{jfdir / name}} && go get -u

pt_tidy name:
	cd {{pkgtest_src / name}} && go mod tidy -e

pt_update name:
	cd {{pkgtest_src / name}} && go get -u

pt_build name:
	cd {{pkgtest_src / name}} && go build -o {{pkgtest_dist / name + executable_suffix}}

pt_run name:
	cd {{ pkgtest_dist }} && {{pkgtest_dist / name + executable_suffix}}

cp_front name:
	cp -r {{ pkgtest_src / name / "dist" / "spa" / "*" }} {{ pkgtest_dist / "wwwroot" }}

