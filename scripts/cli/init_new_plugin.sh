#! /bin/bash
# This script can fast init an demo plugin structs
# And add/change common things

default="\033[0m"
green="\033[32m"
red="\033[31m"
blue="\033[96m"

# macOs
shopt -s expand_aliases
if [[ $(uname) == "Darwin" ]]; then
  echo "darwin"
  alias sed='sed -i ""'
elif [[ $(uname) == "Linux" ]]; then
  echo "linux"
  alias sed='sed -i'
fi

# first you need tell the script the plugin name && plugin language
echo -e "${green}~~~~~~~~~~~~~~~~~~~~~~~~~~~ Welcome to Veinmind Tools ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
echo -e "pwd is $(pwd)"
if pwd | grep -q -E "veinmind-tools$"; then
  echo -e "${green}script is running in correct path"
else
  echo -e "${red}please make sure running this script at veinmind-tools project root path (use \`bash ./scripts/cli/add_a_new_plugin.sh\`)"
  exit 1
fi
echo -e "${default}Please input the plugins language you used${green}"
read -p "enter language(go/python): " language
echo -e "${default}Please input the plugins name ${green}"
read -p "enter name: " name
echo -e "${default}Is this plugins need to add in veinmind-runner? ${green}"
read -p "enter yes/no(default no): " publish

initCommon() {
  # mkdir
  mkdir $dir
  # cp demo
  cp ./example/parallel-container-run.sh ./example/README.md ./example/README.en.md $dir
  # init README
  echo -e "# ${name}  \n\n这是描述文件" >$dir/README.md
  echo -e "# ${name}  \n\nthis is description file" >$dir/README.en.md
  sed "s/veinmind-example/${name}/g" $dir/parallel-container-run.sh
}

initGoPlugin() {
  initCommon
  cp -r ./example/go/* $dir
  # init script
  sed "s/veinmind-example/${name}/g" $dir/Dockerfile
  sed "s/veinmind-example/${name}/g" $dir/script/build_amd64.sh
  sed "s/veinmind-example/${name}/g" $dir/script/build.sh
  sed "s/veinmind-example/${name}/g" $dir/go.mod
  sed "s/veinmind-example/${name}/g" $dir/cmd/cli.go
}

initPythonPlugin() {
  initCommon
  cp -r ./example/python/* $dir
  sed "s/veinmind-example/${name}/g" $dir/Dockerfile
  sed "s/veinmind-example/${name}/g" $dir/scan.py
}

#add2Runner() {
#  if [[ publish == "yes" ]]; then
#
#  fi
#}

name=veinmind-${name}
filename=$(echo $name | sed "s/-/_/g")
# validate
if [[ $language == "go" ]]; then
  dir="plugins/go/${name}"
  echo -e "${blue}init Veinmind GO Plugin ${name} at: $(pwd)/${dir}${red}"
  initGoPlugin
#  add2Runner
elif [[ $language == "python" ]]; then
  dir="plugins/python/${name}"
  echo -e "${blue}init Veinmind GO Plugin ${name} at: $(pwd)/${dir}${red}"
  initPythonPlugin
#  add2Runner
else
  echo -e "${red}Not Support Language"
fi
