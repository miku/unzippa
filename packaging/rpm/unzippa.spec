Summary:    Unzip members of a zipfile faster than unzip.
Name:       unzippa
Version:    0.1.4
Release:    0
License:    GPL
BuildArch:  x86_64
BuildRoot:  %{_tmppath}/%{name}-build
Group:      System/Base
Vendor:     Leipzig University Library, https://www.ub.uni-leipzig.de
URL:        https://github.com/miku/unzippa

%description

Unzip members of a zipfile faster than unzip.

%prep

%build

%pre

%install

mkdir -p $RPM_BUILD_ROOT/usr/local/sbin
install -m 755 unzippa $RPM_BUILD_ROOT/usr/local/sbin
install -m 755 unzippall $RPM_BUILD_ROOT/usr/local/sbin

%post

%clean
rm -rf $RPM_BUILD_ROOT
rm -rf %{_tmppath}/%{name}
rm -rf %{_topdir}/BUILD/%{name}

%files
%defattr(-,root,root)

/usr/local/sbin/unzippa
/usr/local/sbin/unzippall

%changelog

* Mon Apr 9 2018 Martin Czygan
- 0.1.1 initial release
