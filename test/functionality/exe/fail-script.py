import sys

if __name__ == "__main__":
    print("start script", file=sys.stdout, flush=True)
    print("exit with 1", file=sys.stderr, flush=True)
    sys.exit(1)
