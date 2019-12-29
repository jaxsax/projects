import os

def main():
    files = []
    # r=root, d=directories, f = files
    for r, d, f in os.walk('.'):
        for file in f:
            files.append(os.path.join(r, file))

    for f in files:
        print(f)

if __name__ == "__main__":
    main()