#!/usr/bin/env bash
#
#  Copyright 2024 whipcode.app (AnnikaV9)
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on
#  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#  either express or implied. See the License for the specific
#  language governing permissions and limitations under the License.
#

if ! command -v argon2 &> /dev/null || ! command -v openssl &> /dev/null
then
    echo "This script requires argon2 and openssl. You have either one or both missing."
    exit 1
fi

SALT=$(openssl rand -hex 16)

echo -n $(</dev/stdin) | argon2 $SALT -id -t 1 -p 1 -l 32 -k 4096 -r > .masterkey
echo -n $SALT >> .masterkey
