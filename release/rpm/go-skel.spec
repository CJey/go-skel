# rpm spec
%define name replace-go-skel
%define version replace-0.0.1
%define release replace-1
%define pjroot replace-/home/cjey/git/go-skel
%define runroot replace-/opt/%{name}

%global _enable_debug_package 0
%global debug_package %{nil}
%global __os_install_post /usr/lib/rpm/brp-compress %{nil}

Name: %{name}
Version: %{version}
Release: %{release}
Vendor: cjey.hou@gmail.com
License: MIT
Requires: logrotate
Summary: Go skeleton

%description
Go skeleton

%prep
cp %{pjroot}/release/sysvinit sysvinit
cp %{pjroot}/release/logrotate.conf logrotate.conf
cp %{pjroot}/release/%{name}.toml config.toml
echo %{version} > version

%build
MAINFILE="%{pjroot}/main.go" %{pjroot}/build %{name} bin

sed -i 's;${{INSTALL_ROOT}};%{runroot};g' sysvinit logrotate.conf
sed -i 's;^\(\s*file\s*=\s*"\)stderr\(".*\)$;\1%{runroot}/log/%{name}.log\2;' config.toml

%install
# systemd
mkdir -p %{buildroot}/etc/init.d/
cp sysvinit %{buildroot}/etc/init.d/%{name}

# logrotate
mkdir -p %{buildroot}/etc/logrotate.d
cp logrotate.conf %{buildroot}/etc/logrotate.d/%{name}

# package
mkdir -p %{buildroot}%{runroot}
cd %{buildroot}%{runroot}

mkdir -p bin conf log
install %{_builddir}/bin bin/%{name}
cp %{_builddir}/config.toml conf/%{name}.toml.orig
cp %{_builddir}/config.toml conf/%{name}.toml
cp %{_builddir}/version .

%files
/etc/init.d/%{name}
/etc/logrotate.d/%{name}

%{runroot}/version
%{runroot}/bin/%{name}
%{runroot}/conf/%{name}.toml.orig

%config(noreplace) %{runroot}/conf/%{name}.toml

%dir %{runroot}
%dir %{runroot}/bin
%dir %{runroot}/log
%dir %{runroot}/conf

# %pre, %post, %preun, %postun
# 这四个脚本都会有一个参数，即$1
# $1 表明的是执行本脚本时，系统上此包有几个，但值得注意的是
# 安装阶段，$1可以是1或2，分别表示新安装或者升级
# 卸载阶段，$1可以是0或1，分别表示卸载或者升级

%pre

%post
if [ $1 -eq 1 ]; then
    chkconfig %{name} on
fi

service %{name} reload

%preun
if [ $1 -eq 0 ]; then
    chkconfig %{name} off
    service %{name} stop
fi

%postun
