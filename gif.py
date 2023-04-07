# import required libraries
from PIL import Image
import sys, os

KEY_FRAMES = 10 if len(sys.argv) < 3 else int(sys.argv[3])
FILE_NAME = sys.argv[1]
DIR_NAME = FILE_NAME.split("/")[-1].split(".")[0].upper()
SCALE = float(sys.argv[2])

try:
    os.mkdir(DIR_NAME)
except Exception as err:
    print(str(err.with_traceback()))
with Image.open(FILE_NAME) as im:
    for i in range(KEY_FRAMES):
        im.seek(im.n_frames // KEY_FRAMES * i)
        im.save(f"{DIR_NAME}/{i}-temp.png")

for i in range(KEY_FRAMES):
    os.system(f"./ascii --out {DIR_NAME}/{i}.txt --in {DIR_NAME}/{i}-temp.png --scale {SCALE}")
    os.remove(f"{DIR_NAME}/{i}-temp.png")

print(f"done, check directory \"{DIR_NAME}\".")
