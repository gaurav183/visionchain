import cv2
import subprocess as sp
import base64
import sys
import hashlib

# ffmpegCmd = "ffmpeg -i foo.h264 -frames:v 1 -f image2 frame.h264"

# vid_file = open('foo.h264')
# ffmpeg = sp.Popen(ffmpegCmd)


# fileSizeBytes = ffmpeg.stdout.read(6)
# print ffmpeg.stdout.read()
# fileSize = 0
# for i in xrange(4):
#     fileSize += fileSizeBytes[i + 2] * 256 ** i
# bmpData = fileSizeBytes + ffmpeg.stdout.read(fileSize - 6)
# image = cv2.imdecode(np.fromstring(bmpData, dtype = np.uint8), 1)
encoded_string = ""

for i in xrange(30):
	ffmpegCmd = "ffmpeg -i dash_2s.h264 -c:v libx264 -filter:v 'select=gte(n\,%d)' -frames:v 1 -f h264 frames/frame%d.h264"%(i,i)
	output = sp.check_output([ffmpegCmd], shell=True)

	framePath = "frames/frame%d.h264"%i
	with open(framePath, "rb") as image_file:
		encoded_string += base64.b64encode(image_file.read())


print "Hash MD5 = ", (hashlib.md5(encoded_string).hexdigest())
print "Hash SHA = ", (hashlib.sha512(encoded_string).hexdigest())

sys.stdout = open('dash2s30.txt', 'w')
print encoded_string
