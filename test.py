#!/usr/bin/env python3

import argparse
import os
import platform
import psutil
import subprocess

TEST_OUT_DIR = os.path.join("gen", "test")

def parse_args():
    parser = argparse.ArgumentParser(description='Run tests')
    parser.add_argument("-t", "--type", required=False, choices=["unit", "func"])
    parser.add_argument("-p", "--part", required=False)

    return parser.parse_args()

def kill_proc(name):
    for proc in psutil.process_iter(["name"]):
        if proc.info["name"] == name:
            proc.kill()

def is_process_running_by_name(proc_name):
    for proc in psutil.process_iter(["name"]):
        if proc.info["name"] == proc_name:
            return True
    return False

def unit_test(part):
    print("=====================")
    print("Starting unit test...")
    cover_profile = os.path.join(TEST_OUT_DIR, "cover-profile")
    if part:
        cmd = "go test -v {} -coverpkg=./... ./... -coverprofile {}".format(part, cover_profile)
    else:
        cmd = "go test -v -coverpkg=./... ./... -coverprofile {}".format(cover_profile)
    os.system(cmd)

    print("Generating coverage report...")
    coverage_report = os.path.join(TEST_OUT_DIR, "unit-test.html")
    cmd = "go tool cover -html {} -o {}".format(cover_profile, coverage_report)
    os.system(cmd)
    print("See {} for coverage report".format(coverage_report))

    print("Unit test finished.")

def func_test(part):
    print("=====================")
    print("Starting functionality test...")

    print("Making executable...")
    sys_type = platform.system()
    if "Linux" == sys_type:
        subprocess.check_call(["sh", "make.sh"])
    elif "Windows" == sys_type:
        subprocess.check_call(["make.bat"])
    else:
        raise Exception("Unsupported platform: {}".format(sys_type))

    print("Starting executable...")
    sys_type = platform.system()
    if "Linux" == sys_type:
        process_name = "sunagentd"
        subprocess.check_call("gen/sunagentd --config=gen/config.conf &", shell=True)
    elif "Windows" == sys_type:
        process_name = "sunagentd.exe"
        subprocess.check_call("start gen\sunagentd.exe --config=gen/config.conf", shell=True)
    if not is_process_running_by_name(process_name):
        raise Exception("Executable failed to start")

    print("Running test cases...")
    cmd = "python3 -m pytest --capture=tee-sys"
    if part:
        cmd += " " + part
    try:
        subprocess.check_call(cmd, shell=True)
        print("Test succeeded")
    except subprocess.CalledProcessError as err:
        print("Test failed!")

    print("Stopping executable...")
    if "Linux" == sys_type:
        kill_proc("sunagentd")
    elif "Windows" == sys_type:
        kill_proc("sunagentd.exe")

if __name__ == "__main__":
    args = parse_args()
    os.makedirs(TEST_OUT_DIR, exist_ok=True)
    os.environ["TEST_ARTIFACT_DIR"] = os.path.join(os.getcwd(), "test")
    if "unit" == args.type:
        unit_test(args.part)
    elif "func" == args.type:
        func_test(args.part)
    else:
        unit_test(args.part)
        func_test(None)
