#!/bin/bash
cd "$(dirname "$0")"
dotnet publish -r linux-x64 -o output LittleRosie.csproj
if [ ! -e "csharp" ]; then
    ln -s output/LittleRosie csharp
fi
