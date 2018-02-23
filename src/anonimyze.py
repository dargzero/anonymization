#!/usr/bin/env python3
# -*- coding: Utf-8 -*-
# new.py -i tesztdata.txt -x treedata.xml -s recordstruct.txt

import argparse
import shlex
import csv
import sys
import xml.etree.ElementTree as ET
import os
import time
import msvcrt
import json
import importlib
import shutil

parser = argparse.ArgumentParser(description='Example for arguments.')
parser.add_argument('-x','--XMLFILE', help='Az xml bemeneti fajl eleresi utja.', type=str, required=False, default = 'treedata.xml')
parser.add_argument('-s','--STRUCTFILE', help='A rekordstruktúra bemeneti fajl eleresi utja.', type=str, required=False, default = 'recordstruct.json')
parser.add_argument('-w','--WOKRDIR', help='A munkakönyvtár ahova az ideiglenes fájlokat rakja.', type=str, required=False, default='worktempdir')
parser.add_argument('-d','--DATADIR', help='Az temp könyvtár ahova az készfájlokat rakja.', type=str, required=False, default='datadir')
parser.add_argument('-c','--CONFIGDIR', help='A config könyvtár ahol az beállító fájlok vannak.', type=str, required=False, default='configDir')
parser.add_argument('-o','--OUTDIR', help='A kimeneti könyvtár ahova az anonim fájlok kerülnek.', type=str, required=False, default='outdir')



args = parser.parse_args()

#inicializálás

inputList = []  # bemeneti lista
anonymList = [] # anonimizált adatok
notKAnonymColl = [] # nem anonimizálandó adatok
wasAlgorithm = False # volt-e algoritmus megadva
module = {} # az általánosításhoz megadott algoritmusok oszoloponként
kAnonym = 4 # anonimizálás szintje



scriptDir = ''
scriptDir = os.path.dirname(os.path.abspath(__file__))
anonimData = {} # a kulcsok a sorok, azokból hány azonos van
createLog = True

# könyvtár létrehozása ha kell
outDir = os.path.join(scriptDir, args.OUTDIR)
try:
    os.makedirs(outDir)
except OSError:
    pass # already exists

# könyvtár létrehozása ha kell
workDir = os.path.join(scriptDir, args.WOKRDIR)
try:
    os.makedirs(workDir)
except OSError:
    pass # already exists


# könyvtár létrehozása ha kell
dataDir = os.path.join(scriptDir, args.DATADIR)
try:
    os.makedirs(dataDir)
except OSError:
    pass # already exists


configDir = os.path.join(scriptDir, args.CONFIGDIR)
# nem anonimizált adatok fájlja

file_name = 'notanonymcoll.txt'
collpath = os.path.join(outDir, file_name)
	
# adatstruktúra beolvasása listákba
needGenCols = set()   # azon oszlopok melyeken módosítunk
idColls = set() # id oszlopok
hashColls = set() # az hash-elt azonosító oszlopa
genmethmod = ""
genmeth = ""

if args.STRUCTFILE is not None :
	
	try:
		path = os.path.join(configDir, args.STRUCTFILE)
		with open(path) as json_file:
			json_data = json.load(json_file)
		for col in json_data:
			
			for colattr in json_data[col] :
				for attr in colattr :
					if attr == "isId" and colattr[attr] == True :
						idColls = idColls|{col}
					if attr == "isHash" and colattr[attr] == True :
						hashColls = hashColls|{col}
					if attr == "needgen" and colattr[attr] == True :
						needGenCols = needGenCols|{col}
					if attr == "genmethod" and col in needGenCols :
						for genattr in colattr[attr]:
							for genmethod in genattr :
								if genmethod == "module" :
									genmethmod = genattr[genmethod]
								if genmethod == "method" :
									genmeth = genattr[genmethod]
						module[int(col)] = [genmethmod,genmeth]				
	except IOError as e:
		errno, strerror = e.args
		print("ERROR: be/kimeneti hiba({0}): {1}".format(errno,strerror))
		sys.exit(0)
	except:
		print("ERROR: egyéb hiba:", sys.exc_info()[0])
		sys.exit(0)



			

def CreateLogFile ( prefix, data, isset = False  ) :
	
	logfile_name = prefix + str(time.time()).replace('.','') + '.txt'
	logpath = os.path.join(workDir, logfile_name)
	logfile = open(logpath, 'w')
	for line in data :
		if isset == True :
			logfile.write(str(line) + ':' + str(data[line]) + '\n')
		else :
			logfile.write(str(line) + '\n')
	logfile.close()
	

# az anonym csoportok mentése
algorithm = {}
def SaveAnonimData( ):    # 
	path = os.path.join(configDir, 'anonimData.txt')
	f = open(path, 'w')
	json.dump(anonimData,f)
	f.close()

# az anonym csoportok betöltése
def LoadAnonimData( ):    # 
	try:
		path = os.path.join(configDir, "anonimData.txt")
		with open(path) as json_file:
			anonimData = json.load(json_file)
		return anonimData
	except FileNotFoundError as e:
		return {}
		pass

# eddig nem kimentett adatok listájának betöltése
def LoadNotKAnonymColl () :
	try:
		
		collfile = open(collpath, 'r')
		notKAnonymColl = []
		for data in collfile :
			anonymLine = []
			anonymLineTmp = []
			anonymLine = data.split(";")
			for col in anonymLine :
				anonymLineTmp.append(col.strip(" "))
			
			notKAnonymColl.append(anonymLineTmp)
		collfile.close()
		return notKAnonymColl
	except FileNotFoundError as e:
		pass
		

	
# anonimizáló algoritmus mentése
def SaveAlgorithm( ):    # 
	path = os.path.join(configDir, 'algorithm.txt')
	f = open(path, 'w')
	json.dump(algorithm,f)
	f.close()

# anonimizáló algoritmus betöltése
def LoadAlgorithm( ):    # 
	path = os.path.join(configDir, "algorithm.txt")
	with open(path) as json_file:
		algorithm = json.load(json_file)
	return algorithm
	
	
try:

	
	algorithm = LoadAlgorithm()

	wasAlgorithm = True
except FileNotFoundError as e:
	for i in range(len(json_data)) :
		algorithm[str(i)] = 0 # az eddigi anonimizáló algoritmus [[oszlop, iteráció]]	
	pass


if wasAlgorithm == True :
	anonimData = LoadAnonimData()
	notKAnonymColl = LoadNotKAnonymColl()
	
	
# statisztika a megfelelő csoportok létrehozására, k anonimitás ellenőrzésére
def Statistics( dataListStat, iter = '-' ):    
	result = []
	statData = [set() for i in range(len(dataListStat[0]))] # melyik oszlopban mennyi különböző érték van
	line = str
	anonimDataTemp = {}
	for k in anonimData :
		anonimDataTemp[k] = anonimData[k]
	for dataLine in dataListStat :
		line = ''.join(dataLine[(len(hashColls) + len(idColls)):len(needGenCols)])
		if not( line in anonimDataTemp) :
			anonimDataTemp[line] = 1
		else :
			anonimDataTemp[line] = anonimDataTemp[line] + 1
		for i in range(len(dataLine)) :
			statData[i].add(dataLine[i])
	
	anonymity = 0

	for k,v in anonimDataTemp.items() :
		if v < anonymity or anonymity == 0 :
			anonymity = v
	if createLog == True :
		prefix = 'statdata'
		if iter != '-' :
			prefix = prefix + '_' + iter + '_'
		CreateLogFile(prefix, statData)
		prefix = 'anonimDataTemp'
		if iter != '-' :
			prefix = prefix + '_' + iter + '_'
		CreateLogFile(prefix, anonimDataTemp, True)
	result.append(anonymity)
	result.append(statData)
	result.append(anonimDataTemp)

	return result 
	
def CreateAnonymWOAlg ( dataSet ) :
	curAnonymity = 0 # jelenlegi anonimitás szintje
	maxItems = 0 # a legnagyobb elemszámú csoport, amit anonimizálni kell
	statResult = [] # statisztika gyűjtés végeredménye
	anonymDataSet = dataSet # anonimizált adathalmaz
	statResult = Statistics( anonymDataSet, 'pre' )
	curAnonymity = statResult[0]  
	anonymCol = 0 # anonimizálandó oszlop
	for hash in hashColls : # Csak 1 lehet benne 
		intHash = int(hash)
		for i in range(len(anonymDataSet)) :	
			
			for id in idColls :
				intId = int(id)
				anonymDataSet[i][intHash] = anonymDataSet[i][intHash] + anonymDataSet[i][intId]
				method_to_call = getattr(importlib.import_module(module[intId][0]),module[intId][1])
				retData = method_to_call(anonymDataSet[i][intId],anonymDataSet[i][intHash])
				anonymDataSet[i][intId] = retData
			method_to_call = getattr(importlib.import_module(module[intHash][0]),module[intHash][1])
			data = method_to_call(anonymDataSet[i][intHash])

			anonymDataSet[i][intHash] = data
			
	iteration = 0
	while curAnonymity < kAnonym :
		iteration = iteration + 1
		maxItems = 0
		for i in range(0,len(needGenCols) + 1):  # itt majd a beégetett cuccok helyett lehetne változót csinálni
			if not (str(i) in idColls ) and not (str(i) in hashColls) and (str(i) in needGenCols):
				if len(statResult[1][i]) > maxItems :
					maxItems = len(statResult[1][i])
					anonymCol = i 
		if maxItems == 1 :
			print("ERROR: Nem anonimizálható")
			sys.exit(0)
		algorithm[str(anonymCol)] = algorithm[str(anonymCol)] + 1
		
		#az egész kollekció 1 oszlopára megcsinálja a generalizálást
		
		for i in range(len(anonymDataSet)) :
			
			method_to_call = getattr(importlib.import_module(module[anonymCol][0]),module[anonymCol][1])
			anonymDataSet[i][anonymCol] = method_to_call(anonymDataSet[i][anonymCol],algorithm[str(anonymCol)])

			
		statResult = Statistics( anonymDataSet, str(iteration) )
		curAnonymity = statResult[0]
		if createLog == True :
			CreateLogFile('anonymDataSet_'+str(iteration) + '_', anonymDataSet)

	for k in statResult[2]:
		anonimData[k] = statResult[2][k]
	
	return anonymDataSet

# A már meglévő algoritmus használata
def CreateAnonymWAlg ( dataSet ) :

	anonymDataSet = dataSet # anonimizált adathalmaz
	anonymCol = 0 # anonimizálandó oszlop
	statResult = [] # statisztika gyűjtés végeredménye
	for hash in hashColls : # Csak 1 lehet benne 
		intHash = int(hash)
		for i in range(len(anonymDataSet)) :	
			
			for id in idColls :
				intId = int(id)
				anonymDataSet[i][intHash] = anonymDataSet[i][intHash] + anonymDataSet[i][intId]
				method_to_call = getattr(importlib.import_module(module[intId][0]),module[intId][1])
				retData = method_to_call(anonymDataSet[i][intId],anonymDataSet[i][intHash])
				anonymDataSet[i][intId] = retData
			method_to_call = getattr(importlib.import_module(module[intHash][0]),module[intHash][1])
			data = method_to_call(anonymDataSet[i][intHash])

			anonymDataSet[i][intHash] = data
	
	for j in range(1,len(algorithm)) :
		if (not j in idColls) and (not j in hashColls) and  (algorithm[str(j)] > 0):
			for i in range(len(anonymDataSet)) :
				method_to_call = getattr(importlib.import_module(module[j][0]),module[j][1])
				anonymDataSet[i][j] = method_to_call(anonymDataSet[i][j],algorithm[str(j)])
			

	statResult = Statistics( anonymDataSet )
	for k in statResult[2]:
		anonimData[k] = statResult[2][k]
	return anonymDataSet
	
# főprogi
# adatok beolvasása listákba
while True :
	paths = [os.path.join(dataDir,fn) for fn in next(os.walk(dataDir))[2]]

	now =time.time()
	if len(paths) == 0 :
		sys.exit(0)
	fileNum = 0
	anonymList = []

	for file in paths : 
		inputList = []
		try:
			infile = open(file, "r")
			for line in infile : 
				inputList.append(line.split(";"))
		except IOError as e:
			errno, strerror = e.args
			print("ERROR: be/kimeneti hiba({0}): {1}".format(errno,strerror))
			sys.exit(0)
		except:
			print("ERROR: egyéb hiba:", sys.exc_info()[0])
			sys.exit(0)
		infile.close()
		shutil.move(infile.name, workDir)
		anonymLine = []
		for data in inputList :
			anonymLine = []
			anonymLine.append('')
			anonymLine.extend(data[0:])
			anonymList.append(anonymLine)
	
	if len(anonymList) > 0 :
		
		if wasAlgorithm :
			## volt már algoritmus definiálva
			anonymList = CreateAnonymWAlg(anonymList)
		else :
			## volt már algoritmus definiálva
			anonymList = CreateAnonymWOAlg(anonymList)	
			SaveAlgorithm()
			wasAlgorithm = True
		
		if createLog == True :
			CreateLogFile('anonymList_new_', anonymList)
		if notKAnonymColl != None :
			for line in notKAnonymColl:
				anonymList.append(line)
		collfile = open(collpath, 'w')
	
		file_name = 'outfile' + str(time.time()).replace('.','') + '.txt'
		outpath = os.path.join(outDir, file_name)
		anonym = open(outpath, 'w')
		notKAnonymColl = []
		if createLog == True :
			CreateLogFile('anonymList_all_', anonymList)
		for line in anonymList :
			anonimLine = ''.join(line[(len(hashColls) + len(idColls)):len(needGenCols)])
			if anonimData[anonimLine] >= kAnonym :
				anonym.write('; '.join(map(str, line)))
			else :
				notKAnonymColl.append(line)
				collfile.write('; '.join(map(str, line)))
		collfile.close()
		if len(notKAnonymColl) == 0 :
			os.remove(collpath)   
		anonym.close()
		# az üres fájlt töröljük
		statinfo = os.stat(outpath)
		if statinfo.st_size == 0 :
			os.remove(outpath)
		anonymList = []
		SaveAnonimData()
	later = time.time()
	# Idő mérésre
	#print(time.ctime(now))
	#print(time.ctime(later))
	

