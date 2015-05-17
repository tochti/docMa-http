# Kläb a Bäbber

## Ablauf

1. Login Server
2. Installiere Mongodb Datenbank und starte diese.
3. Speicher [kleber-bin](https://drive.google.com/open?id=0BxLCS9PB1fV2fkw3R1dCV3BBZkdOWktNVzFvNFZqVGhGWklkSWZ1Y2xEa0lWa3JDdTJ2OW8&authuser=0).
    Nicht vergessen ```chmod u+x kleber```
4. Download Public Dateien 
    ```
    git clone https://github.com/rrawrriw/bebberPublic.git
    ```.
   Der Benutzer welcher zum Starten des HttpServers (kleber) verwendet wird muss auf das Verzeichnis Zugriff haben.
5. Setze Umgebungsvariablen
6. Starte kleber-bin

```bash
# Alle Umgebungsvariablen müssen gesetzt werden
export BEBBER_IP=Bind IP für HTTP-Server
export BEBBER_PORT=Auf welchen Port soll HTTP-Server lauschen
export BEBBER_DB_SERVER=Mongodb Server IP oder Hostname
export BEBBER_DB_NAME=Mongodb Databasename
export BEBBER_PUBLIC=Pfad zu Public Ordner
```

## Ready to Kläb

```bash
./kleber
```

## Kleber selber machen
1. Installiere go
2. Setze Umgebungsvariable GOPATH.
3. Installiere Packages
```bash
go get github.com/rrawrriw/kleber
go get github.com/rrawrriw/bebber
go get github.com/gin-gonic/gin
go get gopkg.in/mgo.v2
```
4.  Kompiliere und installiere 
```bash
go install github.com/rrawrriw/kleber
```
