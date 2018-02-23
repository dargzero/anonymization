A program felépítese, használt fájlok:
A program python 3-ban íródott, a fõ modul az anonimyze.py python fájjlel indítható. 
Ha bemeneti paraméter nélkül hívjuk meg, akkor az alapértelmezett beállításokat és fájlokat használja a program, ezek a következõk:
	'-d','--DATADIR', az a könyvtár, ahova az adatokat tartalmazó fajlok kerülnek
	'-c','--CONFIGDIR', a beállításokat tartalmazó könyvtár
	'-o','--OUTDIR', az anonimizált adatok ide kerülnek kimentésre
	'-s','--STRUCTFILE', a bemeneti adat struktúráját írja le, és a használandó generalizáló algoritmusokat
	'-x','--XMLFILE' : az anonimizáláshoz használhatunk egy fát ami a generalizáláshoz szükséges struktúrát tartalmazza
	'-w','--WOKRDIR' : ebbe a könyvtárba másolódni ki a datadir-bõl a már feldolgozott fájlok
	

A configdir könyvtár tartalma:
	-	recordstruct.json : leírja a bemenetei adatszerkezetet(id mezõket, nem generalizálandó mezõket, és a generalizált mezõkhöz használt algoritmusokat)
	-	mod.py : ebbe a fájlba lettek kiszervezve a generalizáló algoritmusokat
	-	treedata.xml : azt a fa struktúrát tartalmazza ami alapján a szöveges mezõk általánosíthatóak, jelen esetben csak 1-et tud kezelni a program
	-	anonimData.txt : a program futása során keletkezik, azokat a csoportokat tartalmazza, amelyeket már anonimizáltunk, 
	és vagy teljesítik a k anonym feltételt, vagy nem 
	-	algorithm.txt : az elsõ futtatás után keletkezõ anonimizáló algoritmust tartalmazza, a késõbbi futtatásokban ezzel az algoritmussal dolgozik a program
	
Az outdir könyvtár tartalma:
	-	outfile*.txt : az anonimizálás végeredménye
	-	notanonymcoll.txt : segédfájl, ami tartamazza azokat a sorokat, amelyeket nem tudtuk még az eredményfájlba betenni, mert nem felelnek meg k anonymak

A teszt könyvtárban lévõ readme_teszt.txt-ben leírt beállítások után az anonimyze.py-vel indítható a program.