nettoyer () {
  # Suppression des processus de l'application app
  killall app-base 2> /dev/null

  # Suppression des processus de l'application ctl
  killall app-control 2> /dev/null

  # Suppression des processus tee et cat
  killall tee 2> /dev/null
  killall cat 2> /dev/null

  # Suppression des tubes nommés
  \rm -f /tmp/in* /tmp/out*
  exit 0
}

nettoyer