name=$1
ssh_url=$2
commitId=$3

cd /tmp
git clone $ssh_url
cd $name

BUILD_FILE=release/docker/build.sh
if [ -f "$BUILD_FILE" ]; then
  /bin/bash $BUILD_FILE $commitId
else
  docker build -t $name:$commitId -f release/docker/Dockerfile .
fi

docker rmi `docker images|grep none|awk '{print $3 }'|xargs`

cd -
rm -rf /tmp/$name