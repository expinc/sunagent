import base64
import distro
import json
import os
import platform
import subprocess
from assertpy import assert_that

HOST = "127.0.0.1"
PORT = 5000

# where test data is read
TEST_DATA_DIR = "dir"
TEST_DATA_PATH = os.path.join(os.getcwd(), "test", "functionality", "data", TEST_DATA_DIR)

# where temporary file directory
TEST_TMP_DIR = os.path.join(os.getcwd(), "tmp")

# sample: "2021-07-31T22:44:17.7724489+08:00" or "2021-07-31T22:44:17.7724489Z" RFC3339
TIMESTAMP_PATTERN = r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|(\+\d{2}:\d{2}))$"
# sample: "-rw-rw-rw-"
REGULAR_FILE_MODE_PATTERN = r"^-((r|-)(w|-)(x|-)){3}$"
# sample: "drwxrwxrwx"
DIRECTORY_MODE_PATTERN = r"^d((r|-)(w|-)(x|-)){3}$"

LINUX_DIR_SIZE = 4096

# dummy python process script
TEST_EXE_PATH = os.path.join(os.getcwd(), "test", "functionality", "exe")
TEST_DUMMY_PROC = os.path.join(TEST_EXE_PATH, "dummy-proc.py")

# there are 2 python processes for test script:
# 1. test.py
# 2. pytest
TEST_SCRIPT_PYTHON_PROC_COUNT = 2

# OS family to package {name, file, URL, newer-file, newer-URL}
TEST_PKGS = {
    "ubuntu": {
        "name": "glibc-doc-reference",
        "file": "glibc-doc-reference_2.30-1ubuntu1_all.deb",
        "url": "http://archive.ubuntu.com/ubuntu/pool/main/g/glibc-doc-reference/glibc-doc-reference_2.30-1ubuntu1_all.deb",
        "newFile": "glibc-doc-reference_2.33-0ubuntu1_all.deb",
        "newUrl": "http://archive.ubuntu.com/ubuntu/pool/main/g/glibc-doc-reference/glibc-doc-reference_2.33-0ubuntu1_all.deb"},
    "debian": {
        "name": "glibc-doc",
        "file": "glibc-doc_2.24-11+deb9u4_all.deb",
        "url": "http://ftp.br.debian.org/debian/pool/main/g/glibc/glibc-doc_2.24-11+deb9u4_all.deb",
        "newFile": "glibc-doc_2.28-10_all.deb",
        "newUrl": "http://ftp.br.debian.org/debian/pool/main/g/glibc/glibc-doc_2.28-10_all.deb"},
    "centos": {
        "name": "gdb-doc",
        "file": "gdb-doc-7.6.1-120.el7.noarch.rpm",
        "url": "http://mirror.centos.org/centos/7/os/x86_64/Packages/gdb-doc-7.6.1-120.el7.noarch.rpm",
        "newFile": "gdb-doc-8.2-15.el8.noarch.rpm",
        "newUrl": "http://mirror.centos.org/centos/8/AppStream/aarch64/os/Packages/gdb-doc-8.2-15.el8.noarch.rpm"},
    "opensuse": {
        "name": "git-doc",
        "file": "git-doc-2.26.2-lp152.2.12.1.noarch.rpm",
        "url": "https://ftp.lysator.liu.se/pub/opensuse/update/leap/15.2/oss/noarch/git-doc-2.26.2-lp152.2.12.1.noarch.rpm",
        "newFile": "git-doc-2.33.1-1.1.noarch.rpm",
        "newUrl": "https://ftp.lysator.liu.se/pub/opensuse/ports/aarch64/tumbleweed/repo/oss/noarch/git-doc-2.33.1-1.1.noarch.rpm"},
    "opensuse-leap": {
        "name": "git-doc",
        "file": "git-doc-2.26.2-lp152.2.12.1.noarch.rpm",
        "url": "https://ftp.lysator.liu.se/pub/opensuse/update/leap/15.2/oss/noarch/git-doc-2.26.2-lp152.2.12.1.noarch.rpm",
        "newFile": "git-doc-2.33.1-1.1.noarch.rpm",
        "newUrl": "https://ftp.lysator.liu.se/pub/opensuse/ports/aarch64/tumbleweed/repo/oss/noarch/git-doc-2.33.1-1.1.noarch.rpm"},
}

BASIC_AUTH_TOKEN = base64.b64encode(b"admin:root").decode("ascii")

def assert_successful_response(response, status, data=None):
    assert_that(response.status).is_equal_to(status)

    body = response.read()
    body = json.loads(body)

    assert_that(body["status"]).is_equal_to(status)
    assert_that(body["successful"]).is_equal_to(True)
    if data:
        assert_that(body["data"]).is_equal_to(data)

def assert_failed_response(response, status, data=None):
    assert_that(response.status).is_equal_to(status)

    body = response.read()
    body = json.loads(body)

    assert_that(body["status"]).is_equal_to(status)
    assert_that(body["successful"]).is_equal_to(False)
    if data:
        assert_that(body["data"]).is_equal_to(data)

def get_successful_response(response, status):
    assert_that(response.status).is_equal_to(status)
    body = response.read()
    body = json.loads(body)
    assert_that(body["successful"]).is_equal_to(True)
    return body["data"]

def get_failed_response(response, status):
    assert_that(response.status).is_equal_to(status)
    body = response.read()
    body = json.loads(body)
    assert_that(body["successful"]).is_equal_to(False)
    return body["data"]

def get_binary_response(response, status):
    assert_that(response.status).is_equal_to(status)
    return response.read()

def install_package(name_or_path):
    if "ubuntu" == distro.id() or "debian" == distro.id():
        return 0 == os.system("apt install -y {}".format(name_or_path))
    elif "centos" == distro.id():
        return 0 == os.system("yum -y install {}".format(name_or_path))
    elif "opensuse" in distro.id():
        return 0 == os.system("zypper -n install {}".format(name_or_path))
    else:
        raise Exception("Not supported")

def remove_package(name):
    if "ubuntu" == distro.id() or "debian" == distro.id():
        return 0 == os.system("dpkg -r {}".format(name))
    elif "centos" == distro.id() or "opensuse" in distro.id():
        return 0 == os.system("rpm -e {}".format(name))
    else:
        raise Exception("Not supported")

def is_package_installed(name):
    if "ubuntu" == distro.id() or "debian" == distro.id():
        return 0 == os.system("dpkg -s {}".format(name))
    elif "centos" == distro.id() or "opensuse" in distro.id():
        return 0 == os.system("rpm -qi {}".format(name))
    else:
        raise Exception("Not supported")

def download_file(url, target_path):
    if "Linux" == platform.system():
        return 0 == os.system("wget -P {} {}".format(target_path, url))
    else:
        raise Exception("Not supported")

def get_package_version(name):
    if "ubuntu" == distro.id() or "debian" == distro.id():
        output = subprocess.check_output("dpkg -s {} | grep Version".format(name), shell=True)
    elif "centos" == distro.id() or "opensuse" in distro.id():
        output = subprocess.check_output("rpm -qi {} | grep Version".format(name), shell=True)
    else:
        raise Exception("Not supported")

    # format of output should be "Version : x.y.z"
    return str(output).split(": ")[1]
