# Kläb a Bäbber

## Vorbereitung

1) Installiere Mongodb Datenbank und starte diese :).
2) Download [kleber-bin](https://drive.google.com/open?id=0BxLCS9PB1fV2fkw3R1dCV3BBZkdOWktNVzFvNFZqVGhGWklkSWZ1Y2xEa0lWa3JDdTJ2OW8&authuser=0)
3) Download Public Dateien git clone https://github.com/rrawrriw/bebberPublic.git
4) Setze Umgebungsvariablen

```bash
# Alle Umgebungsvariablen müssen gesetzt werden
export BEBBER_IP=Bind IP für HTTP-Server
export BEBBER_PORT=Auf welchen Port soll HTTP-Server lauschen
export BEBBER_DB_SERVER=Mongodb Server IP oder Hostname
export BEBBER_DB_NAME=Mongodb Databasename
export BEBBER_PUBLIC=Pfad zu Public Ordner
```

## Ready to Kläb

```
bebber
```
