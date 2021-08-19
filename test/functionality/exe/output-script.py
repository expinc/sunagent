import sys

if __name__ == "__main__":
    print("stdout 1", file=sys.stdout, flush=True)
    print("stderr 1", file=sys.stderr, flush=True)
    print("stdout 2", file=sys.stdout, flush=True)
    print("stderr 2", file=sys.stderr, flush=True)
