%global source_version 2.3.0
%global tag 1
%global provider        github
%global provider_tld    com
%global gopath          %{_datadir}/gocode

Name:           cn
%global project         %{name}
%global repo            %{name}
%global provider_prefix %{provider}.%{provider_tld}/%{project}/%{repo}
%global import_path     %{provider_prefix}
Version:        %{source_version}
Release:        %{tag}%{?dist}
Summary:        A client tool to bootstrap S3 gateways that leverages container technologies
License:        Apache-2.0
Group:          System/Filesystems
URL:            https://%{import_path}
Source0:        https://%{import_path}/archive/v%{source_version}.tar.gz
Source1:        %{name}-vendor-%{source_version}.tar.xz
Source2:        rebuild-vendor.sh

%if !%{defined gobuild}
%define gobuild(o:) go build -compiler gc -ldflags "${LDFLAGS:-} -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \\n')" -a -v -x %{?**};
%endif

BuildRequires:  go-srpm-macros
BuildRequires:  golang
BuildRequires:  dep

%description
A client tool to bootstrap S3 gateways that leverages container technologies

%prep
%setup -q -a 1 -n cn-%version

# move content of vendor under Godeps
mkdir -p Godeps/_workspace/src
mv vendor/* Godeps/_workspace/src/

%build
export GOPATH=$(pwd):$(pwd)/Godeps/_workspace:%{gopath}
export LDFLAGS="$LDFLAGS -X main.version=%{source_version}"
%gobuild -o bin/cn main.go

%install
install -D -p -m 755 bin/cn %{buildroot}%{_bindir}/cn
install -D -p -m 644 cn.toml %{buildroot}%{_sysconfdir}/cn/cn.toml
install -D -p -m 644 contrib/cn_completion.sh %{buildroot}%{_sysconfdir}/bash_completion.d/cn_completion.sh

%files
%doc README.md
%{_bindir}/cn
%{_sysconfdir}/cn/cn.toml
%{_sysconfdir}/bash_completion.d/cn_completion.sh

%changelog
* Mon Mar 11 2019  Erwan Velu <evelu@redhat.com> - 2.3.0-1
- update-check: Improving 'update-check' output
- adds update notification for newer version of nano
- README.md - Add TOC and "Enable mgr dashboard" section
- start: bindmounts /run/udev and /run/lvm upon using blockdevice
- Readme: Bump the new release tag: v2.2.0
- Packaging: Update specfile version to v2.2.0
* Tue Jan 22 2019  Erwan Velu <evelu@redhat.com> - 2.2.0-1
- Adds :z option to support SeLinux
- utils: simplify the output of cn cluster status
- utils: do not output ceph health in cn cluster status
- Readme: Bump the new release tag: v2.1.1
- Packaging: Update specfile version to v2.1.1
* Wed Dec 19 2018  Erwan Velu <evelu@redhat.com> - 2.1.1-1
- contrib: %install failed at installing files
- Readme: Bump the new release tag: v2.1.0
- Packaging: Update specfile version to v2.1.0
* Tue Dec 18 2018  Erwan Velu <evelu@redhat.com> - 2.1.0-1
- contrib: Updating cn_completion.sh
- travis: The bash_completion update cannot be done in travis
- Makefile: go test must be run in the local context
- flavors: Improving 'flavors show' output
- config: Merge flavors even if no configuration file found
- main: Print configuration filename to stderr
- contrib: Improving bash_completion generation
- start: Removing wrong statement
- flavors: Improving help message
- image: Split cliShowAliases & listAliases
- config: Update getStringFromConfig() to mimic other get*FromConfig()
- quality: Improving english & style
- list: Reporting cluster flavor
- utils: Using switch() statement in inspectImage()
- config: Adding documentation
- Packaging: Adding cn.toml configuration file
- config: Don't print a message if no configuration file found
- config: Fixing typo
- config: Adding more parallelism in tests
- config: Populating 'flavors show default' command
- config: Reworking merging of default & other flavors
- config: Using isParameterExist() instead of custom IsSet() calls
- utils: Unifying toBytes() conversions
- config: Adding '--work-dir' support in flavors
- config: Adding '--size' support in flavors
- config: Add "--data" support in flavors
- Makefile: Add tests in the prepare target
- start: Removing privileged option from command line
- image: Adding 'image show-aliases' command
- cn.toml: Updating title
- config: Moving high level functions from config_file to utils
- config: Using viper.IsSet() instead of custom code
- config: Adding "flavors {ls|show}" commands
- config: Hardcode more flavors & images
- config: Switching from panic() to log.Fatal()
- config: Adding [images] support
- config: Moving flavors into [flavors]
- config: Renaming MemorySize to memory_size
- config: Adding flavor types to containers
- main: Reporting if no configuration file found
- config: Report memory & cpu settings at start time
- config: Reporting configuration file used
- config: Implement cpu_count
- config: Reading ceph.conf parameters from config file
- config: Enforcing containerName for getMemorySize
- config: Implement use_default
- cmd: Adding configuration file support
- Readme: Bump the new release tag: v2.0.4
- Packaging: Update specfile version to v2.0.4
* Thu Nov 22 2018  Erwan Velu <evelu@redhat.com> - 2.0.4-1
- travis: Don't split edit & commit the README
- contrib/travis.sh: Fixing typo
- contrib/travis: Updating specfile when releasing
- Packaging: Adding rpm support
- Packaging: Simplifying versionning
- doc: add a 'build' section
- Bump README with the new release tag: v2.0.3
* Tue Nov 11 2018 Erwan Velu <evelu@redhat.com> - 2.0.3-1
- Initial Release
