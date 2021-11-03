import distro
import json
import os
import platform
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

# OS family to package [name, file, URL]
TEST_PKGS = {
    "ubuntu": ["glibc-doc-reference", "glibc-doc-reference_2.30-1ubuntu1_all.deb", "http://archive.ubuntu.com/ubuntu/pool/main/g/glibc-doc-reference/glibc-doc-reference_2.30-1ubuntu1_all.deb"],
    "debian": ["glibc-doc", "glibc-doc_2.24-11+deb9u4_all.deb", "http://ftp.br.debian.org/debian/pool/main/g/glibc/glibc-doc_2.24-11+deb9u4_all.deb"],
    "centos": ["gdb-doc", "gdb-doc-7.6.1-120.el7.noarch.rpm", "http://mirror.centos.org/centos/7/os/x86_64/Packages/gdb-doc-7.6.1-120.el7.noarch.rpm"],
    "opensuse": ["git-doc", "git-doc-2.26.2-lp152.2.12.1.noarch.rpm", "https://ftp.lysator.liu.se/pub/opensuse/update/leap/15.2/oss/noarch/git-doc-2.26.2-lp152.2.12.1.noarch.rpm"],
}

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
