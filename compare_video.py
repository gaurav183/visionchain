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

for i in xrange(10):
	for j in xrange(30):
		ffmpegCmd = "ffmpeg -i dash_cam.mp4 -c:v libx264 -filter:v 'select=gte(n\,%d)' -frames:v 1 -f h264 frames/frame%d.h264"%(j,j)
		output = sp.check_output([ffmpegCmd], shell=True)

		framePath = "frames/frame%d.h264"%j
		with open(framePath, "rb") as image_file:
			encoded_string += base64.b64encode(image_file.read())

	md5Hash = (hashlib.md5(encoded_string).hexdigest())
	print "Hash MD5 = ", md5Hash
	shaHash = (hashlib.sha512(encoded_string).hexdigest())
	print "Hash SHA = ", shaHash

	# send POST req to blockchain with shaHash, id = "dash_%d"%i

	encoded_string = ""




sys.stdout = open('dash30.txt', 'w')
print encoded_string

# 579d88bf0a07f3671427e2bf6d42fbe99420b8fb7ebf32ce4f24153b41d5b6c1b3cce86da0f53c1b0c1485df753bb0e21261adb6bc73dc7ef6428667ee19f61b