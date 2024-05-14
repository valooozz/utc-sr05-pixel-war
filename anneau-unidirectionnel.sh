mkfifo /tmp/in_A1 /tmp/out_A1
mkfifo /tmp/in_C1 /tmp/out_C1

mkfifo /tmp/in_A2 /tmp/out_A2
mkfifo /tmp/in_C2 /tmp/out_C2

mkfifo /tmp/in_A3 /tmp/out_A3
mkfifo /tmp/in_C3 /tmp/out_C3

go run app-base -n A1 -m g< /tmp/in_A1 > /tmp/out_A1 &
#go run app-base -n A1< /tmp/in_A1 > /tmp/out_A1 &
go run app-control -n C1 < /tmp/in_C1 > /tmp/out_C1 &

go run app-base -n A2 -m g< /tmp/in_A2 > /tmp/out_A2 &
#go run app-base -n A2< /tmp/in_A2 > /tmp/out_A2 &
go run app-control -n C2 < /tmp/in_C2 > /tmp/out_C2 &

go run app-base -n A3 -m g< /tmp/in_A3 > /tmp/out_A3 &
#go run app-base -n A3< /tmp/in_A3 > /tmp/out_A3 &
go run app-control -n C3 < /tmp/in_C3 > /tmp/out_C3 &

cat /tmp/out_A1 > /tmp/in_C1 &
cat /tmp/out_C1 | tee /tmp/in_A1 > /tmp/in_C2 &

cat /tmp/out_A2 > /tmp/in_C2 &
cat /tmp/out_C2 | tee /tmp/in_A2 > /tmp/in_C3 &

cat /tmp/out_A3 > /tmp/in_C3 &
cat /tmp/out_C3 | tee /tmp/in_A3 > /tmp/in_C1 &