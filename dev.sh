#!/bin/sh
parse_yaml() {
   local prefix=$2
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @ | tr @ '\034')
   sed -ne "s|^\($s\)\($w\)$s:$s\"\(.*\)\"$s\$|\1$fs\2$fs\3|p" \
      -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p" $1 |
      awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
      }
   }'
}

eval $(parse_yaml .gitlab-ci.yml)

echo "DOCKER_VERSION="$variables_DOCKER_VERSION >.env
echo "GOLANG_VERSION="$variables_GOLANG_VERSION >>.env
echo "NODE_VERSION="$variables_NODE_VERSION >>.env
echo "ALPINE_VERSION="$variables_ALPINE_VERSION >>.env
echo "RCLONE_VERSION="$variables_RCLONE_VERSION >>.env
echo "RESTIC_VERSION="$variables_RESTIC_VERSION >>.env
echo "DEBIAN_VERSION="$variables_DEBIAN_VERSION >>.env
echo "TELEPORT_VERSION="$variables_TELEPORT_VERSION >>.env
