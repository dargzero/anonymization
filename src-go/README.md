# Anonimizálás big data környezetben

Diplomaterv mellékletének használati útmutatója

Készítette: Martinek Vilmos

## KÖRNYEZET

Mivel az anonimizáló eszköz elkészítéséhez a Docker eszközt alkalmaztam, így a futtatáshoz az egyetlen szükséges szoftver a Docker. Minden más szükséges komponest automatikusan letölt és használ. Alapvetően a Docker működik Windows és Linux környezetben is, azonban én az előbbin teszteltem. Fontos azonban, hogy Windows környezetben is a Docker Linux konténerek használatára legyen beállítva – ugyanis maga az eszköz Linux alapú konténerekben fut.

A Docker eszköz használatához a Docker Communitiy Edition letöltése és telepítése szükséges, amit a következő oldalról lehet beszerezni: <https://www.docker.com/community-edition#/download>.

## ESZKÖZ FUTTATÁSA

Az anonimizáló eszköz futtatásához egy PowerShell ablak megnyitása után navigáljunk el a mellékletben található AnonymizationServer mappába, és adjuk ki a következő parancsot.

```bash
docker-compose up -d --build
```

Ez a parancs buildeli és teszteli az anonimizáló eszközt, majd pedig elindítja azt, az adatbázisszerverrel együtt. Ezen a ponton a REST API elérhető a localhost 9137-es porton, amit tetszőleges kliens – például Postman – segítségével tesztelhetünk. Amennyiben az adatbázisszervert szeretnénk elérni, az a MongoDB esetében alapértelmezett localhost 27017-es porton található.

Amennyiben le szeretnénk állítani az eszközt, a következő parancsot kell kiadnunk. Ez a parancs eltávolítja az adatbázisszerver által tárolt összes adatot, így a következő indulásnál tiszta lappal indulunk. Amennyiben nem ez a kívánt működés, akkor hagyjuk el a végéről a ’-v’ kapcsolót.

```bash
docker-compose down –v
```

## INTEGRÁCIÓS TESZT FUTTATÁSA

Az integrációs teszt futtatásához először bizonyosodjunk meg róla, hogy fut az anonimizáló eszköz. Ezután a korábbiakhoz hasonlóan navigáljunk el PowerShell segítségével a mellékletben található következő mappába: AnonymizationServer\tests\IntegrationTest. Ebben a mappában a következő parancsot adjuk ki.

```bash
docker-compose up --build
```

Ennek a parancsnak a hatására elindulnak a Python nyelven íródott integrációs tesztek, és kiírják a konzolba az eredményeiket. A teszthez használt konténer eltávolítására a következő parancsot kell kiadnunk.

```bash
docker-compose down
```
