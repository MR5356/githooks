name=$1
ssh_url=$2
commitId=$3

cd /tmp
git clone $ssh_url
cd $name
docker build -t $name:$commitId -f release/docker/Dockerfile .

cd -
rm -rf /tmp/$name