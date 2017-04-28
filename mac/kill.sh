#!/bin/bash
# FlexiKiller
# Copyright (C) 2017 Claudio "nex" Guarnieri
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

# bash <(curl -s https://ops.securitywithoutborders.org/flexispy/kill.sh)

function check() {
	COUNTER=0
	# These are the folders used on mac by FlexiSpy.
	FLEXYS=("blblu" "blbld" "blblw")

	# Check if they all exist.
	for i in "${FLEXYS[@]}"
	do
		if [ -d /usr/libexec/.$i ]; then
			let COUNTER=COUNTER+1
		fi
	done

	# If all three of them exist we're confident there's a FlexiSpy
	# installation active.
	if [ $COUNTER -eq 3 ]; then
		return 0
	else
		return 1
	fi
}

function nuke() {
	echo ""

	# If user wants to uninstall we can use the script provided by
	# FlexiSpy itself :).
	UNINSTALL=/usr/libexec/.blblu/blblu/Contents/Resources/Uninstall.sh
	if [ -f $UNINSTALL ]; then
		sudo bash $UNINSTALL >/dev/null 2>/dev/null
	else
		echo "I can't find the uninstall script!"
	fi

	if ! check; then
		echo "The computer appears to be clean now :-)"
	else
		echo "The uninstall process failed! :-("
	fi
}

if check; then
	echo "This computer appears to be infected with FlexiSpy!"
	echo ""
	echo "Do you wish to uninstall it?"
	select yn in "Yes" "No"; do
		case $yn in
			Yes) nuke; break;;
			No) exit;;
		esac
	done
else
	echo "There doesn't seem to be FlexiSpy on this computer."
fi
