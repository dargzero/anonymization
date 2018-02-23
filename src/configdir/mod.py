import os
import xml.etree.ElementTree as ET

# a fa feldolgozása a nemzetiségek általánosításához

scriptDir = os.path.dirname(os.path.abspath(__file__))
path = os.path.join(scriptDir, 'treedata.xml')
tree = ET.parse(path)
root = tree.getroot()

def Height(node): 
	height = 0
	for child in node :
		height = max(height, Height(child))
	return height + 1
	

def Hash( data) :  
	return hash(data)

#def SuppressData( data, hashData ):    
#	
#	toHash = ''
#	toHash = str(hashData)
#	supprData = []
#	retData = []
#	toHash = toHash + data
#	supprData.append('*')
#	retData.append(toHash)
#	retData = retData + supprData
#	return retData

def SuppressData( data, n ):    
	return "*"
	
def GeneralizeDataByTruncating( data , n ):    
	if len(data) < n :
		return "*"*n
	retData = data[:len(data)-n] + "*"*n
	return retData


def GeneralizeDataByCategorizingInt( data = 0, n = 1, max = 100 ):    
	dataInt = 0
	if '..' in data :
		items = data.split("<=")
		dataInt = (int(items[0]) + int(items[2]) ) / 2
	elif '<' in data:
		return data
	else :
		dataInt = int(data)
	step = 2**(n-1) * 10
	
	if int(max) < dataInt or step > int(max):
		return str(max) + "< "
	
	
	if dataInt > 1 :
		dataInt = dataInt - 1
	min = int((dataInt // step) * step) + 1
	max = int((dataInt // step + 1) * step)
	if min == 1 :
		min = 0
	retData = '{}<=..<={}'.format(min,max)
	return retData
	
def GetParent( node, child):
	if node.find(child) != None :
		return node
	else:
		for f in node :
			found = GetParent( f, child)
			if found != None :
				return found

def GetElement(node, data) :
	if node.tag == data:
		return node
	else:
		for f in node :
			found = GetElement( f, data)
			if found != None :
				return found
		return None

def GeneralizeDataByCategorizingTree( data, n ):    
	tempNode = GetElement(root,data)
	i = Height(tempNode)
	while i <= n :
		node = GetParent(root,data)
		data = node.tag
		i = i + 1
	
	return data
	
