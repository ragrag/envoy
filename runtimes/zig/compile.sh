#!/bin/bash
/usr/local/zig/zig build-exe -Doptimize=ReleaseFast --color off --cache-dir ${BOX_ROOT}/zig/cache --global-cache-dir ${BOX_ROOT}/zig/global-cache --name main main.zig
