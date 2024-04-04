# goTSVParser

# 🔨 Description 

Service for scanning directories with files in .tsv format and rewriting information into pdf/svg files

# 🧩 Config

```json
{
  "host": "0.0.0.0:8081",
  "tls": false,
  "certificate": "",
  "private": "",
  "dir_from": "from",
  "dir_to": "to",
  "dsn": "postgres://user:password@db:5432/dbname?sslmode=disable",
  "refresh_interval": 10,
  "svg_gen": false
}

```

# 🏴‍☠️ Flags
```
a - ip for REST -a=host
f - directory from -f="/User/..."
t - drectory to -t="/User/..."
d - connection string -d=connection string
r - refresh interval of scan directory -r=1
cert - path to certificate -cert=path_to_certificate
key - path to private key -key=path_to_key
tls - enable or disable tls certificate -tls=false/true
svg - enable generation .svg files -svg=true/false

```

# 📞 Request and Response

Request

```http

POST https://localhost:8080/api/all HTTP/1.1
Content-Type: application/json
{
    "unit_guid": "01749246-95f6-57db-b7c3-2ae0e8be671f",
    "limit": 2,
    "page": 2
}

```

Response

```json
{
    [
        {
            "Number": "3",
            "MQTT": "",
            "InventoryID": "G-044322",
            "UnitGUID": "01749246-95f6-57db-b7c3-2ae0e8be671f",
            "MessageID": "cold7_ComprSK_status",
            "MessageText": "Компрессор",
            "Context": "",
            "MessageClass": "working",
            "Level": "100",
            "Area": "LOCAL",
            "Address": "cold7_status.ComprSK_status",
            "Block": "",
            "Type": "",
            "Bit": "",
            "InvertBit": ""
        },
        {
            "Number": "4",
            "MQTT": "",
            "InventoryID": "G-044322",
            "UnitGUID": "01749246-95f6-57db-b7c3-2ae0e8be671f",
            "MessageID": "cold7_Temp_Al_HH",
            "MessageText": "Высокая температура",
            "Context": "",
            "MessageClass": "alarm",
            "Level": "100",
            "Area": "LOCAL",
            "Address": "cold7_status.Temp_Al_HH",
            "Block": "",
            "Type": "",
            "Bit": "",
            "InvertBit": ""
        }
    ],
    [
        {
            "Number": "16",
            "MQTT": "",
            "InventoryID": "G-044322",
            "UnitGUID": "01749246-95f6-57db-b7c3-2ae0e8be671f",
            "MessageID": "test_alarm",
            "MessageText": "Тест Аларм",
            "Context": "",
            "MessageClass": "alarm",
            "Level": "100",
            "Area": "LOCAL",
            "Address": "TestingForMsg.Alarm",
            "Block": "",
            "Type": "",
            "Bit": "",
            "InvertBit": ""
        },
    ],
}
