PUT _ingest/pipeline/msg-parser
{
    "description": "Preprocessing of gotten message string",
    "processors": [
        {
            "grok": {
                "pattern_definitions": {
                    "BOOL": "[0|1]"
                },
                "patterns": [
                    "^%{POSINT:channel_id},%{DATA:channel_name},%{GREEDYDATA:text},%{POSINT:time},%{DATA:user_name},%{POSINT:user_id},%{BOOL:turbo},%{BOOL:sub},%{BOOL:mod}$"
                ],
                "field": "msg"
            },
            "date": {
                "field": "time",
                "formats": [
                    "UNIX"
                ]
            },
            "remove": {
                "field": [
                    "time",
                    "msg"
                ]
            },
            "convert": {
                "field": "channel_id",
                "type": "integer"
            }
        },
        {
            "convert": {
                "field": "mod",
                "type": "integer"
            }
        },
        {
            "convert": {
                "field": "sub",
                "type": "integer"
            }
        },
        {
            "convert": {
                "field": "turbo",
                "type": "integer"
            }
        },
        {
            "convert": {
                "field": "user_id",
                "type": "integer"
            }
        }
    ]
}