import csv

from collections import defaultdict

myarray = []

selectedname = ''

selectedFile = input("Enter Selected Name\n1. KanjiFrequencyListKunyomi.csv\n2. KanjiFrequencyListOnyomi.csv\n3. KunyomiWithHiragana.csv\n")

if selectedFile == '1':
    selectedname = 'KanjiFrequencyListKunyomi.csv'

elif selectedFile == '2':
    selectedname = 'KanjiFrequencyListOnyomi.csv'

elif selectedFile == '3':
    selectedname = 'KunyomiWithHiragana.csv'

# import designated csv to test for duplicates
with open(selectedname, 'r', encoding = 'utf-8') as csvfile:
    kanjireader = csv.reader(csvfile, delimiter = '\t')

    for row in kanjireader:
        myarray.append(row[0])


# create default dict, where the value is going to be an int
dupdict = defaultdict(int)

# we will iterate through every value in the array
for currentValue in myarray:

    # we will add into the duplicate dict, which handles does not exist errors
    dupdict[currentValue] += 1

# iterate through the map through key and value, and print out the duplicates
for index, value in dupdict.items():
    if value > 1:
        print(index, ' -- ', value)
