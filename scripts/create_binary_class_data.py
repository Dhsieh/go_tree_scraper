import collections
import os
import json
import shutil

END_DIR = "./downloads/binary_data"
DATA_DIR = "./downloads/bing"
TREE_SPECIES_JSON = "./downloads/tree_data.json"
all_species = os.listdir(DATA_DIR)

tree_dict = collections.defaultdict(str)
tree_type_dict = collections.defaultdict(str)
tree_type_count_dict = collections.defaultdict(int)
tree_type_set = set()

with open(TREE_SPECIES_JSON) as f:
    tree_json = json.load(f)['Data']

for tree in tree_json:
    tree_dict[tree['ScientificName']] = tree['TreeType']
    tree_type_set.add(tree['TreeType'])

try:
    os.mkdir(END_DIR)

except:
    print("{} has already been created".format(END_DIR))

for tree_type in tree_type_set:
    tree_type_dir = "{}/{}".format(END_DIR, tree_type)
    tree_type_dict[tree_type] = tree_type_dir
    tree_type_count_dict[tree_type] = 0
    try:
        os.mkdir(tree_type_dir)
    except:
        print("{} has already been created!".format(tree_type_dir))

for species in all_species:
    species_dir = "{}/{}".format(DATA_DIR, species)
    fixed_species = str.replace(species, "_", " ")
    tree_type = tree_dict[fixed_species]
    photos = os.listdir(species_dir)
    for photo in photos:
        if tree_type_count_dict[tree_type] == 500:
            break
        tree_type_count_dict[tree_type] = tree_type_count_dict[tree_type] + 1
        photo_file = "{}/{}".format(species_dir, photo)
        new_photo_File = "{}_{}".format(species, photo)
        photo_dest = "{}/{}".format(tree_type_dict[tree_type], new_photo_File)
        shutil.copyfile(photo_file, photo_dest)

print(tree_type_count_dict)
