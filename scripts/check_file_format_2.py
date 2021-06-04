from struct import unpack
import os


marker_mapping = {
        0xffd8: "Start of Image",
        0xffe0: "Application Default Header",
        0xffdb: "Quantization Table",
        0xffc0: "Start of Frame",
        0xffc4: "Define Huffman Table",
        0xffda: "Start of Scan",
        0xffd9: "End of Image"
}

class JPEG:
    def __init__(self, image_file):
        with open(image_file, 'rb') as f:
            self.img_data = f.read()

    def decode(self):
        data = self.img_data
        while(True):
            marker, = unpack(">H", data[0:2])
            if marker == 0xffd8:
                data = data[2:]
            elif marker == 0xffd9:
                return
            elif marker == 0xffda:
                data = data[-2:]
            else:
                lenchunk, = unpack(">H", data[2:4])
                data = data[2+lenchunk:]
            
            if len(data) == 0:
                break

bads = []
source_dir = './downloads/binary_data'
s_list = os.listdir(source_dir)
no_extension = 0
for klass in s_list:
    klass_path=os.path.join(source_dir, klass)
    print('processing class directory ', klass)
    if os.path.isdir(klass_path):
        file_list = os.listdir(klass_path)
        for f in file_list:
            f_path = os.path.join(klass_path, f)
            image = JPEG(f_path)
            try:
                image.decode()
            except:
                bads.append(f_path)


print("Number of bad images: {}".format(len(bads)))
for name in bads:
    os.remove(name)
