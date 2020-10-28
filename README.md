# Page Extractor

## RUN

go run src/main.go

## TEST

go test ./...

## Request

```json
{
    "url": "http://example.com/",
    "template": {
        "selector": "div",
        "listName": "links",
        "method": "table",
        "listData": [
            {
                "selector": "p",
                "method": "text",
                "listName": "text"
            },{
                "selector": "a",
                "method": "link",
                "listName": "link"
            }
        ]
    },
    "html": "\n    \n<!doctype html>\n<html>\n<head>\n    <title>Example Domain</title>\n\n    <meta charset=\"utf-8\" />\n    <meta http-equiv=\"Content-type\" content=\"text/html; charset=utf-8\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n    <style type=\"text/css\">\n    body {\n        background-color: #f0f0f2;\n        margin: 0;\n        padding: 0;\n        font-family: -apple-system, system-ui, BlinkMacSystemFont, \"Segoe UI\", \"Open Sans\", \"Helvetica Neue\", Helvetica, Arial, sans-serif;\n        \n    }\n    div {\n        width: 600px;\n        margin: 5em auto;\n        padding: 2em;\n        background-color: #fdfdff;\n        border-radius: 0.5em;\n        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);\n    }\n    a:link, a:visited {\n        color: #38488f;\n        text-decoration: none;\n    }\n    @media (max-width: 700px) {\n        div {\n            margin: 0 auto;\n            width: auto;\n        }\n    }\n    </style>    \n</head>\n\n<body>\n<div>\n    <h1>Example Domain</h1>\n    <p>This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission.</p>\n    <p><a href=\"https://www.iana.org/domains/example\">More information...</a></p>\n</div>\n</body>\n</html>\n\n    "
}
```

> Template: rule to extract html content

response: 

```json
[
    {
        "link": [
            "https://www.iana.org/domains/example"
        ],
        "text": "This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission.\nMore information..."
    }
]
```