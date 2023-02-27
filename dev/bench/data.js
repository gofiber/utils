window.BENCHMARK_DATA = {
  "lastUpdate": 1677482392400,
  "repoUrl": "https://github.com/gofiber/utils",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "committer": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "distinct": true,
          "id": "a1bdcfa4e17bf96395bd81a7a2762468565f587e",
          "message": ":construction_worker: make go 1.18 minimum version",
          "timestamp": "2022-09-23T20:32:09+03:00",
          "tree_id": "87e1c614d995b44a9ced008013198f861a0bcd04",
          "url": "https://github.com/gofiber/utils/commit/a1bdcfa4e17bf96395bd81a7a2762468565f587e"
        },
        "date": 1663956296257,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 64.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18826969 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 196.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "6219244 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 65.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18660182 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 174.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "6485299 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 83.91,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14519265 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 412.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2940142 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 88.79,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14068340 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 333.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3449799 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 82.02,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "13855543 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 601.1,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2040511 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 42.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "28240450 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.9471,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 32.39,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "37238107 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.997,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "296394441 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 38.92,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "29326492 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 88.56,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "12823032 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 93.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12805435 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 207.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5866706 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 88.82,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "13513023 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 9.647,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "124660440 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.423,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "857034735 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 5.359,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "221690346 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 115.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9373128 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 345.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3548768 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 123.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9912066 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 332.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3462897 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "committer": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "distinct": true,
          "id": "a1bdcfa4e17bf96395bd81a7a2762468565f587e",
          "message": ":construction_worker: make go 1.18 minimum version",
          "timestamp": "2022-09-23T20:32:09+03:00",
          "tree_id": "87e1c614d995b44a9ced008013198f861a0bcd04",
          "url": "https://github.com/gofiber/utils/commit/a1bdcfa4e17bf96395bd81a7a2762468565f587e"
        },
        "date": 1663956387661,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 57.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19912767 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 150.8,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7967053 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 72.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16262934 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 171.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "6933460 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 72.12,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16129975 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 360,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3408698 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 74.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16036868 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 295.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4004860 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 75.85,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "15386491 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 546,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2238800 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 36.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "32811636 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.9167,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 28.88,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "38154003 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.962,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "305740502 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 34.62,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "32778691 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 81.6,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "14391530 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 102.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11608651 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 180.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6496843 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 80.62,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "14569904 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 8.871,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "134564601 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.396,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "862069249 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.342,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "275011036 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 97.03,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10403553 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 298.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3901646 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 109.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10346767 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 301.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3945277 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "committer": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "distinct": true,
          "id": "a1bdcfa4e17bf96395bd81a7a2762468565f587e",
          "message": ":construction_worker: make go 1.18 minimum version",
          "timestamp": "2022-09-23T20:32:09+03:00",
          "tree_id": "87e1c614d995b44a9ced008013198f861a0bcd04",
          "url": "https://github.com/gofiber/utils/commit/a1bdcfa4e17bf96395bd81a7a2762468565f587e"
        },
        "date": 1663956666637,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 47.91,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "25042434 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 127.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9395613 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 61.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19405305 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 141.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8287904 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 61.91,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19308669 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4079541 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18774159 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 250.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4774011 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 64.73,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17940790 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 452.8,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2583669 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38581777 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7543,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 24.85,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "45188552 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.353,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "357956008 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 30.23,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "36624456 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 69.35,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16918502 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 82.45,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14548659 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 152.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7826206 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 67.38,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17534570 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.504,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "159876219 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.174,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.686,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325462695 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 82.76,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14071214 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 253.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4703698 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 93.24,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12682668 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 252.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4727508 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "35188480+wangjq4214@users.noreply.github.com",
            "name": "Jinquan Wang",
            "username": "wangjq4214"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "43258f91674e9fc80f413beef0e6f18d577876fe",
          "message": "Definitions required for client (#29)",
          "timestamp": "2022-09-25T14:49:26+03:00",
          "tree_id": "6f0812bac1641e20cfaa18cde12bdf684b6e195b",
          "url": "https://github.com/gofiber/utils/commit/43258f91674e9fc80f413beef0e6f18d577876fe"
        },
        "date": 1664106623347,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 47.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "25007773 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 127.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9381171 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 61.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19521658 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 140.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8469594 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 61.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19407770 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4067007 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.65,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18779151 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4771422 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 63.93,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "18147613 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 452.3,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2642883 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.55,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "37328883 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.755,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.11,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "48711052 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.351,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "358409592 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.88,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "35480468 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 68.97,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "17098424 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 81.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14722653 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 155,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7781192 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 66.23,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17811570 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.499,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "159866538 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.172,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.688,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325152964 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 82.41,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "13985163 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 255.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4629418 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 93.79,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12825763 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 254.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4740362 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "35188480+wangjq4214@users.noreply.github.com",
            "name": "Jinquan Wang",
            "username": "wangjq4214"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "43258f91674e9fc80f413beef0e6f18d577876fe",
          "message": "Definitions required for client (#29)",
          "timestamp": "2022-09-25T14:49:26+03:00",
          "tree_id": "6f0812bac1641e20cfaa18cde12bdf684b6e195b",
          "url": "https://github.com/gofiber/utils/commit/43258f91674e9fc80f413beef0e6f18d577876fe"
        },
        "date": 1664107837257,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 58.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21089148 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 172,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "6853371 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 54.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20702576 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 138.8,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8328055 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 63.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17295894 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 325.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3627494 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 65.35,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18522790 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 276,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4488397 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 63.53,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "18562561 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 479.7,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2563809 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 32.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "33644683 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7315,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 24.92,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "43484426 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.219,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "366826928 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 30.18,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "34573706 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 79.27,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16935014 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 75.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15146359 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 177.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6912730 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 67.06,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17238037 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.505,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "155460492 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.134,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.269,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "272711414 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 102.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12385129 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 287.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4043355 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 103,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "11816773 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 292.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4042467 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a20485f7af36d2bd83e3f76a7859e7a40ca37808",
          "message": "Bump dependabot/fetch-metadata from 1.3.3 to 1.3.4 (#31)\n\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-10-04T23:28:05+03:00",
          "tree_id": "240d8f3559a75e2c91166d5b1bf5d81d87fa5cde",
          "url": "https://github.com/gofiber/utils/commit/a20485f7af36d2bd83e3f76a7859e7a40ca37808"
        },
        "date": 1664915344096,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 47.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "24886406 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 128.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9383402 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 61.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19532953 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 140.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8436039 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 61.82,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19405966 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4072008 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18827756 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4760814 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 64.42,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17630629 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 452.9,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2646109 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38198331 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7543,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 24.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "43823739 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.362,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "358313522 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.65,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "37312146 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 69.18,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "17023610 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 79.53,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15019935 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 156.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7696616 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 66.38,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17392597 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.499,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "160061383 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.173,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.682,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325943113 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 81.73,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14278410 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 253.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4716909 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 93.38,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12606572 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 250.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4810902 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "6c75e140ddc2b665f677d81c9b3517934211955c",
          "message": "Update test.yml",
          "timestamp": "2022-10-10T08:44:43+02:00",
          "tree_id": "a84102e8795f0493868fc501536b5f577b586c54",
          "url": "https://github.com/gofiber/utils/commit/6c75e140ddc2b665f677d81c9b3517934211955c"
        },
        "date": 1665384338861,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 47.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "24766539 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 128.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9210948 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 61.53,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19528192 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 142.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8277612 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 61.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19388976 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4069314 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18424070 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4744993 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 66.66,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17425996 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 453.5,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2619501 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.53,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38375368 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7539,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "44183170 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.349,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "358135180 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.85,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "36269376 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 69.41,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16683946 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 82.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14685546 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 155.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7748407 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 66.57,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17107054 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.509,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "159600442 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.172,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.684,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325794914 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 83.75,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14359902 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 256.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4676439 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 93.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12632416 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 252.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4711066 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "committer": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "distinct": true,
          "id": "af87ea402459ceb1ca2efd5b8268104ddbb13592",
          "message": ":sparkles: http: add new mime types",
          "timestamp": "2022-10-17T22:28:14+03:00",
          "tree_id": "a78233955701a48d7a7266e0dbff614d0cdca47a",
          "url": "https://github.com/gofiber/utils/commit/af87ea402459ceb1ca2efd5b8268104ddbb13592"
        },
        "date": 1666034978948,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 41.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29040070 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 137.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8537031 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 65.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18196798 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 154.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7723588 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 67.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17614718 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 298.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4025084 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 66.79,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16024676 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 240.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4956345 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 61.37,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "19236889 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 437.6,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2717865 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38808292 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.8034,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 23.82,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "46706137 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 2.821,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "425799676 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 28.97,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "39000631 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 68.74,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16354521 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 74.98,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16000342 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 151.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7915584 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 64.71,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18220290 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 8.096,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "148981582 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.206,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "995559628 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "298645344 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 79.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14731138 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 317.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3762416 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 100.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "11707754 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 269,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4413063 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "committer": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "distinct": true,
          "id": "2691728230fb4ea9ae3dbdd84eb119bf9089e476",
          "message": ":sparkles: http: add new mime types (https://github.com/gofiber/fiber/commit/902e30efb676dbc1012946701db96b170b4bb741)",
          "timestamp": "2022-10-17T22:30:15+03:00",
          "tree_id": "a78233955701a48d7a7266e0dbff614d0cdca47a",
          "url": "https://github.com/gofiber/utils/commit/2691728230fb4ea9ae3dbdd84eb119bf9089e476"
        },
        "date": 1666035087232,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 47.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "24968548 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 132.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9248793 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 61.46,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19540285 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 142.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8435652 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 61.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19391215 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 289.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4078858 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18584845 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 261.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4585856 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 64.18,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "18159046 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 454.9,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2605665 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38434182 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7711,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.72,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "44636530 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.349,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "357425814 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.67,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "37056469 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 70.78,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "17022392 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 82.11,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14563779 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 156.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7616320 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 66.73,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17294870 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.594,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "157872643 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.172,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.691,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325805414 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 82.41,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14280999 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 254.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4707770 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 92.19,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12386997 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 251.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4731067 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "committer": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "distinct": true,
          "id": "77200dae849839af700047d7ef712379ee656628",
          "message": ":sparkles: http: update mime type of js (https://github.com/gofiber/fiber/commit/9a6002056c4c1e76a44f0c68a0d7c27c39e09507\\#diff-7f1417e0dd9749e9228a9aa79d64d8080e662f106d9c491bd8a669a50dd3940a)",
          "timestamp": "2022-10-17T22:32:03+03:00",
          "tree_id": "778db26f9ed234520b24329272fe9ea8dc427888",
          "url": "https://github.com/gofiber/utils/commit/77200dae849839af700047d7ef712379ee656628"
        },
        "date": 1666035192865,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 56.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21787359 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 156.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7117555 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 71.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16776342 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 172.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7014139 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 72.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16247488 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 341.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3445705 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 74.32,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14270066 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 292.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3941223 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 79.32,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "15290523 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 539.1,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2219472 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 36.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "31911806 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.8755,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 29.39,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "37365487 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 4.033,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "300541257 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 34.43,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "31246909 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 80.26,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "14665153 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 93.09,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12695732 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 176.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6605978 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 79.04,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "14519412 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 8.871,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "136970218 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.374,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "886959278 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.348,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "266567185 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 97.65,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12141001 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 292,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4096370 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 107.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10950343 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 294.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4036696 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jozsef@sallai.me",
            "name": "József Sallai",
            "username": "jozsefsallai"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8898e89353ce22390876d3ace387da4d5c6e2f23",
          "message": "Merge pull request #32 from sadfun/master\n\n🔥 Add IsIPv4/v6 util",
          "timestamp": "2022-10-17T23:55:06+03:00",
          "tree_id": "00d1681ab169a5244389f4760c6ed8ff83d459eb",
          "url": "https://github.com/gofiber/utils/commit/8898e89353ce22390876d3ace387da4d5c6e2f23"
        },
        "date": 1666040184538,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 56.54,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21719092 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 156.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7890182 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 54.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21762139 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 146.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8776209 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 73.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16541476 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 369.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3277621 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 73.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16394312 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 295.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4050321 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 69.58,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "16152258 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 523,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2408091 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 37.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30855019 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.8123,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 27.44,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "37536321 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.778,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "301270243 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 34.74,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "34378532 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 83.36,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "15789028 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 92.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11857959 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 174.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6812090 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 72.03,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "16328479 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 9.014,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "136028612 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.193,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "984463508 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 5.223,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "223650475 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 32.24,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "36224406 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 70.11,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "16142113 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 100.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11893856 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 163.5,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7562198 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 101.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12023061 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 305.8,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3908156 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 101.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10984609 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 301,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3990598 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "distinct": true,
          "id": "d278c2e2f632222abd82487bed46e5e1a1f9d774",
          "message": "Add timestamp calculation functionality",
          "timestamp": "2022-10-18T19:17:13+02:00",
          "tree_id": "b0144be6c4e1d50af01c4a633444dca32de04de0",
          "url": "https://github.com/gofiber/utils/commit/d278c2e2f632222abd82487bed46e5e1a1f9d774"
        },
        "date": 1666113526597,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 48.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "25066243 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 126.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9441936 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 61.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19554693 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 140.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8515756 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 61.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19463886 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4080787 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18474435 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4659046 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 64.11,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17557502 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 463.8,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2650884 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 32.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "37449004 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7543,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 24.72,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "42758655 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.352,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "358111832 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.93,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "36927433 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 68.45,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16966338 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 80.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14993166 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 153.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7791758 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 66.58,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17446010 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.817,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "153361610 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.172,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.688,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "321825996 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.28,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "47059917 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 60.07,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "19742810 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 83.48,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14359459 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 141.1,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8387355 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 91.19,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12385647 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 251.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4666575 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 81.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14455098 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 248.8,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4794387 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3718,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.61,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22229545 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "distinct": true,
          "id": "7d092717a38862673e01238def4a17139dd49d8d",
          "message": "Add timestamp calculation functionality",
          "timestamp": "2022-10-19T09:15:29+02:00",
          "tree_id": "8d2ef6e491d813343284373bc233c4f317c43aca",
          "url": "https://github.com/gofiber/utils/commit/7d092717a38862673e01238def4a17139dd49d8d"
        },
        "date": 1666163805126,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 61.42,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19545438 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 129.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9362248 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 51.65,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23199194 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 141.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8368753 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 64.09,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18722818 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 301.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4067887 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18951045 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4773140 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 65.86,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17704990 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 458.7,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2590339 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.82,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "37838400 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.754,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.33,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "41278891 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.35,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "358256898 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 30.47,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "36604554 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 68.61,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "17120554 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 77.99,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15433665 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 150.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7973976 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 67.29,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17429557 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.807,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "154237762 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.174,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.689,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "324189237 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "46209289 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 60.94,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "19341267 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 83.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14290950 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 142.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8327017 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 80.34,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14820088 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 251.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4766175 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 91.61,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "13042135 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 252.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4749548 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3714,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.61,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22339101 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "distinct": true,
          "id": "2de34c77748c8c25ad37740e4279a8104694493d",
          "message": "Add timestamp calculation functionality",
          "timestamp": "2022-10-19T09:20:26+02:00",
          "tree_id": "3bc7756cb51e79106a27b5bd930696bc9e6cba25",
          "url": "https://github.com/gofiber/utils/commit/2de34c77748c8c25ad37740e4279a8104694493d"
        },
        "date": 1666164094580,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 65.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18214611 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 136.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8724030 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 41.17,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29162883 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 154,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7783688 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 67.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17666128 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 309.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3846921 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 67.54,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17971039 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 238.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4983727 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 61.21,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "19011829 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 454.6,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2610448 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.08,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38825023 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.8033,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 23.91,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "46746928 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 2.657,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "451696365 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.22,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "39385449 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 70.39,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "14474710 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 72.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16352301 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 152.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7946035 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 62.05,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18796603 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 8.105,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "146366947 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.205,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "993996640 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.035,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "298617561 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 29.84,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "40086590 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 65.72,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17895824 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 98.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12096916 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 161.8,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7627456 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 94.32,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12263450 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 321.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3733791 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 99.69,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "11777856 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 272.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4354630 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.4224,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22419747 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b802ab8fe427d589a83ccc9407e66dcea36d8259",
          "message": "Merge pull request #33 from gofiber/add_timestamp_calculator\n\nAdd timestamp calculation functionality",
          "timestamp": "2022-10-19T09:25:58+02:00",
          "tree_id": "3bc7756cb51e79106a27b5bd930696bc9e6cba25",
          "url": "https://github.com/gofiber/utils/commit/b802ab8fe427d589a83ccc9407e66dcea36d8259"
        },
        "date": 1666164424632,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 61.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19478707 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 133.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9294303 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 51.64,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23269221 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 143.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8165335 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 64.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18600465 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4075800 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.19,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18966458 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4768768 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 65.4,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17767620 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 454.1,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2549032 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "37351856 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7539,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.06,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "43223984 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.35,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "357460388 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.27,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "38070214 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 70.29,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16803764 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 80.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14899438 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 150.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7976533 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 64.92,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17578916 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.723,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "155356402 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.173,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.682,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325925680 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "47734399 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 62.24,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18502516 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 83.69,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14353524 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 134.3,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8824191 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 82.17,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14442058 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 249.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4438894 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 91.44,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "11619141 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 248.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4808232 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3714,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22426952 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b802ab8fe427d589a83ccc9407e66dcea36d8259",
          "message": "Merge pull request #33 from gofiber/add_timestamp_calculator\n\nAdd timestamp calculation functionality",
          "timestamp": "2022-10-19T09:25:58+02:00",
          "tree_id": "3bc7756cb51e79106a27b5bd930696bc9e6cba25",
          "url": "https://github.com/gofiber/utils/commit/b802ab8fe427d589a83ccc9407e66dcea36d8259"
        },
        "date": 1666165521205,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 61.38,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19492411 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 130.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9229304 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 51.55,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22995708 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 143.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8228671 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 65.85,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18634262 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4075845 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.24,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18755774 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4757365 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 64.34,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17736722 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 461.9,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2638506 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.63,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38412054 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7558,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.11,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "43469856 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.356,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "358193463 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.39,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "38163207 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 72.05,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16701174 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 76.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15616315 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 152,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7927736 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 65.35,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17977240 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.731,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "155205000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.173,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.688,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "324791290 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "46706264 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 63.24,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18377252 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 83.72,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14279312 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 134.5,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8869165 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 81.56,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14297044 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 253,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4772516 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 94.44,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12649466 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 251.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4712584 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3719,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22066156 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "632993922988c11a95f7db52c32a5cb436b7ee97",
          "message": "Merge pull request #34 from gofiber/dependabot/go_modules/github.com/stretchr/testify-1.8.1\n\nBump github.com/stretchr/testify from 1.8.0 to 1.8.1",
          "timestamp": "2022-10-24T09:40:24+02:00",
          "tree_id": "4045d6bf6cba7556f878e759677137cfb672eeae",
          "url": "https://github.com/gofiber/utils/commit/632993922988c11a95f7db52c32a5cb436b7ee97"
        },
        "date": 1666597352806,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 53.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21874444 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 162.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7749369 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 53.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22221134 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 141.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8170270 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 69.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15659692 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 342.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3515427 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 68.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15951538 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 280.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4131756 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 66.69,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17956570 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 499.9,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2398054 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 35.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "33791503 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.8098,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.27,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "44946549 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.439,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "345880723 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 34.25,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "33073888 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 79.92,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "15074744 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 91.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13177741 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 173.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6693291 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 74.75,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "16243348 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 8.474,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "138373214 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.243,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.538,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "273830690 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 30.64,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "40000389 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 68.95,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17126930 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 98.08,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12697796 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 156.8,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7474054 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 98.78,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "11777404 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 272.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4149072 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 98.23,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12825297 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 279.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4397490 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3483,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 49.46,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22941812 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3c95ea8646cd693e31bdcd7ca70eb68a279a29cc",
          "message": "Merge pull request #35 from gofiber/dependabot/github_actions/lewagon/wait-on-check-action-1.2.0\n\nBump lewagon/wait-on-check-action from 1.1.2 to 1.2.0",
          "timestamp": "2022-10-24T09:42:17+02:00",
          "tree_id": "424079cba91307fbbdb895711cd9cfc82353f598",
          "url": "https://github.com/gofiber/utils/commit/3c95ea8646cd693e31bdcd7ca70eb68a279a29cc"
        },
        "date": 1666597429392,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 61.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19550871 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 130.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8705666 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 51.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23300428 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 145.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7989854 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 66.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18697094 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 305.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4078263 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18802120 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4574404 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 63.92,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17944744 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 455.3,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2653316 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "37373913 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7533,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.58,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "44754748 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.451,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "351573258 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.84,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "35709390 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 72.66,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16808637 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 83.07,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14593729 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 153.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7818099 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 64.83,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17511286 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.726,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "155418644 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.172,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.799,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325853245 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "45441260 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 62.7,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18534655 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 83.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14324524 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 138,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8903871 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 81.74,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14157646 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 249.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4785480 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 91.31,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "13109706 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 248.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4824836 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3717,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22446872 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "distinct": true,
          "id": "73837ff49846715248acc83070750b4e46c71199",
          "message": "Changes for version 2 with go1.19",
          "timestamp": "2022-10-30T10:37:52+01:00",
          "tree_id": "52c6379bef972c312d09541b25bb9a4b318c9da3",
          "url": "https://github.com/gofiber/utils/commit/73837ff49846715248acc83070750b4e46c71199"
        },
        "date": 1667122737110,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 65.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18171706 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 137.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8688692 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 41.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29135271 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 154.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7621786 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 67.99,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17612899 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 308.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3787341 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 66.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17975809 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 237.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5019740 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 60.97,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "16583629 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 454,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2631126 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 30.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38116504 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.804,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 23.99,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "46017352 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 2.822,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "424848972 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 28.78,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "38724512 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 70.66,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16761806 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 81.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14714205 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 150.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7949324 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 62.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18808114 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 8.123,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "146466169 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.21,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "992020771 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "298296554 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 30.53,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "40339088 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 67.32,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17744281 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 99.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12106410 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 137,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8664445 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 78.94,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "15165073 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 318.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3756590 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 101.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "11722146 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 275.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4341207 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.4259,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.37,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22430427 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "rene@gofiber.io",
            "name": "René Werner",
            "username": "ReneWerner87"
          },
          "distinct": true,
          "id": "73837ff49846715248acc83070750b4e46c71199",
          "message": "Changes for version 2 with go1.19",
          "timestamp": "2022-10-30T10:37:52+01:00",
          "tree_id": "52c6379bef972c312d09541b25bb9a4b318c9da3",
          "url": "https://github.com/gofiber/utils/commit/73837ff49846715248acc83070750b4e46c71199"
        },
        "date": 1667122938031,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 61.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19562844 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 129.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7864382 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 51.64,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23170467 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 146.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8138312 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 64.42,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18623910 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4073622 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19027116 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4757125 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 67.31,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17813767 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 458.7,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2624942 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "37357984 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7538,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.54,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "42849642 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.354,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "357308095 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.84,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "37450095 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 70.88,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16597884 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 81.07,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14783629 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 151.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7886475 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 65.96,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17946014 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.735,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "155344706 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.174,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.686,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325538328 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "47471610 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 63.64,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18562236 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 83.79,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14312791 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 135.9,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8532138 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 83.64,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "13953360 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 256.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4685205 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 93.26,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12261874 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 254.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4702608 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3727,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.51,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22227138 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "efectn@protonmail.com",
            "name": "M. Efe Çetin",
            "username": "efectn"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7112314b7d034e68b6b94246058f1ddf80bcbcf5",
          "message": "use vulncheck instead of snyk",
          "timestamp": "2022-11-03T19:21:04+03:00",
          "tree_id": "6cf0268257fce49ae4c2cefc1d7a7b1c032ee6f6",
          "url": "https://github.com/gofiber/utils/commit/7112314b7d034e68b6b94246058f1ddf80bcbcf5"
        },
        "date": 1667492556078,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 61.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19548217 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 130.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8956090 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 51.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23227052 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 146,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8205861 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 64.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18665598 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4070619 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 63.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18970458 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4753068 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 64.45,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17247169 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 453.8,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2639006 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38251837 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7537,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.81,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "44524551 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.348,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "358463458 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.45,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "37981702 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 70.6,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16665302 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 79.16,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15184327 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 152.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7877337 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 65.35,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17897349 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.726,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "155336198 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.172,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.685,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325793010 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "47516095 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 62.42,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18590782 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 83.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14440594 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 134.5,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8842462 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 81.42,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14258931 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 250.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4761502 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 91.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12652640 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 250,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4781068 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3717,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22144640 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "679b34224ecce0096188eb099f6c6ed890558746",
          "message": "Bump dependabot/fetch-metadata from 1.3.4 to 1.3.5 (#36)\n\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-11-05T10:36:15+03:00",
          "tree_id": "eaf9bf3173a94bd6180e2fa792d0dde89c3b533f",
          "url": "https://github.com/gofiber/utils/commit/679b34224ecce0096188eb099f6c6ed890558746"
        },
        "date": 1667633856179,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 63.05,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18466425 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 178.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "6488971 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 63.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18415894 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 170.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7119510 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 79.51,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14446412 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 407.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2969244 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 78.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14756608 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 340.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3583694 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 79.73,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "14458840 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 589.7,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2028319 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 41.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "24910654 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.8848,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 29.22,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "36266038 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.873,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "310283817 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 36.05,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "30998917 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 88.66,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "14301360 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 98.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12247036 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 199.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6059954 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 84.53,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "13123826 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 9.405,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "126320557 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.392,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "806147227 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 5.351,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "229981237 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 35.42,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "35194824 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 85.77,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "13730452 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 117.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9679767 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 183.3,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "6248436 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 116.7,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10380852 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 324.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3452540 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 120.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10731938 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 342,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3548901 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.4233,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 59.51,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21453727 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "committer": {
            "email": "efectn@protonmail.com",
            "name": "Muhammed Efe Çetin",
            "username": "efectn"
          },
          "distinct": true,
          "id": "b2d9db23783b771c86d083d37558f93f36af6f1a",
          "message": ":bug: fix function name tests",
          "timestamp": "2022-11-05T10:38:07+03:00",
          "tree_id": "b2ab6529599c63df5fb6b2958e146dc5d63ef010",
          "url": "https://github.com/gofiber/utils/commit/b2d9db23783b771c86d083d37558f93f36af6f1a"
        },
        "date": 1667633956992,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 51.04,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23436334 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 99.52,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "11809431 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 30.99,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38655231 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 114.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10400288 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 48.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "24424264 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 211.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5652176 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 56.41,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21262713 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 188.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6359840 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 43.47,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "27133717 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 322.4,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "3710523 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 25.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "46236529 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.6185,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 19.38,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "61542024 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.089,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "388295848 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 22.96,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "49903456 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 57.04,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "20682270 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 58.72,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20438276 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 126.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9402050 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 48.37,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "24475977 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.409,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "162195976 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 0.9368,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.706,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "316670739 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 20.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "58089704 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 45.66,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "25911318 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 61.53,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19713026 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 108.4,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "11026276 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 60.22,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "19620488 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 191.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "6270624 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 77.07,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "15362524 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 192,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "6255409 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3156,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 55.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21523470 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "git@leonklingele.de",
            "name": "leonklingele",
            "username": "leonklingele"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1dd7023ee3babde1da5823ce9d239c63d729bf62",
          "message": "http: update list of HTTP status codes (#37)",
          "timestamp": "2022-11-12T09:08:07+03:00",
          "tree_id": "803e62fbfe11d9c097dc694ba69dc607f39a1900",
          "url": "https://github.com/gofiber/utils/commit/1dd7023ee3babde1da5823ce9d239c63d729bf62"
        },
        "date": 1668233360970,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 65.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18185593 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 137,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8753215 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 41.17,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29110802 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 153.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7701460 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 67.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17659474 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 403.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2973284 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 66.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17979316 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 238.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5043302 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 61.7,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "19060908 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 454.4,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2628910 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 30.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38245622 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.8035,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 24,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "46406631 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 2.657,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "450955750 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 28.78,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "38890016 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 70.24,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16619151 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 80.09,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14993512 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 152.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7827298 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 62.77,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18684385 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.747,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "155717498 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.207,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "994876030 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.018,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "298607437 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 29.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "40492198 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 67.95,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17851210 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 98.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12119521 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 136.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8669163 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 77.11,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14981193 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 318.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3753704 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 100,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "11870547 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 274.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4345970 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.4261,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.43,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22422494 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e2e2338d82bcfefc579d04d4b2e3332305153bb3",
          "message": "Update dependabot-automerge.yml",
          "timestamp": "2022-11-14T11:48:12+01:00",
          "tree_id": "8bc3cf4decf1c244ba8163ddd8c8018774560bf5",
          "url": "https://github.com/gofiber/utils/commit/e2e2338d82bcfefc579d04d4b2e3332305153bb3"
        },
        "date": 1668422958970,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 73.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16269625 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 156.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7684561 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 61.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19349816 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 181.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "6850639 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 77.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15575300 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 353.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3395832 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 76.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15773995 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 305.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3973216 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 77.32,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "15017618 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 552.4,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2183652 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 38.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29785670 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.908,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 30.02,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "36516628 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 4.037,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "295681680 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 35.57,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "31671583 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 84.41,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "13843791 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 99.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12122481 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 182.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6627640 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 79.36,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "14485921 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 9.263,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "129521588 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.408,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "850333218 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.46,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "268551739 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 30.43,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38457351 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 76.53,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "15489002 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 101,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11974658 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 165.5,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7116351 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 99.33,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12144264 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 315.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3856501 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 111.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10560696 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 306.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3867328 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.4455,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 64.28,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18610632 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "5387204ebac7d6034c84b25dff32c96d99a4fba1",
          "message": "Merge pull request #38 from leonklingele/timestamp-helper-func\n\nv3: time: do not export timestamp var, instead provide helper function to retrieve it",
          "timestamp": "2022-11-14T16:11:37+01:00",
          "tree_id": "2977ee4a50d47583e1cf3c9719bf2ade96c714bd",
          "url": "https://github.com/gofiber/utils/commit/5387204ebac7d6034c84b25dff32c96d99a4fba1"
        },
        "date": 1668438770990,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 65.98,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18102415 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 137.8,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8748690 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 41.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29126733 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 153.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7768797 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 67.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17649802 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 314,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3946083 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 66.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17976522 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 237.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5038852 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 61.1,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "18867175 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 453.3,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2648899 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 30.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "39009142 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.8036,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 23.77,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "46073487 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 2.657,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "451499992 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 30.85,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "34023931 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 69.76,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16837560 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 85.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14028405 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 150.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8046676 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 62.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18647065 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 8.113,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "154108316 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.208,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "993753949 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.035,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "248621989 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 29.99,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "40155469 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 65.7,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17813580 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 99,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12127597 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 136.9,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8735954 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 78.28,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14865114 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 322.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3757095 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 99.89,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "11951054 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 280.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4317130 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.4222,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.09,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22650889 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cd7a17a0e7d323045162447a5826fd8bb29a5d80",
          "message": "Update dependabot-automerge.yml",
          "timestamp": "2022-11-17T09:09:19+01:00",
          "tree_id": "0750c6f1324f2869e7f7548261455f9915961fae",
          "url": "https://github.com/gofiber/utils/commit/cd7a17a0e7d323045162447a5826fd8bb29a5d80"
        },
        "date": 1668672686698,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 61.42,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19556220 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 134.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9188884 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 51.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23291481 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 145.8,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8231587 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 63.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18899944 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4075310 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 61.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19471960 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4775702 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 64.56,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17782711 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 452.6,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2624359 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 31.76,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "37789926 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7534,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.84,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "43967661 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.352,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "358151814 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 29.54,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "37657309 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 70.48,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16583379 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 83.35,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14397122 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 151.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7822125 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 65.54,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18041811 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.729,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "154996867 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.172,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.687,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "325872820 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.39,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "47559710 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 62.77,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18727279 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 83.64,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14379466 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 134.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8841866 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 82.84,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14340064 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 250.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4753902 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 91.38,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12766305 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 254.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4778779 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3727,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 53.39,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22485877 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9f786cfce70ab15115b538d4b170cc7929e62c6e",
          "message": "Merge pull request #39 from gofiber/dependabot/github_actions/dependabot/fetch-metadata-1.3.6\n\nBump dependabot/fetch-metadata from 1.3.5 to 1.3.6",
          "timestamp": "2023-01-24T07:11:57+01:00",
          "tree_id": "226e02b873df9f7065e50cbbd8a9d17dec54d756",
          "url": "https://github.com/gofiber/utils/commit/9f786cfce70ab15115b538d4b170cc7929e62c6e"
        },
        "date": 1674540791441,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 54.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22697654 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 151.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7818835 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 71.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17655537 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 168.9,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7354534 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 72.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17599887 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 345.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3409497 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 77.28,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15750238 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 295.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3813751 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 77.81,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "15193387 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 564.4,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2163542 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 38.37,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "33095157 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.8803,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 32.04,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "33551768 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.905,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "297033190 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 36.42,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "29514802 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 86.29,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "12543454 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 98.35,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11904990 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 182.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6843970 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 85.02,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "13861344 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 9.411,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "125756935 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.361,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "878233174 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 5.597,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "212007685 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 29.84,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "40322823 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 77.87,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "15680604 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 102.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11633373 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 166.8,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7091835 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 110.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10560337 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 291.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4045675 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 98.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12140442 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 311.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "3970441 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.4482,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 64.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18176281 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "6383794738c9eb2742f30bfd4d0e4a058324db34",
          "message": "Merge pull request #40 from gofiber/dependabot/github_actions/lewagon/wait-on-check-action-1.3.1\n\nBump lewagon/wait-on-check-action from 1.2.0 to 1.3.1",
          "timestamp": "2023-02-17T08:06:34+01:00",
          "tree_id": "672834f1b009b4e1982a29731b21939368005f41",
          "url": "https://github.com/gofiber/utils/commit/6383794738c9eb2742f30bfd4d0e4a058324db34"
        },
        "date": 1676617655709,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 47.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "24973864 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 128,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9328321 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 61.38,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19549867 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 142.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8351653 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 62.12,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19354922 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 294.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4076635 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 65.55,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18138361 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4770889 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 63.56,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "18161288 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 468,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2561460 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 32.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "36476643 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7548,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 25.75,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "43304480 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.352,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "358174255 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 30.25,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "37293661 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 70.83,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16658740 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 83.84,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14212420 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 153.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7742740 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 68.93,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "17037272 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.499,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "159985798 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.171,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 4.576,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "263441029 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.87,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "46387534 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 62.43,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18580598 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 87.69,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13719561 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 135.7,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8726176 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 92.71,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12292933 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 251.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4775282 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 81.77,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "14392396 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 253.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4746170 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3721,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 55.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21401775 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rene@gofiber.io",
            "name": "RW",
            "username": "ReneWerner87"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0980153c83f90f0cb25820bb415d4d0697785ea9",
          "message": "Merge pull request #41 from gofiber/dependabot/go_modules/github.com/stretchr/testify-1.8.2\n\nBump github.com/stretchr/testify from 1.8.1 to 1.8.2",
          "timestamp": "2023-02-27T08:18:47+01:00",
          "tree_id": "df0518e40be4d36be672d2989e1cbcd6e2b439f3",
          "url": "https://github.com/gofiber/utils/commit/0980153c83f90f0cb25820bb415d4d0697785ea9"
        },
        "date": 1677482391593,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_ToLowerBytes/fiber",
            "value": 61.38,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19503812 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLowerBytes/default",
            "value": 129.2,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9159158 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/fiber",
            "value": 51.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23291757 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpperBytes/default",
            "value": 148.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8166319 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/fiber",
            "value": 63.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18900606 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFoldBytes/default",
            "value": 289.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4077298 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/fiber",
            "value": 61.76,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19315544 times\n2 procs"
          },
          {
            "name": "Benchmark_EqualFold/default",
            "value": 251.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4774234 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/fiber",
            "value": 66.77,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "17343372 times\n2 procs"
          },
          {
            "name": "Benchmark_UUID/default",
            "value": 467.6,
            "unit": "ns/op\t     168 B/op\t       6 allocs/op",
            "extra": "2564511 times\n2 procs"
          },
          {
            "name": "Benchmark_ConvertToBytes/fiber",
            "value": 30.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38622450 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/unsafe",
            "value": 0.7602,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeString/default",
            "value": 26.68,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "39766300 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/unsafe",
            "value": 3.416,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "357359953 times\n2 procs"
          },
          {
            "name": "Benchmark_UnsafeBytes/default",
            "value": 31.07,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "36112593 times\n2 procs"
          },
          {
            "name": "Benchmark_ToString",
            "value": 71.89,
            "unit": "ns/op\t      40 B/op\t       2 allocs/op",
            "extra": "16189834 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/fiber",
            "value": 83.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14349763 times\n2 procs"
          },
          {
            "name": "Benchmark_GetMIME/default",
            "value": 153.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7797210 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/vendorContentType",
            "value": 70.51,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "16875416 times\n2 procs"
          },
          {
            "name": "Benchmark_ParseVendorSpecificContentType/defaultContentType",
            "value": 7.851,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "153038559 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/fiber",
            "value": 1.174,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_StatusMessage/default",
            "value": 3.692,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "324589120 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/fiber",
            "value": 25.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "46896949 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv4/default",
            "value": 63.92,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "18628087 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/fiber",
            "value": 83.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14283704 times\n2 procs"
          },
          {
            "name": "Benchmark_IsIPv6/default",
            "value": 135.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "8735234 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/fiber",
            "value": 84.66,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "13528677 times\n2 procs"
          },
          {
            "name": "Benchmark_ToUpper/default",
            "value": 265,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4458710 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/fiber",
            "value": 96.78,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "12166359 times\n2 procs"
          },
          {
            "name": "Benchmark_ToLower/default",
            "value": 256.8,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "4593610 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/fiber",
            "value": 0.3714,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_CalculateTimestamp/default",
            "value": 55.98,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21472296 times\n2 procs"
          }
        ]
      }
    ]
  }
}