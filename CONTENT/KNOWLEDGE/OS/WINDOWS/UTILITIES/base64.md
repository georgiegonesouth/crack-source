# **BASE64**


## Encode a File

> ### Print
>
>     [Convert]::ToBase64String([System.IO.File]::ReadAllBytes("<SOURCE_FILE>"))

> ### Save to a File
>
>     [Convert]::ToBase64String([System.IO.File]::ReadAllBytes("<SOURCE_FILE>")) | Set-Content "<DESTINATION_FILE>" -NoNewline


## Encode a String

> ### Print
>
>     [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("<TEXT>"))

> ### Save to a File
>
>     [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("<TEXT>")) | Set-Content "<DESTINATION_FILE>" -NoNewline


## Decode a File

> ### Print
>
>     [System.Text.Encoding]::UTF8.GetString([Convert]::FromBase64String((Get-Content "<SOURCE_FILE>" -Raw)))

> ### Save to a File
>
>     [System.IO.File]::WriteAllBytes("<DESTINATION_FILE>", [Convert]::FromBase64String((Get-Content "<SOURCE_FILE>" -Raw)))


## Decode a String

> ### Print
>
>     [System.Text.Encoding]::UTF8.GetString([Convert]::FromBase64String("<TEXT>"))

> ### Save to a File
>
>     [System.IO.File]::WriteAllBytes("<DESTINATION_FILE>", [Convert]::FromBase64String("<TEXT>"))
