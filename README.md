# ABOUT
Learning Japense Kanji is hard enough, but what is so difficult is figuring out the readings. The purpose of this program is to allow you to enter a reading (romanized or in hiragana), to retrieve all permuations along with links to jisho.org.
# CURRENT FEATURES
The program currently imports 3 csv files, which are subsets of data analyzed on Kunyomi, Onyomi, and Kunyomi/Hiragana words. The program imports these 3 csv files into a hashmap and prompts the user for Hiragana/romaji input. Once input, this will print out all variants so you can further analyze / correlate things for memory and self study.
# HOW IT WORKS
Currently, the program uses multiple CSV files (cleaned, analyzed, and formatted) from analyzing the Wikipedia Joyo Kanji Table along with the PDF provided by the japanese affairs. 
# WORKS CITED
* **Wikipedia Table:**
 https://en.wikipedia.org/wiki/List_of_j%C5%8Dy%C5%8D_kanji
* **Kanji Table: https://www.bunka.go.jp/kokugo_nihongo/sisaku/joho/joho/kijun/naikaku/pdf/joyokanjihyo_20101130.pdf**
https://www.jisho.org - external links to kanji readings with very detailed information

# README.md UNDER CONSTRUCTION

# FEATURES TO ADD
* Make the program accessible through command line interface globally, where calling the name followed by desired reading will immediately print out results.
* Rewrite assistant tools being used in bash and golang.
* Concurrency in reading all 3 csv files in at the same time if possible

# FEEL FREE TO REPORT ANY ERRORS OR FORK OUT
