import csv, json
import sys

def convert_row(row):
	try:
		entry = {} #{'pseudo': '', 'email': '', 'notes': ''}
		entry['name'] = row['name']
		entry['password'] = row['password']
		uname = row['username']
		if '@' in uname:
			entry['email'] = uname 
		else:
			entry['pseudo'] = uname 
		if row['extra'] is not '': 
			entry['notes'] = row['extra']
		return entry
	except:
		print("error parsing row:", row)

if __name__ == '__main__':
	if len(sys.argv) == 1:
		print("missing input file.")
		exit(1)

	pathin = sys.argv[1]

	if not pathin.endswith("csv"):
		print("path should end with .csv.")

	if len(sys.argv) > 2:
		pathout = sys.argv[2]
	else:
		pathout = pathin.replace("csv", "json")

	print("reading from %s..." % pathin)

with open(pathin, 'rb') as fin:
	reader = csv.DictReader(fin, delimiter=',', quotechar='"')
	dico = [ convert_row(r) for r in reader ]

	with open(pathout, 'w') as fout:
		json.dump(dico, fout, indent=4)
		print("exported %d entries to %s." % (len(dico), pathout))



