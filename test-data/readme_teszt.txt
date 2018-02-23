A program teszteléséhez az egy könytárban lévõ adatfájlokat kell használni. 
Ahhoz hogy a tesztek megfelelõen fussanak alap helyzetbe kell állítani a programot, az a configdir tartalma csak a következõ lehet:
 - mod.py
 - treedata.xml
 - az általunk bemásolt szerkezetet leíró json fájl( ha átnevezzük recordstruct.json-rõl, akkor a scriptnek -s kapcsolóval kell megadni a fájlnevet) 

A külön könyvtárakban lévõ json fájlok írják le az adatszerkezetet, a teszthez a megfelelõ könyvtárban lévõ json fájl a configdir 
mappába kell másolni, a mellette lévõ csv állmonányt ami az adatokat tartalmazza pedig a datadir könyvtárba. 
Az anonimyze.py-vel indítható a program, a readme.txt tartalmazza a részletes leírást.