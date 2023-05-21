#!/bin/bash
export PATH="/usr/local/kotlin/bin:/usr/local/openjdk/bin/:${PATH}"
/usr/local/kotlin/bin/kotlinc main.kt -include-runtime -d main.jar