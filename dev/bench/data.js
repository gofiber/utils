window.BENCHMARK_DATA = {
  "lastUpdate": 1663956667130,
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
      }
    ]
  }
}