# import required libraries
from PIL import Image
import sys, os

KEY_FRAMES = 10 if len(sys.argv) < 3 else int(sys.argv[3])
FILE_NAME = sys.argv[1]
DIR_NAME = os.path.splitext(os.path.basename(FILE_NAME))[0].upper()
SCALE = float(sys.argv[2])
TARGET = os.name == "nt" and ".\\ascii.exe" or "./ascii"

try:
    os.makedirs(DIR_NAME, exist_ok=True)
except Exception as err:
    print(str(err.with_traceback()))
with Image.open(FILE_NAME) as im:
    for i in range(KEY_FRAMES):
        im.seek(im.n_frames // KEY_FRAMES * i)
        tmp = os.path.join(DIR_NAME, f"{i}-temp.png")
        im.save(tmp)

for i in range(KEY_FRAMES):
    out_path = os.path.join(".", DIR_NAME, f"{i}.txt")
    in_path = os.path.join(".", DIR_NAME, f"{i}-temp.png")
    os.system(f"\"{TARGET}\" --out {out_path} --in {in_path} --scale {SCALE}")
    os.remove(in_path)

print(f"done, check directory \"{DIR_NAME}\".")
