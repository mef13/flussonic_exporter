#!/bin/sh

cp_mod() {
if [ $# -eq 4 ]
then
  if ! [ -f "$3" ]; then
    echo "cp -n $2 $3"
    cp -n "$2" "$3"
    chown "$4":"$4" "$3"
  fi
else
  echo "cp_mod: invalid parameters count."
fi
}

copy_app() {
if [ $# -eq 4 ]
then
  cp -f "$2" "$3"
  chown "$4":"$4" "$3"
  chmod +x "$3"
else
  echo "copy_app: invalid parameters count."
fi

}

prepare_folder() {
if [ $# -eq 3 ];
then
  if ! [ -d "$2" ]; then
    echo "mkdir -p $2"
    mkdir "$2"
    chown -R "$3":"$3" "$2"
  fi
else
  echo "prepare_folder: invalid parameters count."
fi
}

add_user() {
if [ $# -eq 2 ];
then
  grep "$2:" /etc/passwd >/dev/null
  if [ $? -ne 0 ]; then
    adduser --system --no-create-home --home /dev/null --group "$2"
  else
    echo "User $2 already exist!"
  fi
else
  echo "add_user: invalid parameters count."
fi
}


if [ -n "$1" ]
then
  # shellcheck disable=SC2086
  case "$1" in
    add_user) add_user $* ;;
    cp_mod) cp_mod $* ;;
    prepare_folder) prepare_folder $* ;;
    copy_app) copy_app $* ;;
#    *) echo "$1 is not an option" ;;
  esac
else
  echo "No parameters found."
fi
