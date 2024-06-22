package query

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshalToMinerData(t *testing.T) {
	raw := `
{
"data": [
    {
        "ID": 12,
        "CreatedAt": "2024-05-05T23:56:05.97843Z",
        "UpdatedAt": "2024-05-06T00:01:57.653755Z",
        "DeletedAt": null,
        "Miner": {
            "MacAddress": "AA:AA:AA:AA:AA:AA",
            "IPAddress": "22.22.22.22"
        },
        "Stats": {
            "HashRate": 141946.95,
            "RateIdeal": 141000,
            "Uptime": 130
        },
        "Config": {
            "Username": "username",
            "Password": "password",
            "Firmware": "firmware_version"
        },
        "MinerType": 0,
        "Mode": 0,
        "ModelName" : "Antminer S9",
        "Status": 0,
        "Pools": [
            {
                "ID": 34,
                "CreatedAt": "2024-05-05T23:56:05.978594Z",
                "UpdatedAt": "2024-05-06T00:01:57.652891Z",
                "DeletedAt": null,
                "Pool": {
                    "Url": "pool_url:3333",
                    "User": "pool_user",
                    "Pass": "pool_password",
                    "Status": "",
                    "Accepted": 0,
                    "Rejected": 0,
                    "Stale": 0
                },
                "MinerID": 12
            },
            {
                "ID": 35,
                "CreatedAt": "2024-05-05T23:56:05.978594Z",
                "UpdatedAt": "2024-05-06T00:01:57.653196Z",
                "DeletedAt": null,
                "Pool": {
                    "Url": "pool_url:3333",
                    "User": "pool_user",
                    "Pass": "pool_password",
                    "Status": "",
                    "Accepted": 0,
                    "Rejected": 0,
                    "Stale": 0
                },
                "MinerID": 12
            },
            {
                "ID": 36,
                "CreatedAt": "2024-05-05T23:56:05.978594Z",
                "UpdatedAt": "2024-05-06T00:01:57.653478Z",
                "DeletedAt": null,
                "Pool": {
                    "Url": "pool_url:3333",
                    "User": "pool_user",
                    "Pass": "pool_password",
                    "Status": "",
                    "Accepted": 0,
                    "Rejected": 0,
                    "Stale": 0
                },
                "MinerID": 12
            }
        ],
        "Fan": [
            5940,
            5880,
            5970,
            5820
        ],
        "Temperature": [
            43,
            43,
            58,
            58,
            43,
            43,
            56,
            56,
            42,
            42,
            55,
            55
        ],
        "Log": null,
        "FleetID": 1
    },
    {
        "ID": 10,
        "CreatedAt": "2024-05-05T23:56:05.976459Z",
        "UpdatedAt": "2024-05-06T00:01:57.71841Z",
        "DeletedAt": null,
        "Miner": {
            "MacAddress": "BB:BB:BB:BB:BB:BB:",
            "IPAddress": "22.3.3.3"
        },
        "Stats": {
            "HashRate": 144937.62,
            "RateIdeal": 141000,
            "Uptime": 132
        },
        "Config": {
            "Url": "pool_url:3333",
            "User": "pool_user",
            "Pass": "pool_password",
        },
        "MinerType": 0,
        "Mode": 0,
        "ModelName" : "Antminer S9",
        "Status": 0,
        "Pools": [
            {
                "ID": 28,
                "CreatedAt": "2024-05-05T23:56:05.976694Z",
                "UpdatedAt": "2024-05-06T00:01:57.717495Z",
                "DeletedAt": null,
                "Pool": {
                    "Url": "pool_url:3333",
                    "User": "pool_user",
                    "Pass": "pool_password",
                    "Status": "",
                    "Accepted": 0,
                    "Rejected": 0,
                    "Stale": 0
                },
                "MinerID": 10
            },
            {
                "ID": 29,
                "CreatedAt": "2024-05-05T23:56:05.976694Z",
                "UpdatedAt": "2024-05-06T00:01:57.717754Z",
                "DeletedAt": null,
                "Pool": {
                    "Url": "pool_url:3333",
                    "User": "pool_user",
                    "Pass": "pool_password",
                    "Status": "",
                    "Accepted": 0,
                    "Rejected": 0,
                    "Stale": 0
                },
                "MinerID": 10
            },
            {
                "ID": 30,
                "CreatedAt": "2024-05-05T23:56:05.976694Z",
                "UpdatedAt": "2024-05-06T00:01:57.718055Z",
                "DeletedAt": null,
                "Pool": {
                    "Url": "pool_url:3333",
                    "User": "pool_user",
                    "Pass": "pool_password",
                    "Status": "",
                    "Accepted": 0,
                    "Rejected": 0,
                    "Stale": 0
                },
                "MinerID": 10
            }
        ],
        "Fan": [
            5970,
            6000,
            6000,
            5970
        ],
        "Temperature": [
            46,
            46,
            63,
            63,
            46,
            46,
            61,
            61,
            43,
            43,
            58,
            58
        ],
        "Log": null,
        "FleetID": 1
     }
   ]
}`
	var resp QueryMinerResponse

	err := json.Unmarshal([]byte(raw), &resp)

	if err != nil {
		t.Error(fmt.Errorf("failed to unmarshal: %v", err))
	}
}
