window.BENCHMARK_DATA = {
  "lastUpdate": 1663956297069,
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
      }
    ]
  }
}