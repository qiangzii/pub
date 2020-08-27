#!/bin/bash

echo "Checking Folders ..."
ispDir="/var/isp"
ispCIDRDir="/var/isp/cidr"
ispRangeDir="/var/isp/range"
ispResultDir="/var/isp/results"
ChinaMobile="/var/isp/cidr/ChinaMobile"
ChinaUnicom="/var/isp/cidr/ChinaUnicom"
ChinaTelcom="/var/isp/cidr/ChinaTelcom"
Education="/var/isp/range/ChinaEducation"
CNISP="/var/isp/CN_ISP_RIB" #git项目文件夹

if [ ! -d "$ispDir" ]; then
    mkdir "$ispDir"
fi
if [ ! -d "$ispCIDRDir" ]; then
    mkdir "$ispCIDRDir"
fi
if [ ! -d "$ispResultDir" ]; then
    mkdir "$ispResultDir"
fi
if [ ! -d "$ispRangeDir" ]; then
    mkdir "$ispRangeDir"
fi
echo "Checking folders done"

cd "$ispDir"

echo "Get From https://github.com/yuhaiyang/CN_ISP_RIB/blob/master/StoneOS-User-Defined-ISP-split.DAT ..."
if [  -d "$CNISP" ]; then
    rm -rf "$CNISP"
fi
git clone https://github.com/yuhaiyang/CN_ISP_RIB.git 

cat ./CN_ISP_RIB/StoneOS-User-Defined-ISP-split.DAT | grep "^1:" | awk -F ":" '{print $2}' > "$ChinaMobile"
cat ./CN_ISP_RIB/StoneOS-User-Defined-ISP-split.DAT | grep "^2:" | awk -F ":" '{print $2}' > "$ChinaUnicom"
cat ./CN_ISP_RIB/StoneOS-User-Defined-ISP-split.DAT | grep -E '^3:|^4:'  | awk -F ":" '{print $2}' > "$ChinaTelcom"

echo "Get done"

echo "Get From Whois3 Client ..."
whois3 -h whois.apnic.net -l -i mb MAINT-CERNET-AP | grep 'inetnum' | awk -F ':' '{print $2}'| sed 's/^[ \t]*//g' > "$Education"
echo "Get done"

echo "All files have already prepared"