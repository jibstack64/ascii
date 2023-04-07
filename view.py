import sys, os
import time

DIRECTORY = sys.argv[1]

frames = []
for x in range(0, len(os.listdir(DIRECTORY))):
    with open(f"{DIRECTORY}/{x}.txt", "r") as f:
        frames.append(f.read())
        f.close()

while True:
    for f in frames:
        os.system("cls" if os.name == "nt" else "clear")
        print(f)
        time.sleep(0.25)
