# **BASE64**


## Encode a File

> ### Print
>
>     cat <SOURCE_FILE> | base64 -w 0; echo

> ### Save to a File
>
>     cat <SOURCE_FILE> | base64 > <DESTINATION_FILE>


## Encode a String

> ### Print
>
>     echo "<TEXT>" | base64 -w 0; echo

> ### Save to a File
>
>     echo "<TEXT>" | base64 > <DESTINATION_FILE>


## Decode a String

> ### Print
>
>     echo "<TEXT>" | base64 -d; echo

> ### Save to a File
>
>     echo "<TEXT>" | base64 -d > <DESTINATION_FILE>


## Decode a File

> ### Print
>
>     cat <SOURCE_FILE> | base64 -d; echo

> ### Save to a File
>
>     cat <SOURCE_FILE> | base64 -d > <DESTINATION_FILE>
