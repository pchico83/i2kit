#!/bin/sh

set -e
echo "Building the docker linuxkit distribution..."
linuxkit build -format aws linuxkit.yml
echo "Done!"

for region in $AWS_REGIONS; do
	echo "Pushing AMI to $region..."
  export AWS_REGION=$region
  export AWS_DEFAULT_REGION=$region
	AMI=`linuxkit push aws -timeout 2400 -bucket i2kit -sriov simple linuxkit.raw 2>&1 | awk '{print $3;}'`
  aws ec2 modify-image-attribute --image-id $AMI --launch-permission "{\"Add\": [{\"Group\":\"all\"}]}"
	echo $region $AMI
	echo "Done!"
done
