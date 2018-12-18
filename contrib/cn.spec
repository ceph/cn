%global source_version 2.0.4
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
install -D -p -m 644 cn.toml %{buildroot}%{_sysconfdir}/cn/
install -D -p -m 644 contrib/cn_completion.sh %{buildroot}%{_sysconfdir}/bash_completion.d/

%files
%doc README.md
%{_bindir}/cn
%{_sysconfdir}/cn/cn.toml
%{_sysconfdir}/bash_completion.d/cn_completion.sh

%changelog
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
