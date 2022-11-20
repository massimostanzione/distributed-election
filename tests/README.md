# tests

This package contains node logs for test cases proposed into the project specification document.

Test cases are listed below:

    a) only the coordinator fails;
    b) a non-coordinator node fails;
    c) at least two nodes, which one of them is the coordinator, fails.

All tests are run:
- with $N = 4$ nodes;
- using the scripts that can be found into the [`examples`](https://github.com/massimostanzione/distributed-election/tree/main/examples) folder;
- applying the [`tests.ini`](https://github.com/massimostanzione/distributed-election/tree/main/configs/tests.ini) configuration file, overwritten only for the `NODE_PORT` parameter with the flag `-p`.

## Test case a)
See the `xx_CoordFail` folder.

After 20 seconds from the end of the first election, **node 4 - listening on port 40046 -** is shut down.

Log file contains the subsequent election log and the 20 seconds after that.

## Test case b)
See the `xx_NonCoordFail` folder.

After 20 seconds from the end of the first election, **node 1 - listening on port 40043 -** is shut down.

Log file contains the subsequent election log and the 20 seconds after that.

## Test case c)
See the `xx_MoreFails` folder.

After 20 seconds from the end of the first election, **nodes 3 and 4 - listening on ports 40045 and 40046 -** are shut down.

Log file contains the subsequent election log and the 20 seconds after that.
